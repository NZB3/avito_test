package bannerservice

import (
	"context"
	"errors"
	"project/internal/app/models"
	"project/internal/logger"
)

type bannerStorage interface {
	GetBanner(ctx context.Context, tagID int, featureID int) (models.Banner, error)
	GetBannersByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Banner, error)
	GetBannersByFeatureID(ctx context.Context, featureID int, limit int, offset int) ([]models.Banner, error)
	GetBanners(ctx context.Context, limit int, offset int) ([]models.Banner, error)
	UpdateBanner(ctx context.Context, banner models.Banner) (bool, error)
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	DeleteBanner(ctx context.Context, bannerID int) (bool, error)
}

type bannerCache interface {
	SetBanner(ctx context.Context, tagID int, featureID int, banner models.Banner) error
	GetBanner(ctx context.Context, tagID int, featureID int) (models.Banner, error)
}

type service struct {
	log     logger.Logger
	storage bannerStorage
	cache   bannerCache
}

func New(log logger.Logger, storage bannerStorage, cache bannerCache) *service {
	return &service{
		log:     log,
		storage: storage,
		cache:   cache,
	}
}

func (s *service) GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool) (models.Banner, error) {
	const op = "bannerservice.GetUserBanner"
	if !useLastRevision {
		cachedBanner, err := s.cache.GetBanner(ctx, tagID, featureID)
		if err == nil {
			return cachedBanner, nil
		}
		if !errors.Is(err, models.BannerNotFound) {
			s.log.Errorf("%s Failed to get Banner from cache: %s", op, err)
			return models.Banner{}, err
		}
	}

	storageBanner, err := s.storage.GetBanner(ctx, tagID, featureID)
	if err != nil {
		s.log.Errorf("%s Failed to get Banner from storage: %s", op, err)
		return models.Banner{}, err
	}

	err = s.cache.SetBanner(ctx, tagID, featureID, storageBanner)
	if err != nil {
		s.log.Errorf("%s Failed to set Banner in cache: %s", op, err)
	}

	return storageBanner, nil
}

func (s *service) GetBanners(ctx context.Context, tagID int, featureID int, limit int, offset int) ([]models.Banner, error) {
	const op = "bannerservice.GetBanners"

	if tagID != -1 && featureID != -1 {
		return s.getBannerByTagAndFeatureID(ctx, op, tagID, featureID)
	}

	if tagID != -1 {
		return s.getBannersByTagID(ctx, op, tagID, limit, offset)
	}

	if featureID != -1 {
		return s.getBannersByFeatureID(ctx, op, featureID, limit, offset)
	}

	return s.getBanners(ctx, op, limit, offset)
}

func (s *service) UpdateBanner(ctx context.Context, banner models.Banner) (ok bool, err error) {
	const op = "bannerservice.UpdateBanner"
	ok, err = s.storage.UpdateBanner(ctx, banner)
	if err != nil {
		s.log.Errorf("%s Failed to update banner %d: %v", op, banner.ID, err)
		return false, err
	}

	return ok, nil
}

func (s *service) SaveBanner(ctx context.Context, banner models.Banner) (int, error) {
	const op = "bannerservice.SaveBanner"
	id, err := s.storage.CreateBanner(ctx, banner)
	if err != nil {
		s.log.Errorf("%s Failed to create banner %d: %v", op, banner.ID, err)
		return 0, err
	}

	return id, nil
}

func (s *service) DeleteBanner(ctx context.Context, id int) (ok bool, err error) {
	const op = "bannerservice.DeleteBanner"
	ok, err = s.storage.DeleteBanner(ctx, id)
	if err != nil {
		s.log.Errorf("%s Failed to delete banner %d: %v", op, id, err)
		return false, err
	}

	return ok, nil
}

func (s *service) getBannerByTagAndFeatureID(ctx context.Context, op string, tagID, featureID int) ([]models.Banner, error) {
	banner, err := s.storage.GetBanner(ctx, tagID, featureID)
	if err != nil {
		s.log.Errorf("%s: Failed to get banner for tag %d and feature %d: %v", op, tagID, featureID, err)
		return nil, err
	}

	return []models.Banner{banner}, nil
}

func (s *service) getBannersByTagID(ctx context.Context, op string, tagID, limit, offset int) ([]models.Banner, error) {
	banners, err := s.storage.GetBannersByTagID(ctx, tagID, limit, offset)
	if err != nil {
		s.log.Errorf("%s: Failed to get banners for tag %d: %v", op, tagID, err)
		return nil, err
	}
	return banners, nil
}

func (s *service) getBannersByFeatureID(ctx context.Context, op string, featureID, limit, offset int) ([]models.Banner, error) {
	banners, err := s.storage.GetBannersByFeatureID(ctx, featureID, limit, offset)
	if err != nil {
		s.log.Errorf("%s: Failed to get banners for feature %d: %v", op, featureID, err)
		return nil, err
	}
	return banners, nil
}

func (s *service) getBanners(ctx context.Context, op string, limit, offset int) ([]models.Banner, error) {
	banners, err := s.storage.GetBanners(ctx, limit, offset)
	if err != nil {
		s.log.Errorf("%s: Failed to get banners for limit %d, offset %d: %v", op, limit, offset, err)
		return nil, err
	}

	return banners, nil
}
