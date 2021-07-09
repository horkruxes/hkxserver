package main

import (
	"io/ioutil"

	"github.com/ewenquim/horkruxes/service"
	toml "github.com/pelletier/go-toml"
)

func loadServerConfig() service.ServerConfig {
	createDefaultConfigIfDoesNotExist()
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		panic("Can't load server configuration file (config.toml not found or badly formatted)")
	}

	serverConfig := service.ServerConfig{}
	serverConfig.Name = config.Get("name").(string)
	serverConfig.URL = config.Get("url").(string)
	serverConfig.Description = config.Get("description").(string)
	serverConfig.Info = config.Get("info").(string)
	serverConfig.Private = config.Get("private").(bool)
	serverConfig.Port = config.Get("port").(int64)
	serverConfig.Markdown = config.Get("markdown").(bool)
	serverConfig.Debug = config.Get("testing").(bool)
	sourcesInterface := config.Get("sources").([]interface{})
	serverConfig.TrustedPods = make([]string, len(sourcesInterface))
	for i := range sourcesInterface {
		serverConfig.TrustedPods[i] = sourcesInterface[i].(string)
	}
	keysInterface := config.Get("trusted_keys").([]interface{})
	serverConfig.TrustedKeys = make([]string, len(keysInterface))
	for i := range keysInterface {
		serverConfig.TrustedKeys[i] = keysInterface[i].(string)
	}
	return serverConfig
}

// Creates a default config fil if it doesn't exist
func createDefaultConfigIfDoesNotExist() {
	_, err := ioutil.ReadFile("config.toml")
	if err != nil {
		b := []byte(defaultConfig())
		ioutil.WriteFile("config.toml", b, 0644)
		panic(`It's the first time running Horkruxes, Thank you!

The server needs to be configured with a config.toml file.
Luckily for you, a template has been created <3

Just type "vim config.toml" or "nano config.toml" to custom your server, and you'll be ready to go!`)
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

# Short Information
description = """An Horkrux server administrated by xxx.
Create your own server to defeat censorship!"""

# Give more information to your users
info = """Any abuse of freedom (spam...) on this server will be sanctionned (IP ban).
Please stay serious and informative.

Here is the list of my pods:
- yyy
- zzz

Contact me at aaa@bbb.com or ...
"""

# Is this server private ? 
# If private, only keys added in the /keys folder will be able to post
private = false # default: false

# Port to listen to
port = 80 # default: 80

# Is markdown syntax allowed on this pod ?
markdown = false # default: false

# Sources to listen to when trying to get messages
# These must be trustworthy as everyone that will go to your pod on their browsers
# will see the messages from these sources. Beware of spam and quality content.
sources = ['horkruxes.amethysts.studio',
'hk.quimerch.com']

# Trusted Public Keys
# They have only a 5 minute timer between 2 posts
trusted_keys = ['JL6zyYtrk43MZ+uV7J+y8HFS9MvkI2eZT1RbRnV4Qog=']
`
}
