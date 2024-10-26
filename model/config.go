package model

type Config struct {
	apiUrl    string
	processes []string
	mode      bool
}

type ConfigRepository interface {
	Set(config Config) *Config
	Get() *Config
}
