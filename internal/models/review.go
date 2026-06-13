package models

type Review struct {
	ID     uint   `gorm:"primaryKey"`
	BookID uint
	UserID uint
	Rating int
	Comment string
}
