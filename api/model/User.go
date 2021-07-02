package model

import "time"


type User struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Email  string `gorm:"unique" form:"email" binding:"required" json:"email"`
	Password  string `form:"password" binding:"required" json:"password"` 
	CreatedAt time.Time `json:"createdAt"`
	Parties []*Party `gorm:"many2many:user_parties;" json:"parties"`
}
