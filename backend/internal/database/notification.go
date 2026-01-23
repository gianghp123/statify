package database

import (
	"context"
	"log"
	"time"

	"github.com/gianghp/statify/internal/core/sse"
	"github.com/jackc/pgx/v5"
)

type PostgresNotificationListener struct {
	connStr string
	broker  *sse.Broker
}

func NewPostgresNotificationListener(connStr string, broker *sse.Broker) *PostgresNotificationListener {
	return &PostgresNotificationListener{
		connStr: connStr,
		broker:  broker,
	}
}

func (l *PostgresNotificationListener) Run(ctx context.Context, topic string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Postgres listener stopped")
			return
		default:
		}

		if err := l.listenOnce(ctx, topic); err != nil {
			log.Printf("Listener error: %v. Reconnecting in 5s...", err)
			time.Sleep(5 * time.Second)
		}
	}

}

func (l *PostgresNotificationListener) listenOnce(ctx context.Context, topic string) error {
	conn, err := pgx.Connect(ctx, l.connStr)

	if err != nil {
		log.Printf("Listener connection error: %v", err)
		return err
	}
	defer conn.Close(ctx)

	conn.Exec(ctx, "LISTEN "+topic)

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			log.Printf("Error waiting for notification: %v", err)
			return err
		}

		l.broker.Notify(notification.Payload, topic)
	}
}
