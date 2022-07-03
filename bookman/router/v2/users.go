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
		request.Action,
		&entity.Response[*entity.User]{
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

	request.Pagination.Total = 2
	return &PaginatedActionResponse[*entity.User]{
		Action: request.Action,
		PaginatedResponse: &entity.PaginatedResponse[*entity.User]{
			Pagination: &request.Pagination,
			Data: []*entity.User{
				{
					ID:   123,
					Name: "asdf",
				},
				{
					ID:   567,
					Name: "zxcv",
				},
			},
		},
	}, nil
}
