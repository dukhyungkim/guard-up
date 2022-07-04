package v2

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"
	"log"
)

type BookHandler struct {
	bookService   service.BookService
	rentalService service.RentalService
}

func NewBookHandler(bookService service.BookService, rentalService service.RentalService) *BookHandler {
	return &BookHandler{
		bookService:   bookService,
		rentalService: rentalService,
	}
}

func (h *BookHandler) SaveBook(message []byte) (*ActionResponse[*entity.Book], error) {
	var request BookActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	newBook := &entity.Book{
		ID:      123,
		Name:    request.Book.Name,
		Authors: request.Book.Authors,
		Image:   request.Book.Image,
	}

	return &ActionResponse[*entity.Book]{
		Action: request.Action,
		Response: &entity.Response[*entity.Book]{
			Data: newBook,
		},
	}, nil
}

func (h *BookHandler) ListBooks(message []byte) (*PaginatedActionResponse[*entity.Book], error) {
	var request PaginatedFetchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if request.Pagination.Limit == 0 {
		request.Pagination.Limit = util.DefaultLimit
	}

	books, err := h.bookService.ListBooks(&request.Pagination)
	if err != nil {
		return nil, err
	}

	return &PaginatedActionResponse[*entity.Book]{
		Action: request.Action,
		PaginatedResponse: &entity.PaginatedResponse[*entity.Book]{
			Pagination: &request.Pagination,
			Data:       books,
		},
	}, nil
}
