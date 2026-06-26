package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FilmID   uint
	StudioID uint

	TanggalTayang time.Time
	JamTayang     string
	Harga         float64
	Status        string `gorm:"size:20;default:Aktif"`

	Film   Film   `gorm:"foreignKey:FilmID"`
	Studio Studio `gorm:"foreignKey:StudioID"`
}

func (Schedule) TableName() string {
	return "schedules"
}
