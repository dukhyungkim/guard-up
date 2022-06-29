package entity

type User struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "users"
}
