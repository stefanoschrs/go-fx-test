package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stefanoschrs/go-fx-test/internal/web/controller"
	"github.com/stefanoschrs/go-fx-test/internal/web/router"
	"github.com/stefanoschrs/go-fx-test/internal/web/webserver"

	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

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
			go webserver.Gin.Run()
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

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		res, err := resty.New().R().Get(fmt.Sprintf("http://localhost:%s/ping", port))
		if err != nil {
			log.Fatal(fmt.Errorf("resty.Get: %w", err))
		}
		fmt.Println("Testing Server: " + string(res.Body()))
	}()

	<-app.Wait()
}
