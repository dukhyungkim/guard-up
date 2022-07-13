package entity

type Response[T any] struct {
	Data T `json:"data"`
}

type PaginatedResponse[T any] struct {
	Pagination *Pagination `json:"pagination"`
	Data       []T         `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

var MessageResponseOK = MessageResponse{Message: "OK"}
