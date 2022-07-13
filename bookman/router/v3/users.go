package v3

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"
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

func (h *UserHandler) SaveUser(s socketio.Conn, msg string) {
	var user entity.User
	err := json.Unmarshal([]byte(msg), &user)
	if err != nil {
		sendReply(s, common.ErrInvalidRequestBody(err))
		return
	}

	newUser, err := h.userService.SaveNewUser(&user)
	if err != nil {
		sendReply(s, err)
		return
	}

	response := &entity.Response[*entity.User]{
		Data: newUser,
	}
	sendReply(s, response)
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

func (h *UserHandler) GetUser(s socketio.Conn, msg string) {
	var userID = struct {
		UserID int `json:"userId"`
	}{}
	err := json.Unmarshal([]byte(msg), &userID)
	if err != nil {
		sendReply(s, common.ErrInvalidRequestBody(err))
		return
	}

	user, err := h.userService.GetUser(&entity.User{ID: userID.UserID})
	if err != nil {
		sendReply(s, err)
		return
	}

	response := &entity.Response[*entity.User]{
		Data: user,
	}
	sendReply(s, response)
}

func (h *UserHandler) UpdateUser(s socketio.Conn, msg string) {
	var user entity.User
	err := json.Unmarshal([]byte(msg), &user)
	if err != nil {
		sendReply(s, common.ErrInvalidRequestBody(err))
		return
	}

	updateUser, err := h.userService.UpdateUser(&user)
	if err != nil {
		sendReply(s, err)
		return
	}

	response := &entity.Response[*entity.User]{
		Data: updateUser,
	}
	sendReply(s, response)
}

func (h *UserHandler) DeleteUser(s socketio.Conn, msg string) {
	var user entity.User
	err := json.Unmarshal([]byte(msg), &user)
	if err != nil {
		sendReply(s, common.ErrInvalidRequestBody(err))
		return
	}

	err = h.userService.DeleteUser(user.ID)
	if err != nil {
		sendReply(s, err)
		return
	}

	sendReply(s, entity.MessageResponseOK)
}
