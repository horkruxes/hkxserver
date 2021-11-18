package main

import (
	"embed"
	"flag"
	"fmt"
)

// Embed the entire directory.
//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// @title Horkruxes API
// @version 1.0
// @description Swagger for the Horkruxes API
// @license.name AGPLv3
// @license.url https://www.gnu.org/licenses/agpl-3.0.html
// @BasePath /api
func main() {

	flag.Parse()

	arg := ""
	if len(flag.Args()) > 0 {
		arg = flag.Args()[0]
	}

	switch arg {
	case "version":
		version()
	case "help":
		help()
	case "update":
		err := downloadAndSaveLastVersion()
		if err != nil {
			panic(err)
		}
	default:
		s := setupService()
		app, port := setupServer(s)
		err := app.Listen(fmt.Sprintf(":%v", port))
		if err != nil {
			fmt.Println("Can't start server")
			panic(err)
		}
	}
}
