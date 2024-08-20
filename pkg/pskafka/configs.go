package pskafka

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	// Security protocols
	Plaintext     = "PLAINTEXT"
	SaslPlaintext = "SASL_PLAINTEXT"
	SaslScrum     = "SASL_SCRUM"

	// Scrum algorithms
	ScrumSHA256 = "SHA-256"
	ScrumSHA512 = "SHA-512"
)

// SubscriberConfig is the configuration for the subscriber.
type SubscriberConfig struct {

	// The list of broker addresses used to connect to the kafka cluster.
	Brokers []string `validate:"required,hostname_port"`

	// The security protocol used to communicate with the brokers.
	// Default is PLAINTEXT.
	SecurityProtocol string `validate:"required,oneof=PLAINTEXT SASL_PLAINTEXT SASL_SCRUM" default:"PLAINTEXT"`

	// The group ID of the consumer group.
	GroupID string `validate:"required"`

	// The configuration for SASL_PLAINTEXT security protocol.
	// Required if SecurityProtocol is SASL_PLAINTEXT
	SaslPlaintextConfig *SaslPlaintextConfig

	// The configuration for SASL_SCRUM security protocol.
	// Required if SecurityProtocol is SASL_SCRUM
	SaslScrumConfig *SaslScrumConfig

	// TODO: Implement other security protocols
}

func (c *SubscriberConfig) validate() error {
	v := validator.New()
	err := v.Struct(c)

	failedFields := make([]string, 0)
	if errs, ok := err.(validator.ValidationErrors); ok { //nolint: errorlint
		for _, err := range errs {
			failedFields = append(failedFields, err.Field())
		}
	}

	if len(failedFields) > 0 {
		return fmt.Errorf("%w: failed_keys: %v", ErrInvalidSubscriberConfig, failedFields)
	}

	if c.SecurityProtocol == SaslPlaintext && c.SaslPlaintextConfig == nil {
		return fmt.Errorf("SaslPlaintextConfig is required for SASL_PLAINTEXT security protocol")
	}

	if c.SecurityProtocol == SaslScrum && c.SaslScrumConfig == nil {
		return fmt.Errorf("SaslScrumConfig is required for SASL_SCRUM security protocol")
	}

	return nil
}

// SaslPlaintextConfig is the configuration for SASL_PLAINTEXT security protocol.
type SaslPlaintextConfig struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (c *SaslPlaintextConfig) mechanism() (sasl.Mechanism, error) {
	return plain.Mechanism{
		Username: c.Username,
		Password: c.Password,
	}, nil
}

// SaslScrumConfig is the configuration for SASL_SCRUM security protocol.
type SaslScrumConfig struct {
	Algorithm string `validate:"required,oneof=SHA-256 SHA-512"`
	Username  string `validate:"required"`
	Password  string `validate:"required"`
}

func (c *SaslScrumConfig) mechanism() (sasl.Mechanism, error) {
	switch c.Algorithm {
	case ScrumSHA256:
		return scram.Mechanism(scram.SHA256, c.Username, c.Password)
	case ScrumSHA512:
		return scram.Mechanism(scram.SHA512, c.Username, c.Password)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", c.Algorithm)
	}
}
