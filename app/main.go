package app

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: test path.Clean usage on windows

func Main(cfg Config) {
	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	http.Handle("/", serveRoot("/", cfg.BrowseToken, cfg.Index))
	http.Handle(fmt.Sprintf("/%s/", cfg.BrowseToken), browserHandler{
		host:       ingress,
		root:       cfg.RootDir,
		index:      cfg.Index,
		token:      cfg.BrowseToken,
		indexToken: cfg.IndexToken,
		renderer:   NewRenderer(),
	})

	http.Handle(fmt.Sprintf("/%s/", cfg.IndexToken),
		serveIndex(cfg.IndexToken, cfg.BrowseToken, cfg.RootDir, cfg.ExcludeDirs))

	fmt.Printf(`Listening on: http://%s
Serving from: %s
Default file: %s
`, ingress, cfg.RootDir, cfg.Index)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}
