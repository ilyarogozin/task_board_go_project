package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server   GRPCServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	Logger   LoggerConfig
}

type GRPCServerConfig struct {
	GRPCPort string
}

type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type KafkaConfig struct {
	Brokers   []string
	Topic     string
	Partition int
}

type LoggerConfig struct {
	Level zerolog.Level
}

func LoadConfig() (*Config, error) {
	// инициализируем базовый логгер ДО всего
	initLogger(zerolog.InfoLevel)

	// .env
	envPath := "../configs/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Warn().
			Err(err).
			Str("path", envPath).
			Msg("failed to load .env file")
	}

	// config.yaml
	v := viper.New()
	v.SetConfigFile("../configs/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Warn().
			Err(err).
			Msg("failed to read config.yaml")
	}

	// gRPC Server
	grpcPort := v.GetString("server.grpc_port")
	if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
		grpcPort = envPort
	}

	// Database
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = v.GetString("database.dsn")
	}

	maxOpen := v.GetInt("database.max_open_conns")
	maxIdle := v.GetInt("database.max_idle_conns")

	// Kafka
	var brokers []string
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = []string{envBrokers}
	} else {
		brokers = v.GetStringSlice("kafka.brokers")
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = v.GetString("kafka.topic")
	}

	partition := v.GetInt("kafka.partition")

	// Logger level
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = v.GetString("logging.level")
	}

	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		log.Warn().
			Str("value", logLevelStr).
			Msg("invalid log level, fallback to info")
		level = zerolog.InfoLevel
	}

	cfg := &Config{
		Server: GRPCServerConfig{
			GRPCPort: grpcPort,
		},
		Database: DatabaseConfig{
			DSN:          dsn,
			MaxOpenConns: maxOpen,
			MaxIdleConns: maxIdle,
		},
		Kafka: KafkaConfig{
			Brokers:   brokers,
			Topic:     topic,
			Partition: partition,
		},
		Logger: LoggerConfig{
			Level: level,
		},
	}

	// переинициализируем логгер уже с нужным уровнем
	initLogger(cfg.Logger.Level)

	log.Info().
		Str("grpc_port", grpcPort).
		Str("kafka_topic", topic).
		Msg("configuration loaded successfully")

	return cfg, nil
}

func initLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(level)
}