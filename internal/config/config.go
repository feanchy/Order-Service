package config

import "os"

type Config struct {
	HTTPPort    string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	KafkaBroker string
	KafkaTopic  string
}

func Load() *Config {
	return &Config{
		HTTPPort:    os.Getenv("APP_PORT"),
		DBHost:      os.Getenv("POSTGRES_HOST"),
		DBPort:      os.Getenv("POSTGRES_PORT"),
		DBUser:      os.Getenv("POSTGRES_USER"),
		DBPassword:  os.Getenv("POSTGRES_PASSWORD"),
		DBName:      os.Getenv("POSTGRES_DB"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		KafkaTopic:  os.Getenv("KAFKA_TOPIC"),
	}
}
