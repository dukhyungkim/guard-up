package v3

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) ListUsers(s socketio.Conn, msg string) {
	pagination, err := util.NewPaginationFromMessage(msg)
	if err != nil {
		log.Println(err)
		return
	}

	users, err := h.userService.ListUsers(pagination)
	if err != nil {
		log.Println(err)
		return
	}

	response := &entity.PaginatedResponse[*entity.User]{
		Pagination: pagination,
		Data:       users,
	}
	sendReply(s, response)
}
