package entity

import "time"

type RentalStatus struct {
	BookID    int       `json:"bookId"`
	UserID    int       `json:"userId"`
	Status    string    `json:"status"`
	RentStart time.Time `json:"rentStart"`
	RentEnd   time.Time `json:"rentEnd"`
}

func (RentalStatus) TableName() string {
	return "rental"
}
