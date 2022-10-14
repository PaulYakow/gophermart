package v1

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
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
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest) // Другой возможный вариант - c.AbortWithError
		h.logger.Error(fmt.Errorf("register user: %w", err))
		return
	}

	_, err := h.services.CreateUser(user)
	if err != nil {

		if errors.Is(err, repo.ErrDuplicateKey) {
			h.logger.Error(fmt.Errorf("register user: login already exists"))
			c.AbortWithStatus(http.StatusConflict)
		} else {
			h.logger.Error(fmt.Errorf("register user: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	token, err := h.services.GenerateToken(user.Login, user.Password)
	if err != nil {
		h.logger.Error(fmt.Errorf("login user: %w", err))

		if strings.Contains(err.Error(), "no rows") {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.Header(authorizationHeader, token)
	c.Status(http.StatusOK)
}
