package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
)

type Consumer struct {
	ctx      context.Context
	log      *logger.Logger
	cfg      configs.MQ
	conn     *amqp.Connection
	channel  *amqp.Channel
	delivery <-chan amqp.Delivery //nolint:structcheck,unused
}

func NewConsumer(ctx context.Context, cfg configs.MQ, log *logger.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("cannot connect consumer: %w", err)
	}
	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot create a amqpChannel: %w", err)
	}
	return &Consumer{
		ctx:     ctx,
		log:     log,
		cfg:     cfg,
		channel: amqpChannel,
		conn:    conn,
	}, nil
}

func (c *Consumer) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("cannot close connection: %w", err)
	}
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("cannot close channel: %w", err)
	}

	return nil
}

func (c *Consumer) OpenChannel() (<-chan amqp.Delivery, error) {
	var err error
	if err = c.channel.ExchangeDeclare(
		c.cfg.ExchangeName,
		c.cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("cannot exchange declare: %w", err)
	}
	queue, err := c.channel.QueueDeclare(
		c.cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot queue declare: %w", err)
	}
	err = c.channel.Qos(1, 0, false)
	if err != nil {
		return nil, fmt.Errorf("could not configure QoS: %w", err)
	}
	if err = c.channel.QueueBind(
		queue.Name,
		c.cfg.RoutingKey,
		c.cfg.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("cannot queue bind: %w", err)
	}
	deliveries, err := c.channel.Consume(
		queue.Name,
		c.cfg.Tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.log.Error("cannot consume channel: " + err.Error())
		return nil, fmt.Errorf("cannot consume channel: %w", err)
	}

	return deliveries, nil
}

func (c *Consumer) ReadMsg(delivery <-chan amqp.Delivery, done chan error) {
	var e storage.Event
	c.log.Info("Consumer ready, PID: " + strconv.Itoa(os.Getpid()))
	for d := range delivery {
		c.log.Info("Received a message: " + string(d.Body))
		if err := json.Unmarshal(d.Body, &e); err != nil {
			c.log.Error("unmarshal error for event: " + e.ID + " " + err.Error())
			continue
		}

		if err := d.Ack(false); err != nil {
			c.log.Error("Error acknowledging message : " + err.Error())
		} else {
			c.log.Info("Acknowledged message")
		}
	}
	done <- nil
}
