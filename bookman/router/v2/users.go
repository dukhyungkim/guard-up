package v2

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"
	"log"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SaveUser(message []byte) (*ActionResponse[*entity.User], error) {
	var request UserActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	newBook := &entity.User{
		ID:   123,
		Name: request.User.Name,
	}

	return &ActionResponse[*entity.User]{
		Action: request.Action,
		Response: &entity.Response[*entity.User]{
			Data: newBook,
		},
	}, nil
}

func (h *UserHandler) ListUsers(message []byte) (*PaginatedActionResponse[*entity.User], error) {
	var request PaginatedFetchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if request.Pagination.Limit == 0 {
		request.Pagination.Limit = util.DefaultLimit
	}

	users, err := h.userService.ListUsers(&request.Pagination)
	if err != nil {
		return nil, err
	}

	return &PaginatedActionResponse[*entity.User]{
		Action: request.Action,
		PaginatedResponse: &entity.PaginatedResponse[*entity.User]{
			Pagination: &request.Pagination,
			Data:       users,
		},
	}, nil
}
