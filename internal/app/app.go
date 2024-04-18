package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"project/internal/app/controllers/bannercontroller"
	"project/internal/app/controllers/middleware/authmiddleware"
	"project/internal/app/infrastructure/cache"
	"project/internal/app/infrastructure/repository"
	"project/internal/app/services/authservice"
	"project/internal/app/services/bannerservice"
	"project/internal/logger"
)

type app struct {
	log    logger.Logger
	server *http.Server
}

func New() *app {
	l := logger.New()
	port := os.Getenv("SERVER_PORT")
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}
	return &app{
		log:    l,
		server: srv,
	}
}

func (a *app) Run() error {
	repo, err := repository.New(a.log)
	if err != nil {
		return err
	}

	c, err := cache.New(a.log)
	if err != nil {
		return err
	}

	bannerService := bannerservice.New(a.log, repo, c)
	authService := authservice.New(a.log)

	authMiddleware := authmiddleware.New(a.log, authService)

	bannerController := bannercontroller.New(a.log, bannerService)

	router := gin.Default()

	userBannerRouter := router.Group("/user_banner")
	userBannerRouter.Use(authMiddleware.Auth())
	{
		userBannerRouter.GET("/", bannerController.GetUserBannerHandler())
	}

	bannerGroup := router.Group("/banner")
	bannerGroup.Use(authMiddleware.Auth(), authMiddleware.AdminRequired())
	{
		bannerGroup.GET("/", bannerController.GetHandler())
		bannerGroup.POST("/", bannerController.PostHandler())
		bannerGroup.PATCH("/:id", bannerController.PatchHandler())
		bannerGroup.DELETE("/:id", bannerController.DeleteHandler())
	}

	a.server.Handler = router

	return a.server.ListenAndServe()
}

func (a *app) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
