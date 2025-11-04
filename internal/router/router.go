package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/controller"
	"github.com/pewpowder/url-shortener/internal/middleware"
	"github.com/pewpowder/url-shortener/internal/service"
)

func Init(gin *gin.Engine, container container.Container) {
	config := container.GetConfig()

	// third-party services
	authService := service.NewAuthService(config.Server.AuthServiceURL)

	// middlewares
	authMiddleware := middleware.AuthMiddleware(authService)

	// controllers
	setLinkController(gin, container, authMiddleware)
	setTagController(gin, container, authMiddleware)
}

func setLinkController(router *gin.Engine, container container.Container, authMiddleware gin.HandlerFunc) {
	linkController := controller.NewLinkController(container)

	router.GET("/links", authMiddleware, linkController.GetLinks)
	router.GET("/links/:id", authMiddleware, linkController.GetLinkDetails)
	router.POST("/links", authMiddleware, linkController.CreateLink)
	router.PUT("/links", authMiddleware, linkController.UpdateLink)
	router.DELETE("/links/:id", authMiddleware, linkController.DeleteLink)
}

func setTagController(router *gin.Engine, container container.Container, authMiddleware gin.HandlerFunc) {
	tagController := controller.NewTagController(container)

	router.GET("/tags", authMiddleware, tagController.GetTags)
	router.GET("/tags/:id", authMiddleware, tagController.GetTagByID)
	router.POST("/tags", authMiddleware, tagController.CreateTag)
	router.PUT("/tags/:id", authMiddleware, tagController.UpdateTag)
	router.PATCH("/tags/:id", authMiddleware, tagController.PatchTag)
	router.DELETE("/tags/:id", authMiddleware, tagController.DeleteTag)
}
