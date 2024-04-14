package bannercontroller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type userBannerGetter interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int, notCached bool) (*models.Banner, error)
}

func (c *controller) GetUserBannerHandler() gin.HandlerFunc {
	const op = "bannercontroller.GetUserBanner"
	return func(ctx *gin.Context) {
		tagID, err := controllers.ParseQueryParam(ctx.Params, "tag_id", true, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s : Failed to parse tag_id %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		featureID, err := controllers.ParseQueryParam(ctx.Params, "feature_id", true, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse feature_id: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		notCached, err := controllers.ParseQueryParam(ctx.Params, "use_last_revision", false, controllers.ConvToBool)
		if err != nil {
			c.log.Errorf("%s Failed to parse use_last_revision: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		banner, err := c.bs.GetUserBanner(ctx, tagID, featureID, notCached)
		if err != nil {
			c.log.Errorf("%s Failed to get banner: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if banner == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "banner not found"})
			return
		}

		ctx.JSON(http.StatusOK, banner.Content)
	}
}
