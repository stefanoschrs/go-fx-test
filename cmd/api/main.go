package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/stefanoschrs/go-fx-test/internal/web/controller"
	"github.com/stefanoschrs/go-fx-test/internal/web/router"
	"github.com/stefanoschrs/go-fx-test/internal/web/webserver"

	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func getWebserverAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// e.g 127.0.0.1:8080
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	return port
}

func main() {
	fmt.Println()
	fmt.Println("Starting fx application...")
	fmt.Println()

	app := fx.New(
		webserver.Module,
		router.Module,
		controller.Module,
		fx.Provide(zap.NewExample),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Invoke(func(webserver *webserver.Webserver, logger *zap.Logger) {
			logger.Debug("Webserver module invoked")
			go webserver.Gin.Run(getWebserverAddr())
		}, func(ctrl *controller.Controller, logger *zap.Logger) {
			logger.Debug("Controller module invoked")
		}, func(logger *zap.Logger) {
			logger.Debug("Logger module invoked")
		}),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(fmt.Errorf("app.Start: %w", err))
	}

	/**
	 * Testing if the webserver is running
	 */
	go func() {
		time.Sleep(5 * time.Second)

		res, err := resty.New().R().Get(fmt.Sprintf("http://%s/ping", getWebserverAddr()))
		if err != nil {
			log.Fatal(fmt.Errorf("resty.Get: %w", err))
		}
		fmt.Println("Testing Server: " + string(res.Body()))
	}()

	<-app.Wait()
}
