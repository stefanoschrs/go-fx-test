package router

import (
	"github.com/stefanoschrs/go-fx-test/internal/web/controller"
	"github.com/stefanoschrs/go-fx-test/internal/web/webserver"

	"go.uber.org/fx"
)

func registerRoutes(webserver *webserver.Webserver, controller *controller.Controller) {
	webserver.Gin.GET("/ping", controller.Ping)
}

var Module = fx.Options(fx.Invoke(registerRoutes))
