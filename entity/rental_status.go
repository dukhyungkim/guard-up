package entity

import "time"

type RentalStatus struct {
	BookId    int       `json:"bookId"`
	UserId    int       `json:"userId"`
	Status    string    `json:"status"`
	RentStart time.Time `json:"rent_start"`
	RentEnd   time.Time `json:"rent_end"`
}

func (RentalStatus) TableName() string {
	return "rental"
}
