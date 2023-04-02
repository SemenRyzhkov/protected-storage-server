package config

// Config конфигурация
type Config struct {
	Host            string
	Key             string
	DataBaseAddress string
}

// New конструктор конфига
func New(serverAddress, key, dbAddress string) Config {
	return Config{
		Host:            serverAddress,
		Key:             key,
		DataBaseAddress: dbAddress,
	}
}
