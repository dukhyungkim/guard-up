package v1

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const keyBookID = "bookId"

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
	booksGroup.PUT(":"+keyBookID, r.updateBook)
	booksGroup.DELETE(":"+keyBookID, r.deleteBook)

	booksGroup.POST(":"+keyBookID+"/rent", r.handleRent)
}

func (r *BooksRouter) createBook(c *gin.Context) {
	book, err := util.ParseBody[entity.Book](c)
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
	bookID, err := util.ParseID[int](c, keyBookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	book, err := util.ParseBody[entity.Book](c)
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
	bookID, err := util.ParseID[int](c, keyBookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	book := &entity.Book{ID: bookID}
	err = r.bookService.DeleteBook(book)
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(http.StatusNotFound, customErr)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (r *BooksRouter) handleRent(c *gin.Context) {

}
