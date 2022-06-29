package v1

import (
	"bookman/config"
	"bookman/repository"
	"bookman/service"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, r *gin.Engine) {
	bookRepo, err := repository.NewBookRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	bookService := service.NewBookService(bookRepo)
	bookRouter := NewBookRouter(bookService)
	bookRouter.SetupRouter(r)

	userRepo, err := repository.NewUserRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	userService := service.NewUserService(userRepo)
	userRouter := NewUserRouter(userService)
	userRouter.SetupRouter(r)
}
