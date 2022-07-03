package v2

import (
	"bookman/entity"
)

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

const (
	ActionAddBook    ActionType = "ADD_BOOK"
	ActionListBooks  ActionType = "LIST_BOOKS"
	ActionUpdateBook ActionType = "UPDATE_BOOK"
	ActionDeleteBOOK ActionType = "DELETE_BOOK"

	ActionBookStatus  ActionType = "BOOK_STATUS"
	ActionStartRental ActionType = "START_RENTAL"
	ActionEndRental   ActionType = "END_RENTAL"

	ActionAddUser    ActionType = "ADD_USER"
	ActionListUsers  ActionType = "LIST_USERS"
	ActionGetUser    ActionType = "GET_USER"
	ActionUpdateUSer ActionType = "UPDATE_USER"
	ActionDeleteUser ActionType = "DELETE_USER"
)

type ActionRequest struct {
	Action ActionType `json:"action"`
}

type BookActionRequest struct {
	ActionRequest
	BookID int         `json:"bookId"`
	Book   entity.Book `json:"request"`
}

type UserActionRequest struct {
	ActionRequest
	UserID int         `json:"userId"`
	User   entity.User `json:"request"`
}

type PaginatedFetchRequest struct {
	ActionRequest
	Pagination entity.Pagination `json:"request"`
}

type RentalRequest struct {
	ActionRequest
	BookID int `json:"bookId"`
	UserID int `json:"userId"`
}

type ActionResponse[T any] struct {
	Action ActionType `json:"action"`
	*entity.Response[T]
}

type PaginatedActionResponse[T any] struct {
	Action ActionType `json:"action"`
	*entity.PaginatedResponse[T]
}
