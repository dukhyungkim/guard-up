package service

import (
	"bookman/entity"
	"bookman/events"
	"bookman/repository"
)

var bookServiceInstance BookService

type BookService interface {
	SaveNewBook(book *entity.Book) (*entity.Book, error)
	ListBooks(pagination *entity.Pagination) ([]*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
	DeleteBook(bookID int) error
}

type bookService struct {
	bookRepo    repository.BookRepo
	eventSender events.EventSender
}

func NewBookService(bookRepo repository.BookRepo, eventSender events.EventSender) BookService {
	if bookServiceInstance != nil {
		return bookServiceInstance
	}

	bookServiceInstance = &bookService{
		bookRepo:    bookRepo,
		eventSender: eventSender,
	}
	return bookServiceInstance
}

func (s *bookService) SaveNewBook(book *entity.Book) (*entity.Book, error) {
	saveBook, err := s.bookRepo.SaveBook(book)
	if err != nil {
		return nil, err
	}

	defer func() {
		go s.eventSender(events.EventAddedBook, saveBook)
	}()

	return saveBook, nil
}

func (s *bookService) ListBooks(pagination *entity.Pagination) ([]*entity.Book, error) {
	return s.bookRepo.ListBooks(pagination)
}

func (s *bookService) UpdateBook(book *entity.Book) (*entity.Book, error) {
	_, err := s.bookRepo.FetchBook(&entity.Book{ID: book.ID})
	if err != nil {
		return nil, err
	}

	updateBook, err := s.bookRepo.UpdateBook(book)
	if err != nil {
		return nil, err
	}

	defer func() {
		go s.eventSender(events.EventUpdatedBook, updateBook)
	}()

	return updateBook, nil
}

func (s *bookService) DeleteBook(bookID int) error {
	book := &entity.Book{ID: bookID}
	foundBook, err := s.bookRepo.FetchBook(book)
	if err != nil {
		return err
	}

	err = s.bookRepo.DeleteBook(book)
	if err != nil {
		return err
	}

	defer func() {
		go s.eventSender(events.EventDeletedBook, foundBook)
	}()

	return nil
}
