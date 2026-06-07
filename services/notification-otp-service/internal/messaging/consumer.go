package messaging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"notification_otp_service/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	queue   string
	handler *events.Handler
	logger  *slog.Logger
}

func NewConsumer(conn *amqp.Connection, queue string, handler *events.Handler, logger *slog.Logger) *Consumer {
	return &Consumer{conn: conn, queue: queue, handler: handler, logger: logger}
}

func (c *Consumer) Run(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case delivery, ok := <-messages:
			if !ok {
				return nil
			}
			c.handleDelivery(delivery)
		}
	}
}

func (c *Consumer) handleDelivery(delivery amqp.Delivery) {
	if err := withRetry(3, 100*time.Millisecond, func() error { return c.handler.Handle(delivery.Body) }); err != nil {
		if errors.Is(err, events.ErrDuplicate) {
			ack(delivery, c.logger)
			return
		}
		c.logger.Error("message processing failed", "error", err)
		nack(delivery, c.logger)
		return
	}
	ack(delivery, c.logger)
}

func withRetry(attempts int, delay time.Duration, fn func() error) error {
	var last error
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			last = err
			time.Sleep(delay * time.Duration(i+1))
			continue
		}
		return nil
	}
	return last
}

func ack(delivery amqp.Delivery, logger *slog.Logger) {
	if err := delivery.Ack(false); err != nil {
		logger.Error("ack message", "error", err)
	}
}

func nack(delivery amqp.Delivery, logger *slog.Logger) {
	if err := delivery.Nack(false, true); err != nil {
		logger.Error("nack message", "error", err)
	}
}
