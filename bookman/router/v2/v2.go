package v2

import (
	"bookman/config"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRouter(cfg *config.Config, r *gin.Engine) {
	v2GroupRouter := r.Group("v2")

	v2GroupRouter.GET("websocket", HandleEvent)

}

func HandleEvent(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(mt, string(message))
		pong := struct {
			Message string `json:"message"`
		}{
			Message: "pong",
		}
		b, _ := json.Marshal(pong)

		err = ws.WriteMessage(mt, b)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
