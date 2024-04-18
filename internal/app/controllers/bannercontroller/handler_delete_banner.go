package bannercontroller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
)

type bannerDeleter interface {
	DeleteBanner(ctx context.Context, bannerID int) (bool, error)
}

func (c *controller) DeleteHandler() gin.HandlerFunc {
	const op = "bannercontroller.DeleteBannerHandler"
	return func(ctx *gin.Context) {
		id, err := controllers.ParseQueryParam(ctx, "id", true, -1, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s: Failed to parse param: %v", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		ok, err := c.bs.DeleteBanner(ctx, id)
		if err != nil {
			c.log.Errorf("%s: Failed to delete banner: %v", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		if !ok {
			c.log.Errorf("%s: Not found banner: %d", op, id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": BannerNotFound})
		}

		ctx.JSON(http.StatusNoContent, gin.H{"message": BannerDeleted})
	}
}
