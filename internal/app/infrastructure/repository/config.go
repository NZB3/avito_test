package repository

import (
	"errors"
	"os"
)

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func loadConfig() (*dbConfig, error) {
	user := os.Getenv("DB_USER")
	if user == "" {
		return nil, errors.New("DB_USER environment variable not set")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, errors.New("DB_PASSWORD environment variable not set")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		return nil, errors.New("DB_NAME environment variable not set")
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, errors.New("DB_HOST environment variable not set")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return nil, errors.New("DB_PORT environment variable not set")
	}
	sslMode := os.Getenv("DB_SSL")
	if sslMode == "" {
		return nil, errors.New("DB_SSL environment variable not set")
	}
	return &dbConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
	}, nil
}
