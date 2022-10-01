package v1

import "github.com/gin-gonic/gin"

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

}

func (h *Handler) withdrawBalance(c *gin.Context) {

}
