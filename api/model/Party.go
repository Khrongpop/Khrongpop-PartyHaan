package model

import "time"

type Party struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	NumberOfMember     int64  `json:"numberOfMember"`
	CreatedAt time.Time `json:"createdAt"`
	OwnerID int `json:"ownerId"`

	Owner User `json:"owner"`
	JoinUsers []*User `gorm:"many2many:user_parties;" json:"joinUsers"`
}


type PartyRequest struct {
	Name      string `json:"name"`
	NumberOfMember     string  `json:"numberOfMember"`
}
