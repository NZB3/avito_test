package authmiddleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/app/controllers"
	"project/internal/logger"
)

type authService interface {
	Authenticate(ctx context.Context, token string) (bool, error)
}

type middleware struct {
	log logger.Logger
	as  authService
}

func New(log logger.Logger, as authService) *middleware {
	return &middleware{
		log: log,
		as:  as,
	}
}

func (m *middleware) Auth() gin.HandlerFunc {
	const op = "authmiddleware.Auth"
	return func(ctx *gin.Context) {
		m.log.Info("Authorizing")

		tokenString := ctx.Request.Header.Get("token")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
			ctx.Abort()
			return
		}

		admin, err := m.as.Authenticate(ctx, tokenString)
		if err != nil {
			m.log.Errorf("%s Failed to authenticate: %v", op, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": controllers.BadRequest})
			ctx.Abort()
			return
		}

		ctx.Set("admin", admin)
		ctx.Next()
	}
}

func (m *middleware) AdminRequired() gin.HandlerFunc {
	const op = "middleware.AdminRequired"
	return func(ctx *gin.Context) {
		admin, ok := ctx.Get("admin")
		if !ok {
			m.log.Errorf("%s 'admin' required in context", op)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
			ctx.Abort()
			return
		}

		if !admin.(bool) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": Forbidden})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
