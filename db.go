package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type FileUserInfoDbClient struct {
	connectionString string
}

type FileNewPathDbClient struct {
	connectionString string
}

func (s *FileUserInfoDbClient) CheckDbHealth() bool {
	// 连接数据库
	db, err := sql.Open("mysql", s.connectionString)
	if err != nil {
		logger.Error("数据库连接失败:", zap.Error(err))
		return false
	}
	defer db.Close()
	return true
}

func (s *FileUserInfoDbClient) getDb() (*sql.DB, error) {
	// 连接数据库
	return sql.Open("mysql", s.connectionString)
}

func (s *FileNewPathDbClient) CheckDbHealth() bool {
	// 连接数据库
	db, err := sql.Open("mysql", s.connectionString)
	if err != nil {
		logger.Error("数据库连接失败:", zap.Error(err))
		return false
	}
	defer db.Close()
	return true
}

func (s *FileNewPathDbClient) getDb() (*sql.DB, error) {
	// 连接数据库
	return sql.Open("mysql", s.connectionString)
}
