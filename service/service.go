package service

import (
	"regexp"

	"gorm.io/gorm"
)

type Service struct {
	DB           *gorm.DB
	ServerConfig ServerConfig // Only loaded on start up
	ClientConfig ClientConfig // Can be modified by clients requests
	Regexes      regexes      // It is here because it is loaded only on server startup
}

type regexes struct {
	URLs *regexp.Regexp // Parse urls
}

type ServerConfig struct {
	Name        string // Custom name
	URL         string // Custom name
	Private     bool   // If set true, a public key must have been uploaded on the server
	Port        int64  // Port to listen
	Description string // Free text to explain what is this pod
	Info        string // Contact information and a lot of other things
	Markdown    bool   // Is markdown allowed on this pod ?
	Debug       bool
	TrustedPods []string // Public pods
	TrustedKeys []string // Trusted keys, without heavy limitations
}

type ClientConfig struct {
	PublicPods             bool
	SpecificPods           bool
	SpecificPodsList       []string // List of pods requested by the user
	SpecificPodsListString string   // List of pods requested by the user and displayed back into the field
	AllPodsList            []Pod    // List of all pods (public + requested)
}

type Pod struct {
	URL      string
	Selected bool
}

func (s *Service) UpdateClientPodsList() {

	s.ClientConfig.AllPodsList = []Pod{}

	// Add public pods from server to client
	if s.ClientConfig.PublicPods {
		for _, name := range s.ServerConfig.TrustedPods {
			s.ClientConfig.AllPodsList = append(s.ClientConfig.AllPodsList, Pod{URL: name, Selected: true})
		}
	}

	// Add custom pods from client to themselves
	if s.ClientConfig.SpecificPods {
		for _, name := range s.ClientConfig.SpecificPodsList {

			// Checks for duplicates
			var alreadyIn bool
			for _, alreadyInPod := range s.ClientConfig.AllPodsList {
				if alreadyInPod.URL == name {
					alreadyIn = true
					break
				}
			}
			if !alreadyIn {
				s.ClientConfig.AllPodsList = append(s.ClientConfig.AllPodsList, Pod{URL: name, Selected: true})
			}
		}
	}
}

// InitializeDetectors is heavy and should be initialized once only
func InitializeDetectors() regexes {
	return regexes{
		URLs: regexp.MustCompile(`[\w.-]+\.[\w.-]+`),
	}
}
