package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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
	fmt.Println("Starting non-fx application...")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Println()

	logger := zap.NewExample()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
	})

	router.GET("/ping", func(c *gin.Context) {
		l := c.MustGet("logger").(*zap.Logger)
		l.Debug("Controller.Ping")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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

	if err := router.Run(getWebserverAddr()); err != nil {
		log.Fatal(fmt.Errorf("router.Run: %w", err))
	}
}
