package solo

import (
	"github.com/ewangplay/eventbus/i"
)

// Producer represents Solo Mode Messager Producer
type Producer struct {
	*Context
}

// NewProducer creates a new Solo Mode Messager Producer instance
func NewProducer(ctx *Context) (*Producer, error) {
	p := &Producer{}

	p.Context = ctx

	return p, nil
}

// Close closes the current Producer instance
func (p *Producer) Close() error {
	return nil
}

// Publish ...
func (p *Producer) Publish(msg i.Message) (err error) {
	subject := msg.GetSubject()
	p.Lock()
	queue, found := p.Queues[subject]
	if !found {
		queue = make(chan []byte, QueueMaxSize)
		p.Queues[subject] = queue
	}
	p.Unlock()

	queue <- msg.GetData()

	return nil
}
