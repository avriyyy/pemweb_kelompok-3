package models

import (
	"time"

	"gorm.io/gorm"
)

type Studio struct {
	ID uint `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Tipe        string         `gorm:"column:tipe;size:20;default:Reguler"`
	NamaStudio  string `gorm:"column:nama_studio"`
	JumlahBaris int    `gorm:"column:jumlah_baris"`
	JumlahKolom int    `gorm:"column:jumlah_kolom"`
}

func (Studio) TableName() string {
	return "studios"
}
