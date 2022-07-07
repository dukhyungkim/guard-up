package v3

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"log"

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

func (h *BookHandler) ListBooks(s socketio.Conn, msg string) {
	pagination, err := util.NewPaginationFromMessage(msg)
	if err != nil {
		log.Println(err)
		return
	}

	books, err := h.bookService.ListBooks(pagination)
	if err != nil {
		log.Println(err)
		return
	}

	response := &entity.PaginatedResponse[*entity.Book]{
		Pagination: pagination,
		Data:       books,
	}
	sendReply(s, response)
}
