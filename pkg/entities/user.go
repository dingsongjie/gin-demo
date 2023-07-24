package entities

type User struct {
	Code string `gorm:"column:slug"`
	Id   int32  `gorm:"column:Id,primaryKey"`
}

func (User) TableName() string {
	return "iris_user"
}
