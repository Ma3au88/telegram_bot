package postgres

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) SaveUser(ctx context.Context, chatID int64) error {
	var id int64
	query := `INSERT INTO users (chat_id, date) VALUES ($1, $2) RETURNING id`

	row := r.db.QueryRow(ctx, query, chatID, time.Now())

	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("database scan: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, chatID int64) (int64, error) {
	query := `SELECT chat_id FROM users WHERE chat_id = $1;`

	var user TelegramBot

	row, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return 0, fmt.Errorf("sending query to database: %w", err)
	}

	for row.Next() {
		err := row.Scan(&user.chatID)
		if err != nil {
			return 0, fmt.Errorf("database scan: %w", err)
		}
	}

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}

	return user.chatID, nil
}
