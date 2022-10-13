package v1

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

func (h *Handler) withdrawBalance(c *gin.Context) {

}
