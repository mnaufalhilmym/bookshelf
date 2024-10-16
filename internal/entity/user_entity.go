package entity

type User struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username;not null;unique"`
	Password string `gorm:"column:password"`
}

func (*User) TableName() string {
	return "users"
}
