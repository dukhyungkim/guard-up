package main

import (
	"bookman/config"
	"bookman/repository"
	v1 "bookman/router/v1"
	"bookman/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	opts, err := config.ParseFlags()
	if err != nil {
		log.Panicln(err)
	}

	cfg, err := config.NewConfig(opts.ConfigPath)
	if err != nil {
		log.Panicln(err)
	}

	r := setupBaseRouter(opts.ProductionMode)
	setupRouter(cfg, r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err = r.Run(addr); err != nil {
		log.Panicln(err)
	}
}

func setupBaseRouter(isActivateProdMode bool) *gin.Engine {
	if isActivateProdMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead}

	r.Use(cors.New(corsConfig))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}

func setupRouter(cfg *config.Config, r *gin.Engine) {
	bookRepo, err := repository.NewBookRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	bookService := service.NewBookService(bookRepo)
	bookRouter := v1.NewBookRouter(bookService)
	bookRouter.SetupRouter(r)
}
