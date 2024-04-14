package cache

import (
	"errors"
	"os"
	"strconv"
)

type cacheConfig struct {
	port string
	host string
	db   int
}

func loadConfig() (*cacheConfig, error) {
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return nil, errors.New("REDIS_PORT environment variable not set")
	}

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return nil, errors.New("REDIS_HOST environment variable not set")
	}

	db := os.Getenv("REDIS_DB")
	if db == "" {
		return nil, errors.New("REDIS_DB environment variable not set")
	}

	dbInt, err := strconv.Atoi(db)
	if err != nil {
		return nil, errors.New("REDIS_DB environment variable not valid")
	}

	return &cacheConfig{
		port: port,
		host: host,
		db:   dbInt,
	}, nil
}
