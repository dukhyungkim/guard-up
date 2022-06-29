package entity

import "time"

type RentalStatus struct {
	BookID int       `json:"bookId" gorm:"primaryKey"`
	UserID int       `json:"userId" gorm:"primaryKey"`
	Status string    `json:"status"`
	Start  time.Time `json:"start"`
}

func (RentalStatus) TableName() string {
	return "rental"
}
