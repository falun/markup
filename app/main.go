package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/falun/markup/web"
)

func Main(cfg Config) {
	renderer := web.NewRenderer(cfg.AssetDir)
	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	http.Handle("/", serveRoot("/", cfg.BrowseToken, cfg.Index))

	http.Handle(fmt.Sprintf("/%s/", cfg.BrowseToken), browserHandler{
		Renderer: renderer,
		host:       ingress,
		root:       cfg.RootDir,
		index:      cfg.Index,
		token:      cfg.BrowseToken,
		indexToken: cfg.IndexToken,
	})

	http.Handle(fmt.Sprintf("/%s/", cfg.IndexToken),
		serveIndex(renderer, cfg.IndexToken, cfg.BrowseToken, cfg.RootDir, cfg.ExcludeDirs))

	if cfg.AssetToken != "" && cfg.AssetDir != "" {
		http.Handle(fmt.Sprintf("/%s/", cfg.AssetToken),
			serveAsset(cfg.AssetToken, cfg.AssetDir))
	}

	assetNote := ""
	if cfg.AssetDir != "" {
		assetNote = fmt.Sprintf("\nAssets Token: %s\nAsset Dir: %s", cfg.AssetToken, cfg.AssetDir)
	}
	fmt.Printf(`Listening on: http://%s
Serving from: %s
Default file: %s%s
`, ingress, cfg.RootDir, cfg.Index, assetNote)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}
