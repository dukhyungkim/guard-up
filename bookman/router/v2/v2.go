package v2

import (
	"bookman/common"
	"bookman/config"
	"bookman/events"
	"bookman/repository"
	"bookman/service"
	"encoding/json"
	"log"
	"net/http"

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

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	v2GroupRouter.GET("action", newActionHandler(upgrader, bookHandler, userHandler))
	v2GroupRouter.GET("event", eventManager.HandleWebsocketNotify(upgrader))
}

func newActionHandler(upgrader websocket.Upgrader, bookHandler *BookHandler, userHandler *UserHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("websocket: new client connected")
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
			case ActionListBooks:
				result, err = bookHandler.ListBooks(message)
			case ActionUpdateBook:
				result, err = bookHandler.UpdateBook(message)
			case ActionDeleteBook:
				result, err = bookHandler.DeleteBook(message)

			case ActionBookStatus:
				result, err = bookHandler.Status(message)
			case ActionStartRental:
				result, err = bookHandler.StartRental(message)
			case ActionEndRental:
				result, err = bookHandler.EndRental(message)

			case ActionAddUser:
				result, err = userHandler.SaveUser(message)
			case ActionListUsers:
				result, err = userHandler.ListUsers(message)
			case ActionGetUser:
				result, err = userHandler.GetUser(message)
			case ActionUpdateUser:
				result, err = userHandler.UpdateUser(message)
			case ActionDeleteUser:
				result, err = userHandler.DeleteUser(message)

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
