package config

import "fmt"

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	Driver   string
}

type serverConfig struct {
	Port string
}

// Config struct
type Config struct {
	Server  *serverConfig
	Db      *dbConfig
	Metrics *metrics
}

type metrics struct {
	Gateway string
}

//NewConfig new struct to configurations enviroments
func NewConfig() *Config {
	return &Config{
		Server: &serverConfig{
			Port: GetenvString("SERVER_PORT", "8080"),
		},
		Db: &dbConfig{
			Host:     GetenvString("DB_HOST", "127.0.0.1"),
			Port:     GetenvString("DB_PORT", "5432"),
			User:     GetenvString("DB_USER", "postgres"),
			Name:     GetenvString("DB_NAME", "user_service"),
			Password: GetenvString("DB_PASSWORD", "hash"),
			Driver:   GetenvString("DB_DRIVER", "postgres"),
		},
		Metrics: &metrics{
			Gateway: GetenvString("URL_PUSHGATEWAY", "http://localhost:9091"),
		},
	}
}

func (cfg *dbConfig) ToURL() string {
	url := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable "
	return fmt.Sprintf(url, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
