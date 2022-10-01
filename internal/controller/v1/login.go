package v1

import (
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
    "login": "<login>",
    "password": "<password>"
}

Возможные коды ответа:
    200 — пользователь успешно аутентифицирован;
    400 — неверный формат запроса;
    401 — неверная пара логин/пароль;
    500 — внутренняя ошибка сервера.
*/

func (h *Handler) loginUser(c *gin.Context) {
	var input entity.User
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest) // Другой возможный вариант - c.AbortWithError
		h.logger.Error(fmt.Errorf("handler - login user: %w", err))
		return
	}

	token, err := h.services.GenerateToken(input.Login, input.Password)
	// todo: Выловить ошибку повторяющегося логина (SQLSTATE 23505)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - login user: %w", err))

		if strings.Contains(err.Error(), "no rows") {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError) // Другой возможный вариант - c.AbortWithError
		}

		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
