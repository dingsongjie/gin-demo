package db

import (
	"lenovo-drive-mi-api/log"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type FileUserInfoDbClient struct {
	connectionString string
}

type FileNewPathDbClient struct {
	connectionString string
}

var (
	DefaultFileUserInfoDbClient FileUserInfoDbClient
	DefaultFileNewPathDbClient  FileNewPathDbClient
)

func Config(userInfoConnectionString string, newPathConnectionString string) {
	DefaultFileUserInfoDbClient = FileUserInfoDbClient{userInfoConnectionString}
	DefaultFileNewPathDbClient = FileNewPathDbClient{newPathConnectionString}
}

func (s *FileUserInfoDbClient) CheckDbHealth() bool {
	// 连接数据库
	_, err := s.GetDb()
	if err != nil {

		log.Logger.Error("数据库连接失败:", zap.Error(err))
		return false
	}
	return true
}

func (s *FileUserInfoDbClient) GetDb() (*gorm.DB, error) {
	// 连接数据库
	return gorm.Open(mysql.Open(s.connectionString), &gorm.Config{})
}

func (s *FileNewPathDbClient) CheckDbHealth() bool {
	// 连接数据库
	_, err := s.GetDb()
	if err != nil {
		log.Logger.Error("数据库连接失败:", zap.Error(err))
		return false
	}
	return true
}

func (s *FileNewPathDbClient) GetDb() (*gorm.DB, error) {
	// 连接数据库
	return gorm.Open(mysql.Open(s.connectionString), &gorm.Config{})
}
