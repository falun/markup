package app

import (
	"fmt"
	"net/http"
	"path"
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

const frame = `<html>
	<body>
		<ul>
			%s
		</ul>
	</body>
	</html>`

func serveIndex(indexToken, browseToken, root string, excluding []string) http.HandlerFunc {
	root = path.Clean(root)
	rootLen := len(root)

	cfg := buildindex.Config{
		MaxDepth:       -1,
		FollowSymlinks: true,
		Exclude:        excluding}
	fileset := buildindex.OfDir(root, cfg, buildindex.MatchExt(".md"))
	sort.Sort(buildindex.ByFilepath(fileset))

	mkEntry := func(fe buildindex.FileEntry) string {
		relRoot := fe.Root[rootLen:]

		if fe.Dir {
			if relRoot == "" {
				return ""
			}
			return fmt.Sprintf(
				`<li><a href="/%s%s">%s</a>`,
				indexToken,
				relRoot,
				relRoot,
			)
		}

		return fmt.Sprintf(
			"<li>%s/<a href=\"%s\">%s</a></li>\n",
			relRoot,
			fmt.Sprintf("/%s%s/%s", browseToken, relRoot, fe.Name),
			fe.Name,
		)
	}

	contents := ""
	for _, e := range fileset {
		contents += mkEntry(e)
	}

	body := fmt.Sprintf(frame, contents)

	return func(rw http.ResponseWriter, r *http.Request) {
		// searching is the directory we're looking for; "." if empty
		searching := path.Clean(r.URL.Path[len(indexToken)+1:])

		if !(searching == "." || searching == "/") {
			contents := ""
			for _, c := range fileset {
				if c.Root[len(root):] == searching {
					contents += mkEntry(c)
				}
			}

			if contents == "" {
				http.Redirect(rw, r, fmt.Sprintf("/%s/", indexToken), http.StatusSeeOther)
				return
			}
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf(frame, contents)))
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(body))
		}
	}
}

var serveNotFound = http.NotFound

func serveError(rw http.ResponseWriter, r *http.Request, err error) {
	http.Error(rw, fmt.Sprintf("%v", err), http.StatusServiceUnavailable)
}
