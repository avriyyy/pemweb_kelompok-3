package models

import "time"

type Tiket struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TransaksiID  uint      `gorm:"not null;index" json:"transaksi_id"`
	JadwalID     uint      `gorm:"not null;index" json:"jadwal_id"`
	NomorKursi   string    `gorm:"size:10;not null" json:"nomor_kursi"`
	Harga        int       `gorm:"not null" json:"harga"`
	Status       string    `gorm:"size:30;default:aktif" json:"status"`
	KodeQR       string    `gorm:"size:255" json:"kode_qr"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Transaksi *Transaksi `gorm:"foreignKey:TransaksiID" json:"transaksi,omitempty"`
	Jadwal    *Jadwal    `gorm:"foreignKey:JadwalID" json:"jadwal,omitempty"`
}

func (Tiket) TableName() string { return "tikets" }
