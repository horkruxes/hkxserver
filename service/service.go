package service

import (
	"gorm.io/gorm"
)

type Service struct {
	DB           *gorm.DB
	ServerConfig ServerConfig // Only loaded on start up
	ClientConfig ClientConfig // Can be modified by clients requests
}

type ServerConfig struct {
	Name        string // Custom name
	URL         string // Custom name
	Private     bool   // If set true, a public key must have been uploaded on the server
	Port        int64  // Port to listen
	Info        string // Free text to explain what is this pod
	Markdown    bool   // Is markdown allowed on this pod ?
	Debug       bool
	TrustedPods []string // Public pods
	TrustedKeys []string // Trusted keys, without heavy limitations
}

type ClientConfig struct {
	Pods []Pod // List of pods
}

type Pod struct {
	URL      string
	Selected bool
}
