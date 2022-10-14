package v1

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
		root.POST(registerRoute, h.registerUser)
		root.POST(loginRoute, h.loginUser)

		auth := root.Group("/").Use(h.userIdentity)
		{
			auth.POST(ordersRoute, h.loadOrder)
			auth.GET(ordersRoute, h.getListOfOrders)
			auth.GET(balanceRoute, h.getBalance)
			auth.POST(withdrawRoute, h.withdrawOrder)
			auth.GET(withdrawalsRoute, h.withdrawInfo)
		}
	}

	return handler
}
