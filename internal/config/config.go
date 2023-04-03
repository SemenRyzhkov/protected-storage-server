package config

import "time"

// Config конфигурация
type Config struct {
	Host            string
	Key             string
	DataBaseAddress string
	TokenDuration   time.Duration
}

// New конструктор конфига
func New(serverAddress, key, dbAddress string) Config {
	return Config{
		Host:            serverAddress,
		Key:             key,
		DataBaseAddress: dbAddress,
		TokenDuration:   10 * time.Minute,
	}
}
