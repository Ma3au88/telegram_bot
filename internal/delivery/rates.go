package delivery

import (
	"TestTelegramBot/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	BTC = "BTC"
	ETH = "ETH"
)

type getAllRatesResponse struct {
	Data []postgres.CryptoPrices `json:"data"`
}

func (h *Handler) GetAllRates(c *gin.Context) {
	coins := []string{BTC, ETH}

	price, err := h.services.Currency.GetCurrency(c, coins)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRatesResponse{
		Data: price,
	})
}

func (h *Handler) GetRatesById(c *gin.Context) {
	coin := []string{c.Param("id")}

	price, err := h.services.Currency.GetCurrency(c, coin)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRatesResponse{
		Data: price,
	})
}
