package repository

import (
	"bookman/common"
	"bookman/config"
	"bookman/entity"
	"errors"

	"gorm.io/gorm"
)

type BookRepo interface {
	SaveBook(book *entity.Book) (*entity.Book, error)
	FetchBook(book *entity.Book) (*entity.Book, error)
	ListBooks(pagination *entity.Pagination) ([]*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
	DeleteBook(book *entity.Book) error
}

type bookRepo struct {
}

func NewBookRepo(rdb *config.RDB) (BookRepo, error) {
	err := initClient(rdb)
	if err != nil {
		return nil, err
	}
	return &bookRepo{}, nil
}

func (r *bookRepo) SaveBook(book *entity.Book) (*entity.Book, error) {
	if err := db.Create(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) FetchBook(book *entity.Book) (*entity.Book, error) {
	if err := db.First(book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundBook(err)
		}
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) ListBooks(pagination *entity.Pagination) ([]*entity.Book, error) {
	var books []*entity.Book
	if err := db.Scopes(paginate(books, pagination)).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepo) UpdateBook(book *entity.Book) (*entity.Book, error) {
	if err := db.Save(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) DeleteBook(book *entity.Book) error {
	if err := db.Delete(book).Error; err != nil {
		return err
	}
	return nil
}
