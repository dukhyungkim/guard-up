package entity

import "github.com/lib/pq"

type Book struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Authors pq.StringArray `json:"authors" gorm:"type:text[]"`
	Image   string         `json:"image"`
}

func (Book) TableName() string {
	return "books"
}
