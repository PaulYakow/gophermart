package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/PaulYakow/gophermart/internal/util/workerpool"
	"github.com/imroc/req/v3"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PollResult struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

type PollService struct {
	repo       repo.IUploadOrder
	pool       *workerpool.Pool
	httpclient *req.Client
	endpoint   string
}

func NewPollService(repo repo.IUploadOrder, endpoint string) *PollService {
	client := req.C().
		SetTimeout(1 * time.Second).
		SetCommonRetryCount(10).
		SetCommonRetryInterval(func(resp *req.Response, attempt int) time.Duration {
			if resp.Response != nil {
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					after, err := strconv.Atoi(ra)
					if err == nil {
						return time.Duration(after) * time.Second
					}
				}
			}
			return 5 * time.Second
		}).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if resp.Response != nil {
				result := PollResult{}
				if resp.StatusCode == http.StatusOK && resp.GetHeader("Content-Type") == "application/json" {
					err := resp.UnmarshalJson(&result)
					if err != nil {
						log.Println("PollService - CommonRetryCondition unmarshal error:", err)
						return true
					}
				}
				log.Println("PollService - retry attempt: ", result)

				return err != nil ||
					resp.StatusCode == http.StatusTooManyRequests ||
					resp.StatusCode == http.StatusNoContent ||
					!isFinal(result.Status)
			}
			return true
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			rq := resp.Request.RawRequest
			log.Println("PollService - retry request:", rq.Method, rq.URL)
		})

	return &PollService{
		repo:       repo,
		pool:       workerpool.NewPool(2, 10),
		httpclient: client,
		endpoint:   endpoint,
	}
}

func (s *PollService) Run(ctx context.Context) {
	go s.getResults(ctx)
	s.pool.RunBackground(ctx)
}

func (s *PollService) AddToPoll(number string) {
	/*
		pool add task:
			- Descriptor: ID = "<number>"| TType = ""| Metadata =
			- ExecFn = launch http request (with multiple retry cond - default time & retry-after)? / check status: (INVALID, PROCESSED) -> return PollResult type, else -> wait & retry
			- Args = number
	*/
	task := workerpool.NewTask(
		workerpool.TaskDescriptor{
			ID:       workerpool.TaskID(number),
			TType:    "",
			Metadata: nil,
		},
		s.requestOrderInfo,
		number)
	s.pool.AddTask(*task)
}

func (s *PollService) requestOrderInfo(ctx context.Context, args interface{}) (interface{}, error) {
	number, ok := args.(string)
	if !ok {
		return nil, errors.New("PollService wrong argument type for requestOrderInfo")
	}

	var result PollResult
	resp, err := s.httpclient.R().
		SetPathParam("number", number).
		SetResult(&result).
		EnableDump().
		Get(s.endpoint + "/{number}")

	if err != nil {
		log.Println("PollService - requestOrderInfo error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return nil, err
	}
	//defer resp.Body.Close()

	if resp.IsSuccess() { // Status code is between 200 and 299.
		fmt.Printf("PollService - requestOrderInfo success: result=%v | code=%s\n", result, resp.Status)
		return result, nil
	}

	// Unknown status code.
	log.Println("PollService - requestOrderInfo unknown status:", resp.Status)
	log.Println("raw content:")
	log.Println(resp.Dump())

	return nil, nil
}

func (s *PollService) getResults(ctx context.Context) {
	defer fmt.Println("PollService - getResults exiting")

	for {
		select {
		case r, ok := <-s.pool.Results():
			if !ok {
				log.Println("PollService - getResults <-s.pool.Results(): ", ok)
				return
			}
			log.Println("PollService - getResults: ", r)
			/*
				assert r.Value to PollResult type
				if err -> notify about error
				else -> send value to Results channel
			*/
			if r.Err != nil {
				log.Println("PollService - getResults r.Err: ", r.Err)
			}

			result, ok := r.Value.(PollResult)
			if !ok {
				log.Println("PollService - getResults cannot convert to PollResult: ", result)
				continue
			}

			if err := s.repo.UpdateUploadedOrder(result.Order, result.Status, result.Accrual); err != nil {
				log.Println("PollService - getResults cannot update order: ", err)
				continue
			}

			log.Println("PollService - getResults update order success: ", result)

		case <-ctx.Done():
			fmt.Printf("PollService - getResults context canceled: %v", ctx.Err())
			s.pool.Stop()
			return
		}
	}
}

var finalStatuses = map[string]struct{}{
	"INVALID":   {},
	"PROCESSED": {},
}

func isFinal(s string) bool {
	_, exist := finalStatuses[s]
	return exist
}
