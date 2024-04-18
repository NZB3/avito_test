package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/redis/go-redis/v9"
	"project/internal/app/models"
	"project/internal/logger"
	"strconv"
	"time"
)

type cache struct {
	conn *redis.Client
	log  logger.Logger
}

func New(logger logger.Logger) (*cache, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	conn := redis.NewClient(&redis.Options{
		Addr: cfg.host + ":" + cfg.port,
		DB:   cfg.db,
	})

	return &cache{
		conn: conn,
		log:  logger,
	}, nil
}

func (c *cache) SetBanner(ctx context.Context, tagID int, featureID int, banner models.Banner) error {
	const op = "cache.SetBanner"
	h := hash(tagID, featureID)
	err := c.conn.Set(ctx, h, banner, time.Minute*5).Err()
	if err != nil {
		c.log.Errorf("%s Failed to set banner: %s", op, err)
		return err
	}

	return nil
}

func (c *cache) GetBanner(ctx context.Context, tagID int, featureID int) (models.Banner, error) {
	const op = "cache.GetBanner"

	h := hash(tagID, featureID)

	var banner models.Banner
	err := c.conn.Get(ctx, h).Scan(&banner)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return banner, models.BannerNotFound
		}

		c.log.Errorf("%s Failed to get banner from cache: %s", op, err)
		return banner, err
	}

	return banner, nil
}

func hash(args ...int) string {
	var s string
	for _, arg := range args {
		s += strconv.Itoa(arg)
	}
	h := md5.Sum([]byte(s))
	hashHex := hex.EncodeToString(h[:])

	return hashHex
}
