package main

import (
	"io/ioutil"
	"log"

	"github.com/ewenquim/horkruxes-client/service"
	toml "github.com/pelletier/go-toml"
)

func loadServerConfig() service.ServerConfig {
	getDefaultConfig()
	config, err := toml.LoadFile("config.toml")

	serverConfig := service.ServerConfig{}
	if err == nil {
		serverConfig.Name = config.Get("name").(string)
		serverConfig.URL = config.Get("url").(string)
		serverConfig.Private = config.Get("private").(bool)
		serverConfig.Port = config.Get("port").(int64)
		return serverConfig
	}
	panic("Can't load server configurration file (config.tml)")
}

// Creates a default config fil if it doesn't exist
func getDefaultConfig() {
	_, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Println("Creating default config...")
		b := []byte(defaultConfig())
		ioutil.WriteFile("config.toml", b, 0644)
	}
}

func defaultConfig() string {

	return `### HORKRUXES CONFIGURATION ###

# URL of the server
url = "localhost" # for local use: "localhost" or "127.0.0.1"

# Custom name to display
name = "horkruxes" # default: "horkruxes"

# Is this server private ? 
# If private, only keys added in the /keys folder will be able to post
private = false # default: false

# Port to listen to
port = 80 # default: 80
`
}
