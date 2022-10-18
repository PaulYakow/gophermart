package controller

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/gin-gonic/gin"
	"net/http"
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
	loginRequest, ok := c.Get(loginUserReqKey)
	if !ok {
		h.logger.Error(fmt.Errorf("register: user not found"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("register: user not found")))
		return
	}
	user := loginRequest.(entity.UserDTO)

	userID, err := h.services.CreateUser(user.Login, user.Password)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicateKey) {
			h.logger.Error(fmt.Errorf("register: login already exists"))
			c.AbortWithStatusJSON(http.StatusConflict, errorResponse(err))
		} else {
			h.logger.Error(fmt.Errorf("register user: %w", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.Set(userIDKey, userID)
}
