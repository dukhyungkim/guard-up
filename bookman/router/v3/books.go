package v3

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
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

func (h *BookHandler) CreateBook(s socketio.Conn, msg string) {
	var book entity.Book
	err := json.Unmarshal([]byte(msg), &book)
	if err != nil {
		sendReply(s, common.ErrInvalidRequestBody(err))
		return
	}

	newBook, err := h.bookService.SaveNewBook(&book)
	if err != nil {
		sendReply(s, err)
		return
	}

	response := &entity.Response[*entity.Book]{
		Data: newBook,
	}
	sendReply(s, response)
}

func (h *BookHandler) ListBooks(s socketio.Conn, msg string) {
	pagination, err := util.NewPaginationFromMessage(msg)
	if err != nil {
		sendReply(s, common.ErrInvalidParam(err))
		return
	}

	books, err := h.bookService.ListBooks(pagination)
	if err != nil {
		sendReply(s, err)
		return
	}

	response := &entity.PaginatedResponse[*entity.Book]{
		Pagination: pagination,
		Data:       books,
	}
	sendReply(s, response)
}

func (h *BookHandler) UpdateBook(s socketio.Conn, msg string) {

}

func (h *BookHandler) DeleteBook(s socketio.Conn, msg string) {

}
