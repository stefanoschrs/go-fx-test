package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func (ctrl *Controller) Ping(c *gin.Context) {
	ctrl.logger.Debug("Controller.Ping")

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func New(logger *zap.Logger) *Controller {
	var controller = new(Controller)
	controller.logger = logger

	return controller
}

var Module = fx.Provide(New)
