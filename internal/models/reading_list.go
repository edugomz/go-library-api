package models

type ReadingList struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint
	Name   string
	Books  []Book `gorm:"many2many:reading_list_books;"`
}
