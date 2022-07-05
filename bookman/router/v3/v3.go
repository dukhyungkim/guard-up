package v3

import (
	"bookman/config"
	"bookman/events"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetupRouter(cfg *config.Config, r *gin.Engine, eventManager *events.EventManager) *socketio.Server {
	v3GroupRouter := r.Group("v3")

	server := setupSocketIOServer(v3GroupRouter.BasePath())
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	v3GroupRouter.GET("/socket.io/*any", gin.WrapH(server))
	v3GroupRouter.POST("/socket.io/*any", gin.WrapH(server))

	return server
}

func setupSocketIOServer(basePath string) *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent(basePath+"/", "echo", func(s socketio.Conn, msg string) {
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
