package db

import (
	"github.com/dingsongjie/go-project-template/configs"
	"github.com/dingsongjie/go-project-template/pkg/log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

var (
	FileUserInfoDb *gorm.DB
	FileNewPathDb  *gorm.DB
)

func AddDbClients() {
	FileUserInfoDb, _ = getDb(configs.UserInfoConnectionString)
	FileNewPathDb, _ = getDb(configs.NewPathConnectionString)
}

func getDb(connection string) (*gorm.DB, error) {
	// 连接数据库
	config := GenerateGormConfig()
	return gorm.Open(mysql.Open(connection), config)
}

func GenerateGormConfig() *gorm.Config {
	logger := zapgorm2.New(log.Logger)
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	if configs.IsGinInDebug {
		return &gorm.Config{
			Logger: logger.LogMode(gormlogger.Info)}
	} else {
		return &gorm.Config{}
	}
}
