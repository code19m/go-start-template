package pskafka_test

// import (
// 	"context"
// 	"fmt"
// 	"go-start-template/pkg/pskafka"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"testing"
// 	"time"

// 	"github.com/segmentio/kafka-go"
// )

// func TestConsumer(t *testing.T) {
// 	consumer := pskafka.NewSubscriber([]string{"localhost:9092"}, "test-group")

// 	// consumer.Use(Interceptor1, Interceptor2, Interceptor3)

// 	consumer.Subscribe("test-topic", func(ctx context.Context, msg kafka.Message) error {
// 		fmt.Printf("Received message: %s\n", string(msg.Value))
// 		time.Sleep(15 * time.Second)
// 		fmt.Println("Message processed after 15 seconds")
// 		return nil
// 	})

// 	consumer.SubscribeWithInterceptors(
// 		"test-topic-2",
// 		// pskafka.InterceptorGroup{Interceptor1, Interceptor2, Interceptor3},
// 		nil,
// 		func(ctx context.Context, msg kafka.Message) error {
// 			fmt.Printf("Received message: %s\n", string(msg.Value))
// 			return nil
// 		},
// 	)

// 	fmt.Println("Consuming messages...")
// 	go consumer.Consume()

// 	// Graceful Shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
// 	<-quit

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
// 	defer cancel()
// 	err := consumer.Shutdown(ctx)
// 	if err != nil {
// 		fmt.Println("Error shutting down consumer:", err)
// 	}
// }

// func Interceptor1(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
// 	fmt.Println("Before Interceptor 1")
// 	err := next(ctx, msg)
// 	fmt.Println("After Interceptor 1")
// 	return err
// }

// func Interceptor2(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
// 	fmt.Println("Before Interceptor 2")
// 	err := next(ctx, msg)
// 	fmt.Println("After Interceptor 2")
// 	return err
// }

// func Interceptor3(ctx context.Context, msg kafka.Message, next pskafka.HandleFunc) error {
// 	fmt.Println("Before Interceptor 3")
// 	err := next(ctx, msg)
// 	fmt.Println("After Interceptor 3")
// 	return err
// }
