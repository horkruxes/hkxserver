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

func main() {

	flag.Parse()

	arg := ""
	if len(flag.Args()) > 0 {
		arg = flag.Args()[0]
	}

	switch arg {
	case "version":
		version()
	case "update":
		doUpdate("")
	default:
		app, port := setupServer()
		app.Listen(fmt.Sprintf(":%v", port))
	}
}
