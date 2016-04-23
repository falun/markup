package main

import (
	"flag"
	"log"
	"os"

	"github.com/falun/markup/app"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Could not get working directory: %s", err)
		pwd = ""
	}

	port := 8080
	host := "localhost"
	index := "README.md"

	flag.StringVar(&pwd, "serve.root", pwd, "the root serving directory")
	flag.StringVar(&host, "serve.ip", host, "the interface we should be listening on")
	flag.IntVar(&port, "serve.port", port, "the port markup will listen on")
	flag.StringVar(
		&index, "serve.index", index,
		"the file returned if '/' is requested; resolved relative to server.root")

	flag.Parse()

	cfg := app.Config{
		RootDir:     pwd,
		Host:        host,
		Port:        port,
		Index:       index,
		BrowseToken: "browse",
		IndexToken:  "index",
	}

	app.Main(cfg)
}
