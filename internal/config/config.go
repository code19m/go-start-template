package config

import "time"

type Config struct {
	AppMode    string     `env:"APP_MODE"`
	Project    Project    `yaml:"project"`
	Logger     Logger     `yaml:"logger"`
	HttpServer HttpServer `yaml:"http_server"`
	// Auth       Auth
	Postgres Postgres
	// Mongo      Mongo
}

type Project struct {
	Name    string `yaml:"name"    validate:"required"`
	Domain  string `env:"DOMAIN"   validate:"required"`
	Version string `yaml:"version" validate:"required"`
}

type Logger struct {
	Level  string `yaml:"level"  validate:"required,oneof=debug info warn error"`
	Format string `yaml:"format" validate:"required,oneof=text json"`
}

type HttpServer struct {
	TimeOut         time.Duration `yaml:"timeout"           validate:"required"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"      validate:"required"`
	MaxShutdownTime time.Duration `yaml:"max_shutdown_time" validate:"required"`
}

type Auth struct {
	Host         string `env:"AUTH_HOST"          validate:"required"`
	Port         int32  `env:"AUTH_PORT"          validate:"required"`
	InternalUser string `env:"AUTH_INTERNAL_USER" validate:"required"`
	InternalPass string `env:"AUTH_INTERNAL_PASS" validate:"required"`
	UseTLS       bool   `env:"AUTH_USE_TLS"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"     validate:"required"`
	Port     int32  `env:"POSTGRES_PORT"     validate:"required"`
	Db       string `env:"POSTGRES_DB"       validate:"required"`
	User     string `env:"POSTGRES_USER"     validate:"required"`
	Password string `env:"POSTGRES_PASSWORD" validate:"required"`
}

type Mongo struct {
	Host     string `env:"MONGO_HOST"     validate:"required"`
	Port     int32  `env:"MONGO_PORT"     validate:"required"`
	Db       string `env:"MONGO_DB"       validate:"required"`
	User     string `env:"MONGO_USER"     validate:"required"`
	Password string `env:"MONGO_PASSWORD" validate:"required"`
}
