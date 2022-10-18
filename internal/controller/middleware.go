package controller

import (
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey = "Authorization"
	userIDKey              = "user_id"
	loginUserReqKey        = "login_user_request"
)

func (h *Handler) userAuthentication(c *gin.Context) {
	var loginReq entity.UserDTO
	if err := c.BindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		h.logger.Error(fmt.Errorf("register user: %w", err))
		return
	}

	c.Set(loginUserReqKey, loginReq)
	c.Next()

	if !c.IsAborted() {
		userID, ok := c.Get(userIDKey)
		if !ok {
			h.logger.Error(fmt.Errorf("userAuthentication: user not found"))
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("userAuthentication: user not found")))
			return
		}

		token, err := h.services.GenerateToken(userID.(int))
		if err != nil {
			h.logger.Error(fmt.Errorf("login user: %w", err))

			if strings.Contains(err.Error(), "no rows") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			}

			return
		}

		c.Header(authorizationHeaderKey, token)
		c.Status(http.StatusOK)
	}
}

func (h *Handler) userIdentity(c *gin.Context) {
	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	if authorizationHeader == "" {
		h.logger.Error(fmt.Errorf("user identity: authorization header is not provided"))
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("user identity: authorization header is not provided")))
		return
	}

	userID, err := h.services.ParseToken(authorizationHeader)
	if err != nil {
		h.logger.Error(fmt.Errorf("user identity: %w", err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.Set(userIDKey, userID)
}
