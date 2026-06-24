package models

import "time"

type Jadwal struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FilmID    uint      `gorm:"not null;index" json:"film_id"`
	StudioID  uint      `gorm:"not null;index" json:"studio_id"`
	Tanggal   time.Time `gorm:"not null" json:"tanggal"`
	JamMulai  string    `gorm:"size:10;not null" json:"jam_mulai"`
	JamSelesai string   `gorm:"size:10" json:"jam_selesai"`
	Harga     int       `gorm:"not null" json:"harga"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Film   *Film   `gorm:"foreignKey:FilmID" json:"film,omitempty"`
	Studio *Studio `gorm:"foreignKey:StudioID" json:"studio,omitempty"`
}

func (Jadwal) TableName() string { return "jadwals" }
