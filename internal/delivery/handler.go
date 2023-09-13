package delivery

import (
	"TestTelegramBot/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		rates := api.Group("/rates")
		{
			rates.GET("/", h.GetAllRates)
			rates.GET("/:id", h.GetRatesById)
		}
	}

	return router
}
