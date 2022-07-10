package v3

import (
	"bookman/config"
	"bookman/events"
	"bookman/repository"
	v2 "bookman/router/v2"
	"bookman/service"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetupRouter(cfg *config.Config, r *gin.Engine, eventManager *events.EventManager) *socketio.Server {
	v3GroupRouter := r.Group("v3")

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
	bookHandler := NewBookHandler(bookService, rentalService)

	userRepo, err := repository.NewUserRepo(&cfg.Database)
	if err != nil {
		log.Panicln(err)
	}
	userService := service.NewUserService(userRepo, eventManager.SendEvent)
	userHandler := NewUserHandler(userService)

	server := setupSocketIOServer(v3GroupRouter.BasePath(), bookHandler, userHandler)
	go func() {
		if err = server.Serve(); err != nil {
			log.Panicf("socketio listen error: %s\n", err)
		}
	}()

	v3GroupRouter.GET("/socket.io/*any", gin.WrapH(server))
	v3GroupRouter.POST("/socket.io/*any", gin.WrapH(server))

	return server
}

const EventReply = "reply"

func setupSocketIOServer(basePath string, bookHandler *BookHandler, userHandler *UserHandler) *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent(basePath+"/action", v2.ActionAddBook.String(), bookHandler.CreateBook)
	server.OnEvent(basePath+"/action", v2.ActionListBooks.String(), bookHandler.ListBooks)
	server.OnEvent(basePath+"/action", v2.ActionUpdateBook.String(), bookHandler.UpdateBook)
	server.OnEvent(basePath+"/action", v2.ActionDeleteBook.String(), bookHandler.DeleteBook)

	server.OnEvent(basePath+"/action", v2.ActionListUsers.String(), userHandler.ListUsers)

	server.OnEvent(basePath+"/notify", "", func(s socketio.Conn, msg string) {
		log.Println("echo:", msg)
		s.Emit("reply", "echo "+msg)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", s.ID(), reason)
	})

	return server
}

func sendReply[T any](s socketio.Conn, response T) {
	b, _ := json.Marshal(response)
	s.Emit(EventReply, string(b))
}
