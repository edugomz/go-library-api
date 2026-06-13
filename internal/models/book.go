package models

// type Book struct {
// 	ID       uint   `gorm:"primaryKey"`
// 	Title    string
// 	AuthorID uint
// 	Author   Author
// }

type Book struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	AuthorID uint   `json:"author_id"`
}
