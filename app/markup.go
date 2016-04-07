package app

import (
	"fmt"
	"log"
	"net/http"
)

func Main(cfg Config) {
	http.Handle("/", serveRoot("/", cfg.Token, cfg.Index))
	http.Handle(fmt.Sprintf("/%s/", cfg.Token), browserHandler{
		root:     cfg.RootDir,
		index:    cfg.Index,
		token:    cfg.Token,
		renderer: NewRenderer(),
	})
	http.Handle("/index/", serveIndex(cfg.Token, cfg.RootDir))

	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf(`
Listening on: %s
Serving from: %s
Default file: %s
`, ingress, cfg.RootDir, cfg.Index)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}
