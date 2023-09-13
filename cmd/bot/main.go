package main

import (
	"TestTelegramBot/internal/config"
	"TestTelegramBot/internal/delivery"
	"TestTelegramBot/internal/repository"
	"TestTelegramBot/internal/repository/postgres"
	"TestTelegramBot/internal/service"
	"TestTelegramBot/internal/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("configuration error: %s", err.Error())
	}

	// Инициализация Telegram бота
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logrus.Fatalf("invalid Telegram token: %s", err.Error())
	}

	// Установка Debug в true, чтобы получать информацию об ошибках
	bot.Debug = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfgDB := postgres.Config{
		Host:         cfg.DBHost,
		Port:         cfg.DBPort,
		Username:     cfg.DBUsername,
		Name:         cfg.DBName,
		Password:     cfg.DBPassword,
		SSLmode:      cfg.DBSSLmode,
		PoolMaxConns: cfg.DBPoolMaxConns,
	}

	db, err := postgres.InitDB(ctx, cfgDB)
	if err != nil {
		logrus.Fatalf("error connecting to database: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := delivery.NewHandler(services)
	srv := delivery.NewServer(cfg, services)
	telegramBot := telegram.NewBot(bot, services)

	logrus.Println("server started")
	go func() {
		if err := srv.Run(ctx, cfg.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Println("Telegram bot started")
	go func() {
		if err := telegramBot.Start(ctx); err != nil {
			logrus.Fatalf("Telegram bot start error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Telegram bot Shutting Down")

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("error shutting down http server: %s", err.Error())
	}
}
