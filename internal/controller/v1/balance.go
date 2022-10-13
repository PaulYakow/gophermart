package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*
GET /api/user/balance HTTP/1.1
Content-Length: 0

Возможные коды ответа:
    200 — успешная обработка запроса.
    Формат ответа:

  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  {
      "current": 500.5,
      "withdrawn": 42
  }

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

func (h *Handler) getBalance(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		h.logger.Error(fmt.Errorf("handler - upload order: user id not found"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	balance, err := h.services.GetBalance(ctx, userID.(int))
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - get uploaded orders: invalid request body: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance)
}

type WithdrawRequest struct {
	Order string  `json:"order" db:"number"`
	Sum   float32 `json:"sum" db:"sum"`
}

func (h *Handler) withdrawBalance(c *gin.Context) {
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

	err = h.services.UpdateWithdrawBalance(userID.(int), orderNumber, withdraw.Sum)
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
