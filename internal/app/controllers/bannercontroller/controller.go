package bannercontroller

import (
	"project/internal/logger"
)

type bannerService interface {
	userBannerGetter
	bannerDeleter
	bannerSaver
	bannersGetter
	bannerUpdater
}

type controller struct {
	log logger.Logger
	bs  bannerService
}

func New(log logger.Logger, bs bannerService) *controller {
	return &controller{
		log: log,
		bs:  bs,
	}
}
