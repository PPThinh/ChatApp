package models

import (
	"github.com/google/uuid"
	"time"
)

// NOTE: https://gorm.io/docs/models.html#embedded_struct
type User struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Email       string    `gorm:"type:varchar(100);unique" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(15);unique" json:"phone_number"`
	Password    string    `json:"type:varchar(255);not null" json:"password"`
	CreateAt    time.Time `json:"create_at"`
	UpdateAt    time.Time `json:"update_at"`
}

type PublicUser struct {
}
