package v2

import (
	"bookman/common"
	"bookman/config"
	"bookman/events"
	"bookman/repository"
	"bookman/service"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRouter(cfg *config.Config, r *gin.Engine, eventManager *events.EventManager) {
	v2GroupRouter := r.Group("v2")

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

	upgrader := websocket.Upgrader{}

	v2GroupRouter.GET("action", newActionHandler(upgrader, bookHandler, userHandler))
	v2GroupRouter.GET("event", newEventHandler(upgrader))
}

func newActionHandler(upgrader websocket.Upgrader, bookHandler *BookHandler, userHandler *UserHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			if err = ws.Close(); err != nil {
				log.Println(err)
			}
		}()

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			actionRequest := ActionRequest{}
			err = json.Unmarshal(message, &actionRequest)
			if err != nil {
				log.Println(err)
				writeErr := ws.WriteJSON(common.ErrInvalidRequestBody(err))
				if writeErr != nil {
					log.Println(writeErr)
					return
				}
				continue
			}

			var result any
			switch actionRequest.Action {
			case ActionAddBook:
				result, err = bookHandler.SaveBook(message)
			case ActionAddUser:
				result, err = userHandler.SaveUser(message)
			case ActionListUsers:
				result, err = userHandler.ListUsers(message)
			default:
				err = common.ErrNotFoundAction(nil)
			}
			if err != nil {
				writeErr := ws.WriteJSON(err)
				if writeErr != nil {
					log.Println(writeErr)
					return
				}
				continue
			}

			err = ws.WriteJSON(result)
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func newEventHandler(upgrader websocket.Upgrader) func(c *gin.Context) {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			if err = ws.Close(); err != nil {
				log.Println(err)
				return
			}
		}()

		closeChan := make(chan bool)
		go func() {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				closeChan <- true
			}
		}()

	LOOP:
		for {
			select {
			case <-closeChan:
				log.Println("got closed")
				break LOOP
			default:
				pong := struct {
					Message string `json:"message"`
				}{
					Message: "pong",
				}
				b, _ := json.Marshal(pong)

				err = ws.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					log.Println(err)
					break LOOP
				}
				time.Sleep(time.Second)
			}
		}
	}
}
