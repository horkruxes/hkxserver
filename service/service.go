package service

import "gorm.io/gorm"

type Service struct {
	DB           *gorm.DB
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Name    string // Custom name
	URL     string // Custom name
	Private bool   // If set true, a public key must have been uploaded on the server
}
