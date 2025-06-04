package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/controller"
)

func Init(gin *gin.Engine, container container.Container) {

	// controllers
	setLinkController(gin, container)
}

func setLinkController(router *gin.Engine, container container.Container) {
	linkController := controller.NewLinkController(container)

	router.POST("/links", linkController.CreateLink)
	router.PUT("/links", linkController.UpdateLink)
	router.DELETE("/links/:id", linkController.DeleteLink)
}

/*
1. From where I should take repository?
*/
