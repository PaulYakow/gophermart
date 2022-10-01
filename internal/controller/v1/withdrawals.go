package v1

import "github.com/gin-gonic/gin"

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

}
