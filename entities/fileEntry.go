package entities

import "time"

type FileEntry struct {
	Id                   int32     `gorm:"column:Id,primaryKey"`
	CreatedTime          time.Time `gorm:"column:ctime"`
	LastModifiedTime     time.Time `gorm:"column:mtime"`
	CreatedUserCode      string    `gorm:"column:creator_uid"`
	LastModifiedUserCode string    `gorm:"column:updator_uid"`
	NSid                 string    `gorm:"column:nsid"`
	// CreatedUser          User      `gorm:"foreignKey:CreatedUserCode,references:Id"`
	// LastModifiedUser User `gorm:"foreignKey:LastModifiedUserCode,references:Id"`
}

func (FileEntry) TableName() string {
	return "iris_name_entry"
}
