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

	NamaStudio  string `gorm:"column:nama_studio"`
	JumlahBaris int    `gorm:"column:jumlah_baris"`
	JumlahKolom int    `gorm:"column:jumlah_kolom"`
}

func (Studio) TableName() string {
	return "studios"
}
