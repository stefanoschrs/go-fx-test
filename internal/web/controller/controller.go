package controller

import (
	"net/http"

	"github.com/stefanoschrs/go-fx-test/internal/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Controller struct {
	logger   *zap.Logger
	database *database.Database
}

func (ctrl *Controller) Ping(c *gin.Context) {
	ctrl.logger.Debug("Controller.Ping")

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (ctrl *Controller) GetSomething(c *gin.Context) {
	ctrl.logger.Debug("Controller.GetSomething")

	number, err := ctrl.database.GetRandomNumber()
	if err != nil {
		ctrl.logger.Error("database.GetRandomNumber", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number": number,
	})
}

func New(logger *zap.Logger, database *database.Database) *Controller {
	var controller = new(Controller)
	controller.logger = logger
	controller.database = database

	return controller
}

var Module = fx.Provide(New)
