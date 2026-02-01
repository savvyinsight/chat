package model

import (
	"time"

	"gorm.io/gorm"
)

// Message persists chat messages and delivery metadata.
type Message struct {
	gorm.Model
	From        uint       `json:"from" gorm:"index"`
	To          uint       `json:"to,omitempty" gorm:"index"`
	Room        string     `json:"room,omitempty" gorm:"index"`
	Type        string     `json:"type"`
	Body        string     `json:"body" gorm:"type:text"`
	Delivered   bool       `json:"delivered"`
	DeliveredAt *time.Time `json:"delivered_at"`
}

func (Message) TableName() string {
	return "messages"
}
