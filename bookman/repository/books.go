package repository

import (
	"bookman/common"
	"bookman/config"
	"bookman/entity"
	"errors"

	"gorm.io/gorm"
)

var bookRepoInstance BookRepo

type BookRepo interface {
	SaveBook(book *entity.Book) (*entity.Book, error)
	FetchBook(book *entity.Book) (*entity.Book, error)
	ListBooks(pagination *entity.Pagination) ([]*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
	DeleteBook(book *entity.Book) error
}

type bookRepo struct {
	repo repo[entity.Book]
}

func NewBookRepo(cfg *config.RDB) (BookRepo, error) {
	if bookRepoInstance != nil {
		return bookRepoInstance, nil
	}

	if err := initClient(cfg); err != nil {
		return nil, err
	}

	bookRepoInstance = &bookRepo{}
	return bookRepoInstance, nil
}

func (r *bookRepo) SaveBook(book *entity.Book) (*entity.Book, error) {
	return r.repo.Save(book)
}

func (r *bookRepo) FetchBook(book *entity.Book) (*entity.Book, error) {
	book, err := r.repo.First(book)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundBook(err)
		}
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) ListBooks(pagination *entity.Pagination) ([]*entity.Book, error) {
	return r.repo.Find(pagination)
}

func (r *bookRepo) UpdateBook(book *entity.Book) (*entity.Book, error) {
	return r.repo.Update(book)
}

func (r *bookRepo) DeleteBook(book *entity.Book) error {
	return r.repo.Delete(book)
}
