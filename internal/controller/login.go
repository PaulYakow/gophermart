package controller

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	loginRequest, ok := c.Get(loginUserReqKey)
	if !ok {
		h.logger.Error(fmt.Errorf("login user: user not found"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("login user: user not found")))
		return
	}
	user := loginRequest.(entity.UserDTO)

	userID, err := h.services.GetUser(user.Login, user.Password)
	if err != nil {
		h.logger.Error(fmt.Errorf("login user: %w", err))

		if errors.Is(err, service.ErrLoginNotExist) || errors.Is(err, service.ErrMismatchPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Set(userIDKey, userID)
}
