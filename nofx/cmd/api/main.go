package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stefanoschrs/go-fx-test/nofx/internal/database"

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

	/**/
	// Logger
	/**/

	logger := zap.NewExample()

	/**/
	// Database
	/**/

	db, err := database.New(logger)
	if err != nil {
		log.Fatal(fmt.Errorf("database.New: %w", err))
	}

	/**/
	// Webserver
	/**/

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("database", db)
	})

	router.GET("/ping", func(c *gin.Context) {
		l := c.MustGet("logger").(*zap.Logger)
		l.Debug("Controller.Ping")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/something", func(c *gin.Context) {
		l := c.MustGet("logger").(*zap.Logger)
		l.Debug("Controller.GetSomething")

		d := c.MustGet("database").(*database.Database)

		number, err2 := d.GetRandomNumber()
		if err2 != nil {
			l.Error("database.GetRandomNumber", zap.Error(err2))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number": number,
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

	if err = router.Run(getWebserverAddr()); err != nil {
		log.Fatal(fmt.Errorf("router.Run: %w", err))
	}
}
