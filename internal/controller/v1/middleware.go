package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		h.logger.Error(fmt.Errorf("handler - user identity: empty auth header"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		h.logger.Error(fmt.Errorf("handler - user identity: invalid auth header"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error(fmt.Errorf("handler - user identity: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(userCtx, userID)
}
