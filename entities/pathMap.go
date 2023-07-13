package entities

type PathMap struct {
	OriginPath string `gorm:"column:OriginPath,primaryKey"`
	NewPath    string `gorm:"column:NewPath"`
	Size       int32  `gorm:"column:Size"`
	NeedSync   bool   `gorm:"column:NeedSync"`
}

func (*PathMap) TableName() string {
	return "Path_Maps"
}
