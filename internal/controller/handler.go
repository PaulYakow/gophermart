package controller

import (
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const (
	rootRoute        = "api/user"
	registerRoute    = "/register"
	loginRoute       = "/login"
	ordersRoute      = "/orders"
	balanceRoute     = "/balance"
	withdrawRoute    = "/balance/withdraw"
	withdrawalsRoute = "/withdrawals"
)

type Handler struct {
	services *service.Service
	logger   logger.ILogger
}

func NewHandler(services *service.Service, logger logger.ILogger) *Handler {
	return &Handler{
		services: services,
		logger:   logger.Named("handler"),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(gin.Logger())
	handler.Use(gzip.Gzip(gzip.BestCompression))

	root := handler.Group(rootRoute)
	{
		auth := root.Group("/").Use(h.userAuthentication)
		{
			auth.POST(registerRoute, h.registerUser)
			auth.POST(loginRoute, h.loginUser)
		}

		requireAuth := root.Group("/").Use(h.userIdentity)
		{
			requireAuth.POST(ordersRoute, h.loadOrder)
			requireAuth.GET(ordersRoute, h.getListOfOrders)
			requireAuth.GET(balanceRoute, h.getBalance)
			requireAuth.POST(withdrawRoute, h.withdrawOrder)
			requireAuth.GET(withdrawalsRoute, h.withdrawInfo)
		}
	}

	return handler
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
