package main

import (
	"context"
	"eureka/fargoclient"
	"eureka/interfaces"
	"eureka/vars"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(vars.LocalPort),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	var ec interfaces.EurekaClient
	ec = fargoclient.NewClient()
	// ec = restclient.NewClient()
	// ec = xuanboclient.NewClient()
	// ec = fairwayclient.NewClient()
	// ec = springbootclient.NewClient()
	// ec = hikoqiuclient.NewClient()
	// ec = arthurhltclient.NewClient()
	ec.Register()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ec.DeRegister()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
