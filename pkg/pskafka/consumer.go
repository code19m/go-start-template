package pskafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

// NewSubscriber validates the configuration and returns a new subscriber.
func NewSubscriber(cfg *SubscriberConfig) (*Subscriber, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: subscriber config is nil", ErrInvalidSubscriberConfig)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	dialer := kafka.DefaultDialer

	// Set dialer SASL mechanism based on the configuration
	switch cfg.SecurityProtocol {
	case Plaintext:
		// No SASL mechanism needed
	case SaslPlaintext:
		mechanism, err := cfg.SaslPlaintextConfig.mechanism()
		if err != nil {
			return nil, err
		}
		dialer.SASLMechanism = mechanism
	case SaslScrum:
		mechanism, err := cfg.SaslScrumConfig.mechanism()
		if err != nil {
			return nil, err
		}
		dialer.SASLMechanism = mechanism
	default:
		return nil, fmt.Errorf("unsupported security protocol: %s", cfg.SecurityProtocol)
	}

	subscriber := &Subscriber{
		brokers:    cfg.Brokers,
		groupID:    cfg.GroupID,
		dialer:     dialer,
		shutdownCh: make(chan struct{}),
		doneCh:     make(chan struct{}),
	}

	return subscriber, nil
}

// Subscriber is an abstraction that groups multiple consumers.
// It provides a way to subscribe to multiple topics and consume messages.
// It also provides a way to add global and local interceptors.
type Subscriber struct {
	brokers []string
	groupID string
	dialer  *kafka.Dialer

	interceptors []InterceptorFunc
	consumers    []consumer

	shutdownCh chan struct{}
	doneCh     chan struct{}
}

// Subscribe adds a consumer to the subscriber with a topic and a handler.
func (s *Subscriber) Subscribe(topic string, handler HandleFunc) {
	s.consumers = append(s.consumers, consumer{
		topic:   topic,
		handler: handler,
	})
}

// SubscribeWithInterceptors adds a consumer to the subscriber with a topic, local interceptors, and a handler.
func (s *Subscriber) SubscribeWithInterceptors(topic string, interceptors []InterceptorFunc, handler HandleFunc) {
	s.consumers = append(s.consumers, consumer{
		topic:             topic,
		localInterceptors: interceptors,
		handler:           handler,
	})
}

// Use adds global interceptors to the subscriber that will be applied to all consumers
func (s *Subscriber) Use(interceptors ...InterceptorFunc) {
	s.interceptors = append(s.interceptors, interceptors...)
}

// Consume starts all consumers.
func (c *Subscriber) Consume() {
	var wg sync.WaitGroup

	for _, subscriber := range c.consumers {
		subscriber := subscriber

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.consume(subscriber)
		}()
	}

	wg.Wait()
	close(c.doneCh)
}

// Shutdown closes all consumers and waits for them to finish processing messages.
func (c *Subscriber) Shutdown(ctx context.Context) error {
	close(c.shutdownCh)

	select {
	case <-c.doneCh:
		fmt.Println("All consumers have been shutdown")
		return nil
	case <-ctx.Done():
		fmt.Println("Timeout waiting for consumers to shutdown")
		return ctx.Err()
	}
}

// consumer is an abstraction that groups a topic, a handler, and local interceptors.
type consumer struct {
	topic             string
	handler           HandleFunc
	localInterceptors []InterceptorFunc
}

// chainInterceptors chains global and local interceptors to the consumer's handler.
// It returns the final handler that will be used to consume messages.
func (s *Subscriber) chainInterceptors(subscriber consumer) HandleFunc {
	chain := subscriber.handler
	for i := len(subscriber.localInterceptors) - 1; i >= 0; i-- {
		chain = func(next HandleFunc, i InterceptorFunc) HandleFunc {
			return func(ctx context.Context, msg kafka.Message) error {
				return i(ctx, msg, func(ctx context.Context, msg kafka.Message) error {
					return next(ctx, msg)
				})
			}
		}(chain, subscriber.localInterceptors[i])
	}

	for i := len(s.interceptors) - 1; i >= 0; i-- {
		chain = func(next HandleFunc, i InterceptorFunc) HandleFunc {
			return func(ctx context.Context, msg kafka.Message) error {
				return i(ctx, msg, func(ctx context.Context, msg kafka.Message) error {
					return next(ctx, msg)
				})
			}
		}(chain, s.interceptors[i])
	}

	return chain
}

// consume reads messages from a topic and calls the consumer's handler.
// It stops reading messages when the subscriber is closed.
// It is not enable auto-commit, so the consumer is responsible
// for committing the messages with r.CommitMessages method after processing them
//
// TODO: Add support for auto-commit
// TODO: Add support for concurrent processing
// TODO: Add support for message batching
//
// For now it reads messages one by one and processes them synchronously
func (s *Subscriber) consume(subscriber consumer) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: s.brokers,
		Dialer:  s.dialer,
		GroupID: s.groupID,
		Topic:   subscriber.topic,
	})

	handler := s.chainInterceptors(subscriber)

	go func() {
		<-s.shutdownCh
		err := r.Close()
		if err != nil {
			fmt.Println("Error closing reader:", err)
		} else {
			fmt.Println("Reader closed")
		}
	}()

	for {
		ctx := context.Background()
		m, err := r.FetchMessage(ctx)
		if err != nil {
			return
		}

		if err := handler(ctx, m); err != nil {
			return
		}
	}
}
