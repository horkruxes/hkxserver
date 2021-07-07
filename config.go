package main

import (
	"io/ioutil"
	"log"

	"github.com/ewenquim/horkruxes/service"
	toml "github.com/pelletier/go-toml"
)

func loadServerConfig() service.ServerConfig {
	getDefaultConfig()
	config, err := toml.LoadFile("config.toml")

	serverConfig := service.ServerConfig{}
	if err != nil {
		panic("Can't load server configuration file (config.toml not found)")
	}
	serverConfig.Name = config.Get("name").(string)
	serverConfig.URL = config.Get("url").(string)
	serverConfig.Info = config.Get("info").(string)
	serverConfig.Private = config.Get("private").(bool)
	serverConfig.Port = config.Get("port").(int64)
	serverConfig.Markdown = config.Get("markdown").(bool)
	serverConfig.Debug = config.Get("testing").(bool)
	serverConfig.TrustedPods = []string{"horkruxes.amethysts.studio", "hk.quimerch.com"}
	return serverConfig
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

# Testing
# Set to false if you want deploy an horkrux instance,
# or true if you want just to run the client or test anything
testing = true # default: true

# URL of the server
url = "localhost" # for local use: "localhost" or "127.0.0.1"

# Custom name to display
name = "horkruxes" # default: "horkruxes"

# Give information to your users
info = """An Horkrux server administrated by xxx.
Create your own server to defeat censorship!
"""

# Is this server private ? 
# If private, only keys added in the /keys folder will be able to post
private = false # default: false

# Port to listen to
port = 80 # default: 80

# Is markdown syntax allowed on this pod ?
markdown = false # default: false
`
}
