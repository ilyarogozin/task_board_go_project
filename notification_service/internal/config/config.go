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
	Database DatabaseConfig
	Kafka    KafkaConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type LoggerConfig struct {
	Level zerolog.Level
}

func LoadConfig() (*Config, error) {
	// базовая инициализация логгера (до чтения конфига)
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
		Database: DatabaseConfig{
			DSN:          dsn,
			MaxOpenConns: maxOpen,
			MaxIdleConns: maxIdle,
		},
		Kafka: KafkaConfig{
			Brokers: brokers,
			Topic:   topic,
		},
		Logger: LoggerConfig{
			Level: level,
		},
	}

	// финальная инициализация логгера с уровнем из конфига
	initLogger(cfg.Logger.Level)

	log.Info().
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