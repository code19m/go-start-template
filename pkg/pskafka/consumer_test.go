package pskafka_test

import (
	"context"
	"fmt"
	"go-start-template/pkg/pskafka"
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestConsumer(t *testing.T) {
	consumer := pskafka.NewConsumer([]string{"localhost:9092"}, "test-group")

	// consumer.Use(Interceptor1, Interceptor2, Interceptor3)

	consumer.Subscribe("test-topic", func(ctx context.Context, msg kafka.Message) error {
		fmt.Printf("Received message: %s\n", string(msg.Value))
		return nil
	})

	consumer.SubscribeWithInterceptors(
		"test-topic-2",
		// pskafka.InterceptorGroup{Interceptor1, Interceptor2, Interceptor3},
		nil,
		func(ctx context.Context, msg kafka.Message) error {
			fmt.Printf("Received message: %s\n", string(msg.Value))
			return nil
		},
	)

	fmt.Println("Consuming messages...")
	consumer.Consume()
}

func Interceptor1(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
	fmt.Println("Before Interceptor 1")
	err := next(ctx, msg)
	fmt.Println("After Interceptor 1")
	return err
}

func Interceptor2(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
	fmt.Println("Before Interceptor 2")
	err := next(ctx, msg)
	fmt.Println("After Interceptor 2")
	return err
}

func Interceptor3(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
	fmt.Println("Before Interceptor 3")
	err := next(ctx, msg)
	fmt.Println("After Interceptor 3")
	return err
}
