package entity

type Book struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Title    string `gorm:"column:title"`
	ISBN     string `gorm:"column:isbn;not null;unique"`
	AuthorID int    `gorn:"column:author_id"`

	Author Author `gorm:"foreignKey:author_id;references:id"`
}

func (*Book) TableName() string {
	return "books"
}
