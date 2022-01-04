package service

import (
	"regexp"

	"gorm.io/gorm"
)

// The Service struct is used to store server information.
// Note that it is (should be) only used as a data provider.
// Storing state is a bad idea.
type Service struct {
	DB            *gorm.DB
	GeneralConfig GeneralConfig // Only loaded on start up - public information
	ServerConfig  ServerConfig  // Only loaded on start up - private configuration
	ClientConfig  ClientConfig  // Can be modified by clients requests and is reset between 2 requests
	Regexes       regexes       // It is here because it is loaded only on server startup
}

type regexes struct {
	URLs *regexp.Regexp // Parse urls
}

type GeneralConfig struct {
	Name        string   // Custom name
	URL         string   // URL of the pod
	Description string   // Free text to explain what is this pod
	Info        string   // Contact information and a lot of other things
	TrustedPods []string // Public pods trusted by the admin of this pod
}

type ServerConfig struct {
	Enabled         bool
	Port            int64    // Port to listen
	Debug           bool     // Enable/disable some parameters (limiter, logging...)
	Private         bool     // If set true, a public key must have been uploaded on the server
	LockedByDefault bool     // If set to true, the post route will be disabled
	TrustedKeys     []string // Trusted keys, without heavy limitations
}

type ClientConfig struct {
	Enabled    bool
	Markdown   bool // Is markdown allowed on this pod ?
	PublicPods bool
}

// InitializeDetectors is heavy and should be initialized once only
func InitializeDetectors() regexes {
	return regexes{}
}
