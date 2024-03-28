package pskafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func NewConsumer(brokers []string, groupID string) *Consumer {
	return &Consumer{
		brokers:  brokers,
		groupID:  groupID,
		shutdown: make(chan struct{}),
	}
}

type Consumer struct {
	brokers []string

	interceptors []InterceptorFunc

	subscribers []subscriber

	groupID string

	shutdown chan struct{}
}

type subscriber struct {
	topic             string
	handler           HandleFunc
	localInterceptors []InterceptorFunc
}

func (c *Consumer) Subscribe(topic string, handler HandleFunc) {
	c.subscribers = append(c.subscribers, subscriber{
		topic:   topic,
		handler: handler,
	})

	log.Printf("Subscribed to topic %s", topic)
}

func (c *Consumer) SubscribeWithInterceptors(topic string, interceptors []InterceptorFunc, handler HandleFunc) {
	c.subscribers = append(c.subscribers, subscriber{
		topic:             topic,
		handler:           handler,
		localInterceptors: interceptors,
	})

	log.Printf("Subscribed to topic %s with interceptors", topic)
}

func (c *Consumer) Use(interceptors ...InterceptorFunc) {
	c.interceptors = append(c.interceptors, interceptors...)
}

func (c *Consumer) chainInterceptors(subscriber subscriber) HandleFunc {
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

	for i := len(c.interceptors) - 1; i >= 0; i-- {
		chain = func(next HandleFunc, i InterceptorFunc) HandleFunc {
			return func(ctx context.Context, msg kafka.Message) error {
				return i(ctx, msg, func(ctx context.Context, msg kafka.Message) error {
					return next(ctx, msg)
				})
			}
		}(chain, c.interceptors[i])
	}

	return chain
}

// HandleFunc is a function that handles a message.
type HandleFunc func(ctx context.Context, msg kafka.Message) error

// InterceptorFunc is a function that intercepts a message.
type InterceptorFunc func(ctx context.Context, msg kafka.Message, next HandleFunc) error

func (c *Consumer) consume(ctx context.Context, dialer *kafka.Dialer, subscriber subscriber) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		GroupID: c.groupID,
		Topic:   subscriber.topic,
		Dialer:  dialer,
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Error closing reader: %v", err)
		}
	}()

	for {
		select {
		case <-c.shutdown:
			fmt.Println("Shutting down consumer")
			return
		case <-ctx.Done():
			fmt.Println("Context done")
			return
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				return
			}

			chain := c.chainInterceptors(subscriber)
			if err := chain(ctx, m); err != nil {
				return
			}
		}
	}
}

func (c *Consumer) Consume() {
	dialer := kafka.DefaultDialer
	mechanism := plain.Mechanism{
		Username: "sharedKraftBrokerUser",
		Password: "/x2i9'0k!8Y3",
	}
	dialer.SASLMechanism = mechanism

	var wg sync.WaitGroup

	for _, subscriber := range c.subscribers {
		subscriber := subscriber
		ctx := context.Background()

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.consume(ctx, dialer, subscriber)
		}()
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = ctx

	c.Close()
	wg.Wait()
}

func (c *Consumer) Close() {
	close(c.shutdown)
}
