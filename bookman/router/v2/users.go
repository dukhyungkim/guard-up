package v2

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"
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
		return nil, err
	}

	newUser, err := h.userService.SaveNewUser(&request.User)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.User]{
		Action: request.Action,
		Response: &entity.Response[*entity.User]{
			Data: newUser,
		},
	}, nil
}

func (h *UserHandler) ListUsers(message []byte) (*PaginatedActionResponse[*entity.User], error) {
	var request PaginatedFetchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
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

func (h *UserHandler) GetUser(message []byte) (*ActionResponse[*entity.User], error) {
	var request UserActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	request.User.ID = request.UserID
	user, err := h.userService.GetUser(&request.User)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.User]{
		Action: request.Action,
		Response: &entity.Response[*entity.User]{
			Data: user,
		},
	}, nil
}

func (h *UserHandler) UpdateUser(message []byte) (*ActionResponse[*entity.User], error) {
	var request UserActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	request.User.ID = request.UserID
	updateUser, err := h.userService.UpdateUser(&request.User)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.User]{
		Action: request.Action,
		Response: &entity.Response[*entity.User]{
			Data: updateUser,
		},
	}, nil
}

func (h *UserHandler) DeleteUser(message []byte) (*MessageOKResponse, error) {
	var request UserActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	err = h.userService.DeleteUser(request.UserID)
	if err != nil {
		return nil, err
	}

	return NewMessageOKResponse(request.Action), nil
}
