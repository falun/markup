package app

import (
	"fmt"
	"log"
	"net/http"
)

func Main(cfg Config) {
	http.Handle("/", serveRoot("/", cfg.BrowseToken, cfg.Index))
	http.Handle(fmt.Sprintf("/%s/", cfg.BrowseToken), browserHandler{
		root:     cfg.RootDir,
		index:    cfg.Index,
		token:    cfg.BrowseToken,
		renderer: NewRenderer(),
	})
	http.Handle(fmt.Sprintf("/%s/", cfg.IndexToken),
		serveIndex(cfg.IndexToken, cfg.BrowseToken, cfg.RootDir))

	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf(`
Listening on: %s
Serving from: %s
Default file: %s
`, ingress, cfg.RootDir, cfg.Index)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}
