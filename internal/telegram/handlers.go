package telegram

import (
	"TestTelegramBot/internal/repository/postgres"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

const (
	commandStart     = "start"
	commandStartAuto = "start_auto"
	commandStopAuto  = "stop_auto"
	commandHelp      = "help"
	commandRates     = "rates"
	commandRate      = "rate"

	BTC = "BTC"
	ETH = "ETH"
)

func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(ctx, message)
	case commandStartAuto:
		return b.handleStartAutoCommand(ctx, message)
	case commandStopAuto:
		return b.handleStopAutoCommand(ctx, message)
	case commandRate:
		return b.handleRateCommand(ctx, message)
	case commandRates:
		return b.handleRatesCommand(ctx, message)
	case commandHelp:
		return b.handleHelpCommand(ctx, message)
	default:
		return b.handleUnknownCommand(ctx, message)
	}
}

func (b *Bot) handleStartCommand(ctx context.Context, message *tgbotapi.Message) error {
	user, err := b.services.User.GetUser(ctx, message.Chat.ID)
	if err != nil {
		return err //логику проверки юзера сюда положить
	}

	if user == 0 {
		if err := b.services.User.SaveUser(ctx, message.Chat.ID); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы добавлены в базу пользователей бота")
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}
	}

	return nil
}

func (b *Bot) handleStartAutoCommand(ctx context.Context, message *tgbotapi.Message) error {
	duration, err := strconv.Atoi(message.CommandArguments())
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введено неверное число. "+
			"Введите корректное число минут, через которое трубется обновление валют")
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}
	} else {
		user := postgres.Subscribers{
			ChatID:          message.Chat.ID,
			Interval:        time.Duration(duration) * time.Minute,
			LastMessageTime: time.Now(),
			Status:          true,
		}

		if err := b.services.Subscribers.StartAuto(ctx, user); err != nil {
			return fmt.Errorf("adding to subscribers: %w", err)
		}
	}

	return nil
}

func (b *Bot) handleStopAutoCommand(ctx context.Context, message *tgbotapi.Message) error {
	if err := b.services.Subscribers.StopAuto(ctx, message.Chat.ID); err != nil {
		return fmt.Errorf("removal from subscribers: %w", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы отключены от подписки")
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message to user: %w", err)
	}

	return nil
}

func (b *Bot) handleRateCommand(ctx context.Context, message *tgbotapi.Message) error {
	userCoin := message.CommandArguments()
	coin := []string{userCoin}

	price, err := b.services.Currency.GetCurrency(ctx, coin)
	if err != nil {
		return fmt.Errorf("getting currency data: %w", err)
	}

	for _, p := range price {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Цена %s в RUB: %.2f\n", p.Coin, p.Value))
		_, err = b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}
	}

	return nil
}

func (b *Bot) handleRatesCommand(ctx context.Context, message *tgbotapi.Message) error {
	prices, err := b.services.Currency.GetCurrency(ctx, []string{BTC, ETH})
	if err != nil {
		return fmt.Errorf("getting currency data: %w", err)
	}

	var msgArr []string
	for _, price := range prices {
		msgArr = append(msgArr, fmt.Sprintf("Цена %s в RUB: %.2f\n", price.Coin, price.Value))
	}

	msgDuo := strings.Join(msgArr, "")
	msg := tgbotapi.NewMessage(message.Chat.ID, msgDuo)
	_, err = b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message to user: %w", err)
	}

	return nil
}

func (b *Bot) handleHelpCommand(ctx context.Context, message *tgbotapi.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancel signal: %w", ctx.Err())
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вот список доступных команд:\n/start - Начать\n"+
			"/help - Помощь\n/rates - Получение курса валют BTC и ETH\n/rate {currency} - Получение курса {currency} "+
			"валюты\n/start_auto {число в минутах} - Подписка на автоматическое оповещение о валютах"+
			" через {число в минутах}\n/stop_auto - Отключение от автоматического оповещения")
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}
	}

	return nil
}

func (b *Bot) handleUnknownCommand(ctx context.Context, message *tgbotapi.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancel signal: %w", ctx.Err())
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. "+
			"Попробуйте /help для получения списка команд")
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}
	}

	return nil
}
