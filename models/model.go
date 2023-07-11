package models

import "time"

type RequestModel struct {
	Paths []string `validate:"required"`
}

type ResponseModel struct {
	Paths []string
}

type FileUserInformation struct {
	OriginPath           string
	NewPath              string
	CreatedUserCode      string
	CreatedTime          time.Time
	LastModifiedUserCode string
	LastModifiedTime     time.Time
}
