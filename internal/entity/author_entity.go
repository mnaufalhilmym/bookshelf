package entity

import "time"

type Author struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Birthdate time.Time `gorm:"column:birthdate"`
}

func (*Author) TableName() string {
	return "authors"
}
