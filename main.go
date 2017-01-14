package main

import (
	"flag"
	"log"
	"os"
	"strings"

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
	exclude := ""

	flag.StringVar(&pwd, "index.root", pwd, "[DEPRECATED] the root serving directory")
	flag.StringVar(&host, "serve.ip", host, "[DEPRECATED] the interface we should be listening on")
	flag.IntVar(&port, "serve.port", port, "[DEPRECATED] the port markup will listen on")
	flag.StringVar(
		&index, "serve.default", index,
		"[DEPRECATED] the file returned if '/' is requested; resolved relative to server.root")

	flag.StringVar(&pwd, "root", pwd, "the root serving directory")
	flag.StringVar(&host, "ip", host, "the interface we should be listening on")
	flag.IntVar(&port, "port", port, "the port markup will listen on")
	flag.StringVar(
		&index, "default-file", index,
		"the file returned if '/' is requested; resolved relative to root")

	flag.StringVar(
		&exclude, "index.exclude", exclude,
		"[UNUSED] these directories (comma separated) will not be included in the generated " +
		"/index. They are resolved relative to index.root")

	flag.Parse()

	excludeDirs := []string{}
	if exclude != "" {
		excludeDirs = strings.Split(exclude, ",")
	}

	cfg := app.Config{
		RootDir:     pwd,
		Host:        host,
		Port:        port,
		Index:       index,
		BrowseToken: "browse",
		IndexToken:  "index",
		ExcludeDirs: excludeDirs,
	}

	app.Main(cfg)
}
