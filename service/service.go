package service

import "gorm.io/gorm"

type Service struct {
	DB *gorm.DB
}
