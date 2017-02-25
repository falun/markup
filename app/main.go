package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/falun/markup/app/scanner"
	"github.com/falun/markup/web"
)

func Main(cfg Config, mimeMap map[string]string) {
	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	renderer := web.NewRenderer(cfg.AssetDir)

	// -- register endpoint handlers

	// redirect to default file
	http.Handle("/", serveRoot("/", cfg.BrowseToken, cfg.Index))

	// view specific markdown
	http.Handle(fmt.Sprintf("/%s/", cfg.BrowseToken), browserHandler{
		Renderer:   renderer,
		host:       ingress,
		root:       cfg.RootDir,
		index:      cfg.Index,
		token:      cfg.BrowseToken,
		indexToken: cfg.IndexToken,
	})

	// directory index view
	http.Handle(fmt.Sprintf("/%s/", cfg.IndexToken),
		serveIndex(renderer, cfg.IndexToken, cfg.BrowseToken, cfg.RootDir, cfg.ExcludeDirs))

	// search indexed files
	http.Handle("/i/search", searchHandler(cfg.RootDir))

	// static assets
	if cfg.AssetToken != "" && cfg.AssetDir != "" {
		http.Handle(fmt.Sprintf("/%s/", cfg.AssetToken),
			serveAsset(cfg.AssetToken, cfg.AssetDir, mimeMap))
	}

	// start root dir scanner, if necessary
	scanNote := ""
	if cfg.RootScanFreq != nil {
		scanner.Start(cfg.RootDir, *cfg.RootScanFreq, cfg.ExcludeDirs)
		scanNote = fmt.Sprintf("Root Dir scan frequency: %v\n", *cfg.RootScanFreq)
	}

	// print some shit saying we're starting up
	assetNote := ""
	if cfg.AssetDir != "" {
		assetNote = fmt.Sprintf("Asset Dir: %s\n", cfg.AssetDir)
	}

	fmt.Printf(
		"Listening on: http://%s\nRoot Dir: %s\nDefault file: %s\n%s%s",
		ingress, cfg.RootDir, cfg.Index, assetNote, scanNote)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}
