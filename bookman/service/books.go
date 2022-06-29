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

func (s *bookService) SaveNewBook(book *entity.Book) (*entity.Book, error) {
	return s.bookRepo.SaveBook(book)
}

func (s *bookService) ListBooks(pagination *entity.Pagination) ([]*entity.Book, error) {
	return s.bookRepo.ListBooks(pagination)
}

func (s *bookService) UpdateBook(book *entity.Book) (*entity.Book, error) {
	_, err := s.bookRepo.FetchBook(&entity.Book{ID: book.ID})
	if err != nil {
		return nil, err
	}
	return s.bookRepo.UpdateBook(book)
}

func (s *bookService) DeleteBook(book *entity.Book) error {
	_, err := s.bookRepo.FetchBook(&entity.Book{ID: book.ID})
	if err != nil {
		return err
	}
	return s.bookRepo.DeleteBook(book)
}
