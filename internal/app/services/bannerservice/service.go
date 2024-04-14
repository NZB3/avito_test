package bannerservice

import (
	"context"
	"project/internal/app/models"
	"project/internal/logger"
)

type bannerStorage interface {
	GetBanner(ctx context.Context, tagID int, featureID int) (*models.Banner, error)
	GetBannersByTag(ctx context.Context, tagID int, limit int, offset int) ([]*models.Banner, error)
	GetBannersByFeatureID(ctx context.Context, featureID int, limit int, offset int) ([]*models.Banner, error)
	UpdateBanner(ctx context.Context, banner models.Banner) (bool, error)
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	DeleteBanner(ctx context.Context, bannerID int) (bool, error)
}

type bannerCache interface {
	SetBanner(ctx context.Context, tagID int, featureID int, banner *models.Banner) error
	GetBanner(ctx context.Context, tagID int, featureID int) (*models.Banner, error)
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

func (s *service) GetUserBanner(ctx context.Context, tagID int, featureID int, notCached bool) (*models.Banner, error) {
	const op = "bannerservice.GetUserBanner"
	var banner *models.Banner
	var err error

	if notCached == false {
		banner, err = s.cache.GetBanner(ctx, tagID, featureID)
		if err != nil {
			s.log.Errorf("%s Failed to get Banner from cache: %s", op, err)
			return nil, err
		}
	}

	if banner == nil {
		banner, err = s.storage.GetBanner(ctx, tagID, featureID)
		if err != nil {
			s.log.Errorf("%s Failed to get Banner from storage: %s", op, err)
			return nil, err
		}

		err = s.cache.SetBanner(ctx, tagID, featureID, banner)
	}

	return banner, nil
}

func (s *service) GetBanners(ctx context.Context, tagID int, featureID int, limit int, offset int) ([]*models.Banner, error) {
	const op = "bannerservice.GetBanners"
	if tagID != -1 && featureID != -1 {
		banner, err := s.storage.GetBanner(ctx, tagID, featureID)
		if err != nil {
			s.log.Errorf("%s: Failed to get banner for tag %d and feature %d: %v", op, tagID, featureID, err)
			return nil, err
		}
		return []*models.Banner{banner}, nil
	}

	if tagID != -1 {
		banners, err := s.storage.GetBannersByTag(ctx, tagID, limit, offset)
		if err != nil {
			s.log.Errorf("%s: Failed to get banners for tag %d: %v", op, tagID, err)
			return nil, err
		}
		return banners, nil
	}

	banners, err := s.storage.GetBannersByFeatureID(ctx, featureID, limit, offset)
	if err != nil {
		s.log.Errorf("%s: Failed to get banners for feature %d: %v", op, featureID, err)
		return nil, err
	}
	return banners, nil
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
