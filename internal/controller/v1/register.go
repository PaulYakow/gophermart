package v1

import (
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
    "login": "<login>",
    "password": "<password>"
}

Возможные коды ответа:
    200 — пользователь успешно зарегистрирован и аутентифицирован;
    400 — неверный формат запроса;
    409 — логин уже занят;
    500 — внутренняя ошибка сервера.
*/

func (h *Handler) registerUser(c *gin.Context) {
	var input entity.User
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest) // Другой возможный вариант - c.AbortWithError
		h.logger.Error(fmt.Errorf("handler - register user: %w", err))
		return
	}

	_, err := h.services.CreateUser(input)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - register user: %w", err))

		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			c.AbortWithStatus(http.StatusConflict)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError) // Другой возможный вариант - c.AbortWithError
		}

		return
	}

	token, err := h.services.GenerateToken(input.Login, input.Password)
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - login user: %w", err))

		if strings.Contains(err.Error(), "no rows") {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError) // Другой возможный вариант - c.AbortWithError
		}

		return
	}

	c.Header(authorizationHeader, token)
	c.Status(http.StatusOK)
}
