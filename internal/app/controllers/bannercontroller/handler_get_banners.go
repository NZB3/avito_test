package bannercontroller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type bannersGetter interface {
	GetBanners(ctx context.Context, featureID int, tagID int, limit int, offset int) ([]models.Banner, error)
}

func (c *controller) GetHandler() gin.HandlerFunc {
	const op = "bannercontroller.GetHandler"
	return func(ctx *gin.Context) {
		featureID, err := controllers.ParseQueryParam(ctx, "feature_id", false, -1, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse params: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		tagID, err := controllers.ParseQueryParam(ctx, "tag_id", false, -1, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse params: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		limit, err := controllers.ParseQueryParam(ctx, "limit", false, 10, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse params: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		offset, err := controllers.ParseQueryParam(ctx, "offset", false, 0, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse params: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		banners, err := c.bs.GetBanners(ctx, featureID, tagID, limit, offset)
		if err != nil {
			c.log.Errorf("%s Failed to get banners: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		ctx.IndentedJSON(http.StatusOK, &banners)
	}
}
