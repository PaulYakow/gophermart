package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		h.logger.Error(fmt.Errorf("user identity: empty auth header"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := h.services.ParseToken(header)
	if err != nil {
		h.logger.Error(fmt.Errorf("user identity: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(userCtx, userID)
}
