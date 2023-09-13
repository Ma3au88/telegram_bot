package delivery

import (
	"TestTelegramBot/internal/config"
	"TestTelegramBot/internal/service"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
	services   *service.Service
}

func NewServer(cfg *config.Config, services *service.Service) *Server {
	return &Server{
		config:   cfg,
		services: services,
	}
}

func (s *Server) Run(ctx context.Context, port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		if err := s.autoUpdatesCurrency(ctx, s.config.CurrencyUpdate); err != nil {
			logrus.Fatalf("updating currency or save in database: %v", err)
		}
	}()

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// autoUpdatesPrice создает тикер, который забирает данные о валюте и сохраняет в БД
func (s *Server) autoUpdatesCurrency(ctx context.Context, timeUpdate time.Duration) error {
	ticker := time.NewTicker(timeUpdate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancel signal: %w", ctx.Err())
		case <-ticker.C:
			price, err := s.services.GetPrices.GetPrice(ctx, []string{BTC, ETH})
			if err != nil {
				return fmt.Errorf("getting currency data from API: %w", err)
			}
			if err := s.services.Currency.SaveCurrency(ctx, price); err != nil {
				return fmt.Errorf("database save error: %w", err)
			}
		}
	}
}
