package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SpicyChickenFLY/gin-torrent/controller"
	"github.com/SpicyChickenFLY/gin-torrent/pkgs/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.Cors())
	router.Static("/static", "static")
	router.LoadHTMLGlob("templates/*")

	groupWeb := router.Group("")
	{
		groupTorrent := groupWeb.Group("/torrent")
		{
			groupTorrent.GET("/task", controller.taskUI)
		}
	}

	groupAPI := router.Group("/api")
	{
		groupTorrent := groupAPI.Group("/torrent")
		{
			groupTorrent.GET("/task", controller.RetriveAllTorrents)
			groupTorrent.GET("/task/:id", controller.RetriveTorrent)
			groupTorrent.POST("/task", controller.CreateTorrent)
			groupTorrent.POST("/task:id", controller.DeleteTorrent)
		}
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}
	log.Println("Server exiting")
}
