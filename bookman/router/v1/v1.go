package v1

import (
	"bookman/config"
	"bookman/events"
	"bookman/repository"
	"bookman/service"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, r *gin.Engine, eventManager *events.EventManager) {
	v1GroupRouter := r.Group("v1")

	bookRepo, err := repository.NewBookRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	bookService := service.NewBookService(bookRepo, eventManager.SendEvent)

	rentalRepo, err := repository.NewRentalRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	rentalService := service.NewRentalService(rentalRepo, eventManager.SendEvent)

	bookRouter := NewBookRouter(bookService, rentalService)
	bookRouter.SetupRouter(v1GroupRouter)

	userRepo, err := repository.NewUserRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	userService := service.NewUserService(userRepo, eventManager.SendEvent)
	userRouter := NewUserRouter(userService)
	userRouter.SetupRouter(v1GroupRouter)

	v1GroupRouter.GET("sse", eventManager.HandleSSE)
}
