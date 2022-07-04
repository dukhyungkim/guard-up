package v2

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"encoding/json"
	"errors"
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
		return nil, err
	}

	newBook, err := h.bookService.SaveNewBook(&request.Book)
	if err != nil {
		return nil, err
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

func (h *BookHandler) UpdateBook(message []byte) (*ActionResponse[*entity.Book], error) {
	var request BookActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	request.Book.ID = request.BookID
	updateBook, err := h.bookService.UpdateBook(&request.Book)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.Book]{
		Action: request.Action,
		Response: &entity.Response[*entity.Book]{
			Data: updateBook,
		},
	}, nil
}

func (h *BookHandler) DeleteBook(message []byte) (*MessageOKResponse, error) {
	var request BookActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	err = h.bookService.DeleteBook(request.BookID)
	if err != nil {
		return nil, err
	}

	return NewMessageOKResponse(request.Action), nil
}

func (h *BookHandler) Status(message []byte) (*ActionResponse[*entity.RentalStatus], error) {
	var request BookActionRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	rentalStatus, err := h.rentalService.GetRentStatus(request.BookID)
	if err != nil {
		customErr := &common.Err{}
		if errors.As(err, &customErr) {
			if customErr.Code == common.ErrNotFoundRentalStatus(nil).Code {
				return &ActionResponse[*entity.RentalStatus]{
					Action: request.Action,
				}, nil
			}
		}

		return nil, err
	}

	return &ActionResponse[*entity.RentalStatus]{
		Action: request.Action,
		Response: &entity.Response[*entity.RentalStatus]{
			Data: rentalStatus,
		},
	}, nil
}

func (h *BookHandler) StartRental(message []byte) (*ActionResponse[*entity.RentalStatus], error) {
	var request RentalRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	rentalStatus, err := h.rentalService.StartRentBook(request.BookID, request.UserID)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.RentalStatus]{
		Action: request.Action,
		Response: &entity.Response[*entity.RentalStatus]{
			Data: rentalStatus,
		},
	}, nil
}

func (h *BookHandler) EndRental(message []byte) (*ActionResponse[*entity.RentalStatus], error) {
	var request RentalRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, err
	}

	err = h.rentalService.EndRentBook(request.BookID, request.UserID)
	if err != nil {
		return nil, err
	}

	return &ActionResponse[*entity.RentalStatus]{
		Action: request.Action,
	}, nil
}
