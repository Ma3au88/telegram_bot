package telegram

import (
	"TestTelegramBot/internal/service"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	services *service.Service
}

func NewBot(bot *tgbotapi.BotAPI, serv *service.Service) *Bot {
	return &Bot{
		bot:      bot,
		services: serv,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	logrus.Printf("account autorization: %s", b.bot.Self.UserName)

	go func() {
		if err := b.sendingToSubscribers(ctx); err != nil {
			logrus.Fatalf("updating data from Subscribers database: %v", err)
		}
	}()

	updates := b.initUpdatesChannel()
	if err := b.handleUpdates(ctx, updates); err != nil {
		logrus.Fatalf("updating data from Telegram: %v", err)
	}

	return nil
}

// initUpdatesChannel создаёт обработчика команд
func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	return b.bot.GetUpdatesChan(updateConfig)
}

// В handleUpdates будут обрабатываться все новые сообщения
func (b *Bot) handleUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if err := b.handleCommand(ctx, update.Message); err != nil {
			return fmt.Errorf("receiving command from Telegram: %w", err)
		}
	}

	return nil
}

// sendingToSubscribers создает тикер и каждый тик происходит оповещение подписчиков о курсе валют
func (b *Bot) sendingToSubscribers(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancel signal: %w", ctx.Err())
		case <-ticker.C:
			subs, err := b.services.Subscribers.AutoUpdateUser(ctx)
			if err != nil {
				return fmt.Errorf("getting subscribers data: %w", err)
			}

			for _, sub := range subs {
				prices, err := b.services.Currency.GetCurrency(ctx, []string{BTC, ETH})
				if err != nil {
					return fmt.Errorf("getting currency data: %w", err)
				}

				var msgArr []string
				for _, price := range prices {
					msgArr = append(msgArr, fmt.Sprintf("Цена %s в RUB: %.2f\n", price.Coin, price.Value))
				}

				msgDuo := strings.Join(msgArr, "")
				msg := tgbotapi.NewMessage(sub.ChatID, msgDuo)
				_, err = b.bot.Send(msg)
				if err != nil {
					return fmt.Errorf("sending message to user: %w", err)
				}

				if err := b.services.Subscribers.UpdateLastMessageTime(ctx, sub.ChatID); err != nil {
					return fmt.Errorf("updating subscribers data: %w", err)
				}
			}
		}
	}
}
