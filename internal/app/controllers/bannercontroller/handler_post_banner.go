package bannercontroller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type bannerSaver interface {
	SaveBanner(ctx context.Context, banner models.Banner) (int, error)
}

type postBannerRequest struct {
	FeatureID int            `json:"feature_id"`
	TagIDs    []int          `json:"tag_ids"`
	Content   map[string]any `json:"content"`
	IsActive  bool           `json:"is_active"`
}

func (c *controller) PostHandler() gin.HandlerFunc {
	const op = "bannerHandler.PostBannerHandler"
	return func(ctx *gin.Context) {
		var req postBannerRequest
		if err := ctx.ShouldBind(&req); err != nil {
			c.log.Errorf("%s : Failed to parse body: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		id, err := c.bs.SaveBanner(ctx, models.Banner{
			TagIDs:    req.TagIDs,
			FeatureID: req.FeatureID,
			Content:   req.Content,
			IsActive:  req.IsActive,
		})

		if err != nil {
			c.log.Errorf("%s : Failed to save banner: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"banner_id": id})
	}
}
