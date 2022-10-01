package v1

import (
	"github.com/gin-gonic/gin"
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
*/

func (h *Handler) loadOrder(c *gin.Context) {
	//userId, ok := c.Get(userCtx)
	//if !ok {
	//	h.logger.Error(fmt.Errorf("handler - load order: user id not found"))
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}

}

func (h *Handler) getListOfOrders(c *gin.Context) {

}