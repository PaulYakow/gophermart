package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0

Возможные коды ответа:
    200 — успешная обработка запроса.
    Формат ответа:

  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
      {
          "order": "2377225624",
          "sum": 500,
          "processed_at": "2020-12-09T16:09:57+03:00"
      }
  ]

204 — нет ни одного списания.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/

func (h *Handler) withdrawInfo(c *gin.Context) {
	userID, ok := c.Get(userIDKey)
	if !ok {
		h.logger.Error(fmt.Errorf("upload order: user id not found"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("upload order: user id not found")))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	withdrawOrders, err := h.services.GetWithdrawOrders(ctx, userID.(int))
	if err != nil {
		h.logger.Error(fmt.Errorf("get uploaded orders: invalid request body: %w", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, withdrawOrders)
}
