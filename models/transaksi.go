package models

import "time"

type Transaksi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	JadwalID    uint      `gorm:"not null;index" json:"jadwal_id"`
	TotalHarga  int       `gorm:"not null" json:"total_harga"`
	MetodeBayar string    `gorm:"size:50" json:"metode_bayar"`
	Status      string    `gorm:"size:30;default:pending" json:"status"`
	KodeBooking string    `gorm:"size:50;uniqueIndex" json:"kode_booking"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Schedule *Schedule `gorm:"foreignKey:JadwalID" json:"jadwal,omitempty"`
	Tiket    []Tiket   `gorm:"foreignKey:TransaksiID" json:"tiket,omitempty"`
}

func (Transaksi) TableName() string { return "transaksis" }
