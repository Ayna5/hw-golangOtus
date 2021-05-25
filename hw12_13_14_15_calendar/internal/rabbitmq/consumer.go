package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/streadway/amqp"
)

type Consumer struct {
	cfg     configs.MQ
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewConsumer(cfg configs.MQ) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("cannot connect consumer: %w", err)
	}
	return &Consumer{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (c *Consumer) CloseConn() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("cannot close connection: %w", err)
	}

	return nil
}

func (c *Consumer) OpenChannel() error {
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("cannot open channel: %w", err)
	}
	if err = c.channel.ExchangeDeclare(
		c.cfg.ExchangeName,
		c.cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("cannot exchange declare: %w", err)
	}
	c.queue, err = c.channel.QueueDeclare(
		c.cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("cannot queue declare: %w", err)
	}
	if err = c.channel.QueueBind(
		c.queue.Name,
		c.cfg.RoutingKey,
		c.cfg.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("cannot queue bind: %w", err)
	}

	return nil
}

func (c *Consumer) ReadMsg() (configs.MQNotification, error) {
	deliveries, err := c.channel.Consume(
		c.queue.Name,
		c.cfg.Tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return configs.MQNotification{}, fmt.Errorf("cannot consume channel: %w", err)
	}
	var event configs.MQNotification
	for d := range deliveries {
		if err = json.Unmarshal(d.Body, &event); err != nil {
			continue
		}
		_ = d.Ack(false)
	}

	return event, nil
}

func (c *Consumer) CloseChannel() error {
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("cannot close channel: %w", err)
	}

	return nil
}
