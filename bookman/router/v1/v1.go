package v1

import (
	"bookman/config"
	"bookman/repository"
	"bookman/service"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, r *gin.Engine) {
	v1GroupRouter := r.Group("v1")

	bookRepo, err := repository.NewBookRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	bookService := service.NewBookService(bookRepo)

	rentalRepo, err := repository.NewRentalRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	rentalService := service.NewRentalService(rentalRepo)

	bookRouter := NewBookRouter(bookService, rentalService)
	bookRouter.SetupRouter(v1GroupRouter)

	userRepo, err := repository.NewUserRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	userService := service.NewUserService(userRepo)
	userRouter := NewUserRouter(userService)
	userRouter.SetupRouter(v1GroupRouter)

	v1GroupRouter.GET("sse", notifyEvent)
}

func notifyEvent(c *gin.Context) {
	stream := make(chan string)

	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now()
			currentTime := fmt.Sprintf("The Current Time is %v", now)
			stream <- currentTime
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-stream; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
