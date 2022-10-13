package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

/*
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903

Возможные коды ответа:
    200 — номер заказа уже был загружен этим пользователем;
    202 — новый номер заказа принят в обработку;
    400 — неверный формат запроса;
    401 — пользователь не аутентифицирован;
    409 — номер заказа уже был загружен другим пользователем;
    422 — неверный формат номера заказа;
    500 — внутренняя ошибка сервера.
-----------------------------------------------------------

GET /api/user/orders HTTP/1.1
Content-Length: 0

Доступные статусы обработки расчётов:
    NEW — заказ загружен в систему, но не попал в обработку;
    PROCESSING — вознаграждение за заказ рассчитывается;
    INVALID — система расчёта вознаграждений отказала в расчёте;
    PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.

Возможные коды ответа:
    200 — успешная обработка запроса.
    Формат ответа:

  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
      {
          "number": "9278923470",
          "status": "PROCESSED",
          "accrual": 500,
          "uploaded_at": "2020-12-10T15:15:45+03:00"
      },
  ]

204 — нет данных для ответа.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
---------------------------------------

POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
    "order": "2377225624",
    "sum": 751
}

Возможные коды ответа:
    200 — успешная обработка запроса;
    401 — пользователь не авторизован;
    402 — на счету недостаточно средств;
    422 — неверный номер заказа;
    500 — внутренняя ошибка сервера.
*/

func (h *Handler) loadOrder(c *gin.Context) {
	h.logger.Info("handler - request: %v", *c.Request)

	if c.Request.Header.Get("Content-Type") != "text/plain" {
		h.logger.Error(fmt.Errorf("handler - upload order: content-type not text/plain"))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	number, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - load order: invalid request body: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(string(number))
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - load order: cannot convert data in request body: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !checkOrderNumber(orderNumber) {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	userID, ok := c.Get(userCtx)
	if !ok {
		h.logger.Error(fmt.Errorf("handler - upload order: user id not found"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//h.logger.Info("user_id: %v | order number: %v", userID.(int), orderNumber)
	userIDInOrder, err := h.services.CreateUploadedOrder(userID.(int), orderNumber)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - upload order: failed create in storage: %w", err))
		return
	}

	if userIDInOrder == 0 {
		h.services.Polling.AddToPoll(strconv.Itoa(orderNumber))

		h.logger.Info("handler - upload order: order accepted")
		c.Status(http.StatusAccepted)
	} else if userIDInOrder == userID.(int) {
		h.logger.Info("handler - upload order: order has already been loaded by this user")
		c.Status(http.StatusOK)
	} else {
		h.logger.Info("handler - upload order: order has already been loaded by another user")
		c.Status(http.StatusConflict)
	}
}

func (h *Handler) getListOfOrders(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		h.logger.Error(fmt.Errorf("handler - upload order: user id not found"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	uploadedOrders, err := h.services.GetUploadedOrders(ctx, userID.(int))
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - get uploaded orders: invalid request body: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, uploadedOrders)
}

func (h *Handler) withdrawOrder(c *gin.Context) {
	var withdraw WithdrawRequest
	if err := c.BindJSON(&withdraw); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		h.logger.Error(fmt.Errorf("handler - register withdraw: %w", err))
		return
	}

	orderNumber, err := strconv.Atoi(string(withdraw.Order))
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - register withdraw: cannot convert data in request body: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !checkOrderNumber(orderNumber) {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	userID, ok := c.Get(userCtx)
	if !ok {
		h.logger.Error(fmt.Errorf("handler - register withdraw: user id not found"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.services.CreateWithdrawOrder(userID.(int), withdraw.Order, withdraw.Sum)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - register withdraw: failed create in storage: %w", err))

		if errors.Is(err, repo.ErrNoFunds) {
			c.AbortWithStatus(http.StatusPaymentRequired)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func checkOrderNumber(orderNumber int) bool {
	var luhn int
	number := orderNumber / 10

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}

	return (orderNumber%10+luhn%10)%10 == 0
}
