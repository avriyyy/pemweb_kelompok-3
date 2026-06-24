package models

import "time"

type Film struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Judul       string    `gorm:"size:200;not null" json:"judul"`
	Genre       string    `gorm:"size:100" json:"genre"`
	Durasi      int       `gorm:"not null" json:"durasi"`
	Sinopsis    string    `gorm:"type:text" json:"sinopsis"`
	Poster      string    `gorm:"size:255" json:"poster"`
	Rating      string    `gorm:"size:10" json:"rating"`
	TanggalRilis time.Time `json:"tanggal_rilis"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Film) TableName() string { return "films" }
