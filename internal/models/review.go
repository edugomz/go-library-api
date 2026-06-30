package models

type Review struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	BookID  uint   `gorm:"not null" json:"book_id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	Rating  int    `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment string `json:"comment"`
}
