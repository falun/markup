package app

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/falun/markup/web"
)

type indexHandler struct {
	web.Renderer

	indexToken  string
	browseToken string
	rootDir     string
	excluding   []string
}

// TODO: better error handling
func (h indexHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	urlPath := path.Clean(r.URL.Path[len(h.indexToken)+1:])

	// gets the url path cleaned of the index token
	urlPath = urlPath[1:] + "/"
	// constructs a file path relative to markup root
	fp := path.Join(h.rootDir, urlPath)

	// enumerate children of the requested directory
	infos, err := ioutil.ReadDir(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	view := web.NewIndexView(urlPath, h.indexToken, h.browseToken)

	for _, i := range infos {
		// construct the file path of a given child
		ifp := path.Join(fp, i.Name())
		isDir := i.IsDir()

		// if it's a symlink we can't trust the supplied info's isDir
		symlink := 0 != int(i.Mode()&os.ModeSymlink)
		if symlink {
			// TODO: should readlink instead of relying on path.Join?
			info, err := os.Stat(ifp)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			isDir = info.IsDir()
		}

		view.AddChild(i.Name(), isDir)
	}

	view.Sort()

	output, err := h.Render(web.IndexTmpl, view)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Write([]byte(output))
}
