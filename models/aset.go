package models

import "time"

type Aset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	KodeAset  string    `gorm:"type:varchar(50);unique;not null" json:"kode_aset"`
	NamaAset  string    `gorm:"type:varchar(100);not null" json:"nama_aset"`
	Jenis     string    `gorm:"type:enum('sarana', 'prasarana');default:'sarana'" json:"jenis"`
	Kategori  string    `gorm:"type:varchar(50);not null" json:"kategori"`
	Jumlah    int       `gorm:"not null;default:1" json:"jumlah"`
	Kondisi   string    `gorm:"type:enum('baik', 'rusak_ringan', 'rusak_berat');default:'baik'" json:"kondisi"`
	Lokasi    string    `gorm:"type:varchar(100);not null" json:"lokasi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}