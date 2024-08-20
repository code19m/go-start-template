package pskafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// HandleFunc is a function that handles a message.
type HandleFunc func(ctx context.Context, msg kafka.Message) error

// InterceptorFunc is a function that intercepts a message.
type InterceptorFunc func(ctx context.Context, msg kafka.Message, next HandleFunc) error
