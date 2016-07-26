package app

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/falun/markup/buildindex"
)

func serveRoot(prefix, browseToken, index string) http.HandlerFunc {
	rootPage := fmt.Sprintf("/%s/%s", browseToken, index)
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == prefix {
			http.Redirect(rw, r, rootPage, http.StatusSeeOther)
		} else {
			serveNotFound(rw, r)
		}
	}
}

func serveIndex(token, root string, excluding []string) http.HandlerFunc {
	l := len(root)
	cfg := buildindex.Config{
		MaxDepth:       -1,
		FollowSymlinks: true,
		Exclude:        excluding}

	fileset := buildindex.OfDir(root, cfg, buildindex.MatchExt(".md"))
	sort.Sort(buildindex.ByFilepath(fileset))

	frame := `<html>
	<body>
		<ul>
			%s
		</ul>
	</body>
	</html>`

	mkEntry := func(fe buildindex.FileEntry) string {
		if fe.Dir {
			return ""
		}

		relRoot := fe.Root[l:]
		return fmt.Sprintf(
			"<li>%s/<a href=\"%s\">%s</a></li>\n",
			relRoot,
			fmt.Sprintf("../%s/%s/%s", token, relRoot, fe.Name),
			fe.Name,
		)
	}

	contents := ""
	for _, e := range fileset {
		contents += mkEntry(e)
	}

	body := fmt.Sprintf(frame, contents)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(body))
	}
}

var serveNotFound = http.NotFound

func serveError(rw http.ResponseWriter, r *http.Request, err error) {
	http.Error(rw, fmt.Sprintf("%v", err), http.StatusServiceUnavailable)
}
