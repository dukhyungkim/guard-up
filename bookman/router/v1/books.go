package v1

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BooksRouter struct {
	bookService service.BookService
}

func NewBookRouter(bookService service.BookService) *BooksRouter {
	return &BooksRouter{
		bookService: bookService,
	}
}

func (r *BooksRouter) SetupRouter(router *gin.Engine) {
	booksGroup := router.Group("/v1/books")
	booksGroup.POST("", r.createBook)
	booksGroup.GET("", r.listBooks)
	booksGroup.PUT(":bookId", r.updateBook)
	booksGroup.DELETE(":bookId", r.deleteBook)

	booksGroup.POST(":bookId/rent", r.handleRent)
}

func (r *BooksRouter) createBook(c *gin.Context) {
	book, err := parseBody[entity.Book](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newBook, err := r.bookService.SaveNewBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
		return
	}

	result := entity.Response[*entity.Book]{Data: newBook}
	c.JSON(http.StatusCreated, result)
}

func (r *BooksRouter) listBooks(c *gin.Context) {
	pagination := util.NewPaginationFromRequest(c)

	books, err := r.bookService.ListBooks(pagination)
	if err != nil {
		return
	}

	result := entity.PaginatedResponse[*entity.Book]{
		Pagination: pagination,
		Data:       books,
	}
	c.JSON(http.StatusOK, result)
}

func (r *BooksRouter) updateBook(c *gin.Context) {
	bookID, err := parseBookID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	book, err := parseBody[entity.Book](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	book.ID = bookID
	updateBook, err := r.bookService.UpdateBook(book)
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(http.StatusNotFound, customErr)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, updateBook)
}

func (r *BooksRouter) deleteBook(c *gin.Context) {
	bookID, err := parseBookID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	book := &entity.Book{ID: bookID}
	err = r.bookService.DeleteBook(book)
	if err != nil {
		if errors.As(err, &common.Err{}) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (r *BooksRouter) handleRent(c *gin.Context) {

}

func parseBookID(c *gin.Context) (int, error) {
	bookID := c.Param("bookId")
	id, err := strconv.Atoi(bookID)
	if err != nil {
		return 0, common.ErrInvalidParam(err)
	}
	return id, nil
}

func parseBody[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.ShouldBind(&data); err != nil {
		return nil, common.ErrInvalidRequestBody(err)
	}
	return &data, nil
}
