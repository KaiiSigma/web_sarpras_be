package models

type User struct {
	ID       uint   `gorm:"primaryKey;column:id" json:"id"`
	Nama     string `gorm:"size:255;column:nama" json:"nama"`
	Username string `gorm:"size:255;unique;not null;column:username" json:"username"`
	Password string `gorm:"size:255;not null;column:password" json:"-"`
	Role     string `gorm:"size:255;column:role" json:"role"`
}

func (User) TableName() string {
	return "user"
}