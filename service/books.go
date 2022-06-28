package service

import (
	"bookman/entity"
	"bookman/repository"
)

type BookService interface {
	SaveNewBook(book *entity.Book) (*entity.Book, error)
	ListBooks(pagination *entity.Pagination) ([]*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
	DeleteBook(book *entity.Book) error
}

type bookService struct {
	bookRepo repository.BookRepo
}

func NewBookService(bookRepo repository.BookRepo) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (b *bookService) SaveNewBook(book *entity.Book) (*entity.Book, error) {
	return b.bookRepo.SaveBook(book)
}

func (b *bookService) ListBooks(pagination *entity.Pagination) ([]*entity.Book, error) {
	return b.bookRepo.ListBooks(pagination)
}

func (b *bookService) UpdateBook(book *entity.Book) (*entity.Book, error) {
	fetchBook, err := b.bookRepo.FetchBook(book)
	if err != nil {
		return nil, err
	}
	fetchBook.Name = book.Name
	return b.bookRepo.SaveBook(fetchBook)
}

func (b *bookService) DeleteBook(book *entity.Book) error {
	//TODO implement me
	panic("implement me")
}
