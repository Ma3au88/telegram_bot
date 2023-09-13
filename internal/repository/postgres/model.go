package postgres

import "time"

type CryptoPrices struct {
	Id    int       `json:"-" db:"id"`
	Coin  string    `json:"coin" db:"coin"`
	Value float64   `json:"value" db:"value"`
	Date  time.Time `json:"-" db:"data"`
}

type TelegramBot struct {
	id     int
	chatID int64
	date   time.Time
}

type Subscribers struct {
	Id              int           `db:"id"`
	ChatID          int64         `db:"chat_id"`
	Interval        time.Duration `db:"interval"`
	LastMessageTime time.Time     `db:"last_message_time"`
	Status          bool          `db:"status"`
}
