package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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
	// .env
	envPath := "../configs/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	// config.yaml
	v := viper.New()
	v.SetConfigFile("../configs/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Printf("warning: failed to read config.yaml: %v", err)
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
	brokers := []string{}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = append(brokers, envBrokers)
	} else {
		brokers = v.GetStringSlice("kafka.brokers")
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = v.GetString("kafka.topic")
	}
	partition := v.GetInt("kafka.partition")

	// логгер
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = v.GetString("logging.level")
	}
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
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

	initLogger(cfg.Logger.Level)

	return cfg, nil
}

func initLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(level)
	zerolog.SetGlobalLevel(level)
	_ = log
}