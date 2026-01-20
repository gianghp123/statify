package database

import (
	"context"
	"log"

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

func (l *PostgresNotificationListener) Run(topic string) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, l.connStr)
	if err != nil {
		log.Printf("Listener connection error: %v", err)
		return
	}
	defer conn.Close(ctx)

	conn.Exec(ctx, "LISTEN "+topic)

	for {
		// This line blocks the goroutine until a message arrives
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			log.Printf("Error waiting for notification: %v", err)
			// In production, add a "sleep and retry" here to handle DB restarts
			return
		}

		// Process the message (e.g., send to a Go channel or trigger a function)
		l.broker.Notify(notification.Payload, topic)
	}
}
