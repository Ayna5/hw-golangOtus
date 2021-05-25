package rabbitmq

import (
	"fmt"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Producer struct {
	cfg     configs.MQ
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(cfg configs.MQ) (*Producer, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("cannot connect producer: %w", err)
	}
	return &Producer{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (p *Producer) CloseConn() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("cannot close connection: %w", err)
	}

	return nil
}

func (p *Producer) OpenChannel() error {
	var err error
	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("cannot open channel: %w", err)
	}
	if err = p.channel.ExchangeDeclare(
		p.cfg.ExchangeName,
		p.cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("cannot exchange declare: %w", err)
	}

	return nil
}

func (p *Producer) CloseChannel() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("cannot close channel: %w", err)
	}

	return nil
}

func (p *Producer) Send(body []byte) error {
	if p.channel == nil {
		return errors.New("channel is nil")
	}
	if err := p.channel.Publish(
		p.cfg.ExchangeName,
		p.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Transient,
			Priority:     0,
		},
	); err != nil {
		return fmt.Errorf("cannot exchange publish: %w", err)
	}

	return nil
}
