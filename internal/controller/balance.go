package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
*/

func (h *Handler) getBalance(c *gin.Context) {
	userID, ok := c.Get(userIDKey)
	if !ok {
		h.logger.Error(fmt.Errorf("get balance: user id not found"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("get balance: user id not found")))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	balance, err := h.services.GetBalance(ctx, userID.(int))
	if err != nil {
		h.logger.Error(fmt.Errorf("get uploaded orders: invalid request body: %w", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, balance)
}
