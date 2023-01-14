package webserver

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Webserver struct {
	Gin *gin.Engine
}

func New() *Webserver {
	g := gin.New()
	g.SetTrustedProxies(nil)

	var webserver = new(Webserver)
	webserver.Gin = g

	return webserver
}

var Module = fx.Options(fx.Provide(New))
