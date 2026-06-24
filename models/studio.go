package models

import "time"

type Studio struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"size:120;not null" json:"nama"`
	Kapasitas int       `gorm:"not null" json:"kapasitas"`
	Tipe      string    `gorm:"size:50" json:"tipe"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Studio) TableName() string { return "studios" }
