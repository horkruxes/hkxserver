package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/horkruxes/hkxserver/service"
	toml "github.com/pelletier/go-toml"
)

func loadConfig() service.Service {
	createDefaultConfigIfDoesNotExist()
	config, err := toml.LoadFile("hkxconfig.toml")
	if err != nil {
		panicWithErrorMessage()
	}

	serviceConfig := service.Service{}
	serviceConfig.GeneralConfig.Name = config.Get("general.name").(string)
	serviceConfig.GeneralConfig.URL = config.Get("general.url").(string)
	serviceConfig.GeneralConfig.Description = config.Get("general.description").(string)
	serviceConfig.GeneralConfig.Info = config.Get("general.info").(string)
	sourcesInterface := config.Get("general.sources").([]interface{})
	serviceConfig.GeneralConfig.TrustedPods = make([]string, len(sourcesInterface))
	for i := range sourcesInterface {
		serviceConfig.GeneralConfig.TrustedPods[i] = sourcesInterface[i].(string)
	}

	serviceConfig.ServerConfig.Enabled = config.Get("server.enabled").(bool)
	serviceConfig.ServerConfig.Private = config.Get("server.private").(bool)
	serviceConfig.ServerConfig.Port = config.Get("server.port").(int64)
	serviceConfig.ServerConfig.Debug = config.Get("server.testing").(bool)
	keysInterface := config.Get("server.trusted_keys").([]interface{})
	serviceConfig.ServerConfig.TrustedKeys = make([]string, len(keysInterface))
	for i := range keysInterface {
		serviceConfig.ServerConfig.TrustedKeys[i] = keysInterface[i].(string)
	}

	serviceConfig.ClientConfig.Enabled = config.Get("client.enabled").(bool)
	serviceConfig.ClientConfig.Markdown = config.Get("client.markdown").(bool)

	return serviceConfig
}

// Creates a default config fil if it doesn't exist
func createDefaultConfigIfDoesNotExist() {
	_, err := ioutil.ReadFile("hkxconfig.toml")
	if err != nil {
		b := []byte(defaultConfig())
		ioutil.WriteFile("hkxconfig.toml", b, 0600)
		fmt.Println(`It's the first time running Horkruxes, Thank you!

The server needs to be configured with a "hkxconfig.toml" file.
Luckily for you, a template has been created <3

Just type "vim hkxconfig.toml" or "nano hkxconfig.toml" to custom your server, and you'll be ready to go!

Then, re-run horkruxes. No further configuration is required!

(of course, you could simply re-run horkruxes without editing the hkxconfig.toml file, but that's not the spirit)`)
	}
}

func defaultConfig() string {
	_, lock, _ := ed25519.GenerateKey(nil)
	_, unlock, _ := ed25519.GenerateKey(nil)
	lock_url := base64.URLEncoding.EncodeToString(lock)
	unlock_url := base64.URLEncoding.EncodeToString(unlock)

	return fmt.Sprintf(`###############################
### HORKRUXES CONFIGURATION ###
###############################

#####################################################
### GENERAL - Public information about the server ###
#####################################################

[general]

# URL of the server (ex: hk.quimerch.com, do not type http/https)
url = "localhost" # for local use: "localhost" or "127.0.0.1"

# Custom name to display
name = "Horkruxes" # default: "Horkruxes"

# Short Description
description = """An Horkrux server administrated by xxx.
Create your own server to defeat censorship!"""

# More Information (other nodes, contact, policy...)
info = """Any abuse of freedom (spam...) on this server will be sanctionned (IP ban).
Please stay serious and informative.

Here is the list of my pods:

- yyy
- zzz

Contact me at aaa@bbb.com or ...
"""

# Sources to listen to when trying to get messages
# These must be trustworthy as everyone that will go to your pod on their browsers
# will see the messages from these sources. Beware of spam and quality content.
sources = ['horkruxes.amethysts.studio',
'hk.quimerch.com']



#############################################
### SERVER - Private server configuration ###
#############################################

[server]
# If false, will only act as a "relay" client : no data on this server, only data from other servers
enabled = true # default: true 

# Port to listen to
port = 8000 # default: 8000

# Set to false if you want deploy an horkrux instance,
# or true if you want just to run the client or test anything
testing = true # default: true

# Is this server private ? 
# If private, only trusted keys will be able to post (details below)
private = false # default: false

# Trusted Public Keys
# They have only a 5 minute timer between 2 posts and some privileged rights
trusted_keys = ['JL6zyYtrk43MZ+uV7J+y8HFS9MvkI2eZT1RbRnV4Qog=']

# If the pod is locked, it is set to readonly: nobody can post on the pod THROUGH THE API
# If set to true, emergency lock/unlock will not work. Your data will remain unalterable no matter what.
# You might want to create a bot to post and sign directly to the database.
locked = false # default: false

# If enabled, going to https://pod_url.com/lock/<emergency_lock_url> will set the pod to readonly.
# Can be set quickly in case of emergency
# To recover, you can use the emergency_unlock_url below, or reset the server.
# The url must be really complicated as anyone can lock your pod.
# The one already existing is randomly generated but you can set it to whatever you like
emergency_lock_enabled = false
emergency_lock_url = "%v"

# If enabled, going to https://pod_url.com/unlock/<emergency_unlock_url> will unset the pod readonly.
# Can be set to recover an emergency lock
# Must be really complicated as anyone can unlock your pod.
# The one already existing is randomly generated but you can set it to whatever you like
emergency_unlock_enabled = false
emergency_unlock_url  = "%v"



#################################################################
### CLIENT - Client information, useful only to the interface ###
#################################################################

[client]
# If false, it will only behave like an API
# Can be used to limit charge on the server
enabled = true # default: true 

# Is markdown syntax allowed on this pod ?
markdown = false # default: false
`, lock_url, unlock_url)
}

func panicWithErrorMessage() {
	panic(`Can't load server configuration file (hkxconfig.toml not found or badly formatted).
You might simply have forgotten a quote or an equal sign.

If you delete or rename the 'hkxconfig.toml' file, and then relaunch the app,
Horxruxes will automatically regenerate a correct 'hkxconfig.toml' file.
`)
}
