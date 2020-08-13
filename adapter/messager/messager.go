package messager

import (
	"fmt"

	"github.com/ewangplay/eventbus/adapter/messager/nsq"
	"github.com/ewangplay/eventbus/adapter/messager/rabbitmq"
	"github.com/ewangplay/eventbus/adapter/messager/solo"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

// Messager struct define, implement Messager interface
type Messager struct {
	i.Logger
	context   i.Context
	producer  i.Producer
	opts      *config.EBOptions
	consumers map[string]i.Consumer
}

// NewMessager ...
func NewMessager(opts *config.EBOptions, logger i.Logger) (*Messager, error) {
	m := &Messager{}

	m.opts = opts
	m.Logger = logger

	var err error
	var context i.Context
	var producer i.Producer
	if m.opts.NSQEnable {
		m.Info("Using NSQ messager ...")

		context, err = nsq.GetContext(opts, logger)
		if err != nil {
			m.Error("Get nsq context instance error: %v", err)
			return nil, err
		}

		producer, err = context.CreateProducer()
		if err != nil {
			m.Error("Create nsq producer error: %v", err)
			return nil, err
		}

	} else if m.opts.RabbitmqEnable {
		m.Info("Using Rabbitmq messager ...")

		context, err = rabbitmq.GetContext(opts, logger)
		if err != nil {
			m.Error("Get rabbitmq context instance error: %v", err)
			return nil, err
		}

		producer, err = context.CreateProducer()
		if err != nil {
			m.Error("Create rabbitmq producer error: %v", err)
			return nil, err
		}

	} else {
		m.Info("Using Solo messager ...")

		context, err = solo.GetContext(opts, logger)
		if err != nil {
			m.Error("Get solo context instance error: %v", err)
			return nil, err
		}

		producer, err = context.CreateProducer()
		if err != nil {
			m.Error("Create solo producer error: %v", err)
			return nil, err
		}

	}

	m.context = context
	m.producer = producer
	m.consumers = make(map[string]i.Consumer, 1)

	return m, nil
}

// Close ...
func (m *Messager) Close() error {
	m.Info("Messager will close")
	m.producer.Close()
	for _, consumer := range m.consumers {
		consumer.Close()
	}
	return nil
}

// Publish message to messager
func (m *Messager) Publish(msg i.Message) error {
	if m.producer == nil {
		return fmt.Errorf("producer not valid")
	}
	return m.producer.Publish(msg)
}

// Subscribe messages from messager based on subject
func (m *Messager) Subscribe(subject string) (<-chan i.Message, error) {

	//If consumer has already existed, get message chan directly
	oldConsumer, ok := m.consumers[subject]
	if ok {
		return oldConsumer.GetMessage(), nil
	}

	var err error
	var consumer i.Consumer
	if m.opts.NSQEnable {

		if m.context == nil {
			m.context, err = nsq.GetContext(m.opts, m.Logger)
			if err != nil {
				m.Error("Get nsq context instance error: %v", err)
				return nil, err
			}
		}

		consumer, err = m.context.CreateConsumer(subject)
		if err != nil {
			m.Error("Create nsq consumer error: %v", err)
			return nil, err
		}

	} else if m.opts.RabbitmqEnable {

		if m.context == nil {
			m.context, err = rabbitmq.GetContext(m.opts, m.Logger)
			if err != nil {
				m.Error("Get rabbitmq context instance error: %v", err)
				return nil, err
			}
		}

		consumer, err = m.context.CreateConsumer(subject)
		if err != nil {
			m.Error("Create rabbitmq consumer error: %v", err)
			return nil, err
		}

	} else {

		if m.context == nil {
			m.context, err = solo.GetContext(m.opts, m.Logger)
			if err != nil {
				m.Error("Get solo context instance error: %v", err)
				return nil, err
			}
		}

		consumer, err = m.context.CreateConsumer(subject)
		if err != nil {
			m.Error("Create solo consumer error: %v", err)
			return nil, err
		}
	}

	m.consumers[subject] = consumer

	return consumer.GetMessage(), nil
}

// Unsubscribe ...
func (m *Messager) Unsubscribe(subject string) error {
	consumer, ok := m.consumers[subject]
	if !ok {
		m.Info("Unsubscribe a consumer which does not exist: %v", subject)
		return nil
	}

	err := consumer.Close()
	if err != nil {
		m.Error("Close consumer[%v] error: %v", subject, err)
		return err
	}

	delete(m.consumers, subject)

	return nil
}
