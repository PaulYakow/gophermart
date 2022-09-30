package v1

import (
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/usecase"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const (
	rootRoute = "api/user"
)

func NewRouter(uc usecase.IServer, l logger.ILogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(gin.Logger())
	handler.Use(gzip.Gzip(gzip.BestCompression))

	//root := handler.Group(rootRoute)

	return handler
}
