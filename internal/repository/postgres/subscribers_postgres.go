package postgres

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) StartAuto(ctx context.Context, subs Subscribers) error {
	var id int64
	query := `INSERT INTO
				subscribers (chat_id, interval, last_message_time, status)
			VALUES 
			    ($1, $2, $3, $4)
			ON CONFLICT (chat_id)
			DO UPDATE SET 
			    interval = EXCLUDED.interval, last_message_time = EXCLUDED.last_message_time, status = EXCLUDED.status 
			RETURNING id;`

	row := r.db.QueryRow(ctx, query, subs.ChatID, subs.Interval, subs.LastMessageTime, subs.Status)

	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("database scan: %w", err)
	}

	return nil
}

func (r *Repository) StopAuto(ctx context.Context, chatID int64) error {
	query := `UPDATE subscribers SET status = FALSE WHERE chat_id = $1;`

	_, err := r.db.Exec(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("sending query to database: %w", err)
	}

	return nil
}

func (r *Repository) AutoUpdateUser(ctx context.Context) ([]Subscribers, error) {
	query := `SELECT
				chat_id
			FROM 
			    subscribers
			WHERE 
			    status = TRUE
			AND
			    (interval + last_message_time) <= $1;`

	var sub Subscribers
	subs := make([]Subscribers, 0)

	row, err := r.db.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("sending query to database: %w", err)
	}

	for row.Next() {
		err := row.Scan(&sub.ChatID)
		if err != nil {
			return nil, fmt.Errorf("database scan: %w", err)
		}

		subs = append(subs, sub)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return subs, nil
}

func (r *Repository) UpdateLastMessageTime(ctx context.Context, chatID int64) error {
	query := `UPDATE subscribers SET last_message_time = $1 WHERE chat_id = $2;`

	_, err := r.db.Exec(ctx, query, time.Now(), chatID)
	if err != nil {
		return fmt.Errorf("sending query to database: %w", err)
	}

	return nil
}
