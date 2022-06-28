package v1

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"net/http"

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
	var book entity.Book
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrInvalidRequestBody(err))
		return
	}

	newBook, err := r.bookService.SaveNewBook(&book)
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

}

func (r *BooksRouter) deleteBook(c *gin.Context) {

}

func (r *BooksRouter) handleRent(c *gin.Context) {

}
