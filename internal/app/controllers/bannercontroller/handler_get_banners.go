package bannercontroller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type bannersGetter interface {
	GetBanners(ctx context.Context, featureID int, tagID int, limit int, offset int) ([]*models.Banner, error)
}

func (c *controller) GetHandler() gin.HandlerFunc {
	const op = "bannercontroller.GetHandler"
	return func(ctx *gin.Context) {
		featureID, err := controllers.ParseQueryParam(ctx.Params, "feature_id", false, controllers.ConvToInt)
		if err != nil {
			if !errors.Is(err, controllers.ParameterIsNotSpecified) {
				c.log.Errorf("%s Failed to parse params: %s", op, err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
				return
			}

			featureID = -1
		}

		tagID, err := controllers.ParseQueryParam(ctx.Params, "tag_id", false, controllers.ConvToInt)
		if err != nil {
			if !errors.Is(err, controllers.ParameterIsNotSpecified) {
				c.log.Errorf("%s Failed to parse params: %s", op, err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
				return
			}

			tagID = -1
		}

		limit, err := controllers.ParseQueryParam(ctx.Params, "limit", false, controllers.ConvToInt)
		if err != nil {
			if !errors.Is(err, controllers.ParameterIsNotSpecified) {
				c.log.Errorf("%s Failed to parse params: %s", op, err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
				return
			}

			limit = 10
		}

		offset, err := controllers.ParseQueryParam(ctx.Params, "offset", false, controllers.ConvToInt)
		if err != nil {
			if !errors.Is(err, controllers.ParameterIsNotSpecified) {
				c.log.Errorf("%s Failed to parse params: %s", op, err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
				return
			}

			offset = 0
		}

		banners, err := c.bs.GetBanners(ctx, featureID, tagID, limit, offset)
		if err != nil {
			c.log.Errorf("%s Failed to get banners: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		bannersJSON, err := json.Marshal(banners)
		if err != nil {
			c.log.Errorf("%s Failed to marshal banners: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, bannersJSON)
	}
}
