package repository

import (
	"bookman/config"
	"bookman/entity"
)

type BookRepo interface {
	SaveBook(book *entity.Book) (*entity.Book, error)
	FetchBook(book *entity.Book) (*entity.Book, error)
	ListBooks(pagination *entity.Pagination) ([]*entity.Book, error)
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
	if err := db.Save(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) FetchBook(book *entity.Book) (*entity.Book, error) {
	if err := db.First(book).Error; err != nil {
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

func (r *bookRepo) DeleteBook(book *entity.Book) error {
	//TODO implement me
	panic("implement me")
}
