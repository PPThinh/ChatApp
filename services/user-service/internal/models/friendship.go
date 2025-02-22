package models

import "github.com/google/uuid"

type FriendShip struct {
	UserID   uuid.UUID `gorm:"type:char(36);index" json:"user_id"`
	FriendID uuid.UUID `gorm:"type:char(36);index" json:"friend_id"`
}

type FriendRequest struct {
	FromUserID uuid.UUID `gorm:"type:char(36);index" json:"from_user_id"`
	ToUserID   uuid.UUID `gorm:"type:char(36);index" json:"to_user_id"`
	// todo: implemenet enum status here
	Status string `gorm:"size:10" json:"status"`
}
