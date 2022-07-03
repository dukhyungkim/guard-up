package v2

import (
	"bookman/entity"
	"bookman/service"
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
		request.Action,
		&entity.Response[*entity.Book]{
			Data: newBook,
		},
	}, nil
}
