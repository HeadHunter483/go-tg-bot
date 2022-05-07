package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type WebConfig struct {
	Host string
	Port string
}

type Config struct {
	Token  string
	Domain string
	Web    WebConfig
	DB     DBConfig
}

// New creates a new Config instance.
func New() *Config {
	return &Config{
		Token:  os.Getenv("TOKEN"),
		Domain: os.Getenv("DOMAIN"),
		Web: WebConfig{
			Host: os.Getenv("WEB_HOST"),
			Port: os.Getenv("WEB_PORT"),
		},
		DB: DBConfig{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
	}
}

// GetDBInfo forms a database connection configuration string.
func (conf *Config) GetDBInfo() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Pass, conf.DB.Name)
}
