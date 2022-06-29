package v1

import (
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const keyBookID = "bookId"

type BooksRouter struct {
	bookService   service.BookService
	rentalService service.RentalService
}

func NewBookRouter(bookService service.BookService, rentalService service.RentalService) *BooksRouter {
	return &BooksRouter{
		bookService:   bookService,
		rentalService: rentalService,
	}
}

func (r *BooksRouter) SetupRouter(router *gin.Engine) {
	booksGroup := router.Group("/v1/books")
	booksGroup.POST("", r.createBook)
	booksGroup.GET("", r.listBooks)
	booksGroup.PUT(":"+keyBookID, r.updateBook)
	booksGroup.DELETE(":"+keyBookID, r.deleteBook)

	booksGroup.GET(":"+keyBookID+"/status", r.getRentStatus)
	booksGroup.POST(":"+keyBookID+"/rent", r.handleRent)
}

func (r *BooksRouter) createBook(c *gin.Context) {
	book, err := util.ParseBody[entity.Book](c)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	newBook, err := r.bookService.SaveNewBook(book)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	result := entity.Response[*entity.Book]{Data: newBook}
	c.JSON(http.StatusCreated, result)
}

func (r *BooksRouter) listBooks(c *gin.Context) {
	pagination := util.NewPaginationFromRequest(c)

	books, err := r.bookService.ListBooks(pagination)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	result := entity.PaginatedResponse[*entity.Book]{
		Pagination: pagination,
		Data:       books,
	}
	c.JSON(http.StatusOK, result)
}

func (r *BooksRouter) updateBook(c *gin.Context) {
	bookID, err := util.ParseID(c, keyBookID)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	book, err := util.ParseBody[entity.Book](c)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	book.ID = bookID
	updateBook, err := r.bookService.UpdateBook(book)
	if err != nil {
		util.ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, updateBook)
}

func (r *BooksRouter) deleteBook(c *gin.Context) {
	bookID, err := util.ParseID(c, keyBookID)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	book := &entity.Book{ID: bookID}
	err = r.bookService.DeleteBook(book)
	if err != nil {
		util.ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

const (
	rentalStatusStart = "start"
	rentalStatusEnd   = "end"
)

type rentReqBody struct {
	UserID int `json:"userId"`
}

func (r *BooksRouter) getRentStatus(c *gin.Context) {
	bookID, err := util.ParseID(c, keyBookID)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	rentStatus, err := r.rentalService.GetRentStatus(bookID)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	result := entity.Response[*entity.RentalStatus]{
		Data: rentStatus,
	}
	c.JSON(http.StatusOK, result)
}

func (r *BooksRouter) handleRent(c *gin.Context) {
	bookID, err := util.ParseID(c, keyBookID)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	userID, err := util.ParseBody[rentReqBody](c)
	if err != nil {
		util.ResponseError(c, err)
		return
	}

	status := strings.ToLower(c.Query("status"))
	switch status {
	case rentalStatusStart:
		rentStatus, err := r.rentalService.StartRentBook(bookID, userID.UserID)
		if err != nil {
			util.ResponseError(c, err)
			return
		}
		c.JSON(http.StatusOK, rentStatus)
		return

	case rentalStatusEnd:
		err = r.rentalService.EndRentBook(bookID, userID.UserID)
		if err != nil {
			util.ResponseError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
