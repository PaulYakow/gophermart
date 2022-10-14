package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/PaulYakow/gophermart/internal/util/workerpool"
	"github.com/imroc/req/v3"
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
	pool       *workerpool.Pool
	repo       repo.IOrder
	logger     logger.ILogger
	httpclient *req.Client
	endpoint   string
}

func NewPollService(repo repo.IOrder, logger logger.ILogger, endpoint string) *PollService {
	pollingLogger := logger.Named("polling")

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
			return 2 * time.Second
		}).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if resp.Response != nil {
				result := PollResult{}
				if resp.StatusCode == http.StatusOK && resp.GetHeader("Content-Type") == "application/json" {
					err := resp.UnmarshalJson(&result)
					if err != nil {
						pollingLogger.Error(fmt.Errorf("client - CommonRetryCondition unmarshal error: %w", err))
						return true
					}
				}
				pollingLogger.Info("client - retry attempt: ", result)

				return err != nil ||
					resp.StatusCode == http.StatusTooManyRequests ||
					resp.StatusCode == http.StatusNoContent ||
					!isFinal(result.Status)
			}
			return true
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			rq := resp.Request.RawRequest
			pollingLogger.Info("client - retry request: ", rq.Method, rq.URL)
		})

	return &PollService{
		pool:       workerpool.NewPool(2, 10),
		repo:       repo,
		logger:     pollingLogger,
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
		s.logger.Error(fmt.Errorf("requestOrderInfo - error: %w", err))
		s.logger.Info("raw content: ", resp.Dump()) // Record raw content when error occurs.
		return nil, err
	}
	//defer resp.Body.Close()

	if resp.IsSuccess() { // Status code is between 200 and 299.
		s.logger.Info("requestOrderInfo - success: result=%v | code=%s\n", result, resp.Status)
		return result, nil
	}

	// Unknown status code.
	s.logger.Info("requestOrderInfo - unknown status:", resp.Status)
	s.logger.Info("raw content: ", resp.Dump()) // Record raw content when error occurs.

	return nil, nil
}

func (s *PollService) getResults(ctx context.Context) {
	for {
		select {
		case r, ok := <-s.pool.Results():
			if !ok {
				s.logger.Info("getResults <-s.pool.Results(): ", ok)
				return
			}
			s.logger.Info("getResults: ", r)

			if r.Err != nil {
				s.logger.Error(fmt.Errorf("getResults - r.Err: %w", r.Err))
			}

			result, ok := r.Value.(PollResult)
			if !ok {
				s.logger.Error(fmt.Errorf("getResults - cannot convert to PollResult: %v", result))
				continue
			}

			if err := s.repo.UpdateUploadedOrder(result.Order, result.Status, result.Accrual); err != nil {
				s.logger.Error(fmt.Errorf("getResults - cannot update order: %w", err))
				continue
			}
			s.logger.Info("getResults - update upload order success")

		case <-ctx.Done():
			s.logger.Info("getResults - context canceled: %v", ctx.Err())
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
