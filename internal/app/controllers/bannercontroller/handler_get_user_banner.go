package bannercontroller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type userBannerGetter interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool) (models.Banner, error)
}

func (c *controller) GetUserBannerHandler() gin.HandlerFunc {
	const op = "bannercontroller.GetUserBanner"
	return func(ctx *gin.Context) {
		tagID, err := controllers.ParseQueryParam(ctx, "tag_id", true, -1, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s : Failed to parse tag_id %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		featureID, err := controllers.ParseQueryParam(ctx, "feature_id", true, -1, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s Failed to parse feature_id: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		useLastRevision, err := controllers.ParseQueryParam(ctx, "use_last_revision", false, false, controllers.ConvToBool)
		if err != nil {
			c.log.Errorf("%s Failed to parse use_last_revision: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		banner, err := c.bs.GetUserBanner(ctx, tagID, featureID, useLastRevision)
		if errors.Is(err, models.BannerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": BannerNotFound})
			return
		}

		if err != nil {
			c.log.Errorf("%s Failed to get banner: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		admin, err := controllers.CheckAdminStatus(ctx)
		if err != nil {
			c.log.Errorf("%s : Failed to check admin status: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		if admin {
			ctx.IndentedJSON(http.StatusOK, banner)
		} else {
			ctx.IndentedJSON(http.StatusOK, banner.Content)
		}
	}
}
