package bannercontroller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/app/models"
)

type bannerUpdater interface {
	UpdateBanner(ctx context.Context, banner models.Banner) (bool, error)
}

type patchBannerRequest struct {
	FeatureID int    `json:"feature_id,omitempty"`
	TagIDs    []int  `json:"tag_ids,omitempty"`
	Content   []byte `json:"content,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
}

func (c *controller) PatchHandler() gin.HandlerFunc {
	const op = "bannercontroller.PatchBannerHandler"
	return func(ctx *gin.Context) {
		id, err := controllers.ParseQueryParam(ctx.Params, "id", true, controllers.ConvToInt)
		if err != nil {
			c.log.Errorf("%s : Failed to parse id: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		var req patchBannerRequest
		if err := ctx.ShouldBind(&req); err != nil {
			c.log.Errorf("%s : Failed to parse body: %s", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			return
		}

		banner := models.Banner{
			ID:       id,
			TagIDs:   req.TagIDs,
			Content:  req.Content,
			IsActive: req.IsActive,
		}

		ok, err := c.bs.UpdateBanner(ctx, banner)
		if err != nil {
			c.log.Errorf("%s : Failed to update banner: %s", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": controllers.InternalServerError})
			return
		}

		if !ok {
			c.log.Errorf("%s : Failed to update banner: %s", op, err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": BannerNotFound})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": controllers.OK})
	}
}
