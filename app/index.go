package app

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/falun/markup/assets"
)

type indexHandler struct {
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

	dir := Dir{urlPath, nil}
	addChild := func(name string, isDir bool) {
		dir.Children = append(dir.Children, child{name, isDir})
	}

	for _, i := range infos {
		// construct the file path of a given child
		ifp := path.Join(fp, i.Name())
		isDir := i.IsDir()
		// if it's a symlink we can't trust the supplied info's isDir
		symlink := 0 != int(i.Mode()&os.ModeSymlink)
		if symlink {
			info, err := os.Stat(ifp)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			isDir = info.IsDir()
		}

		addChild(i.Name(), isDir)
	}

	sort.Sort(dir)

	output, err := assets.RenderIndexPage(dir)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Write([]byte(output))
}

type Dir struct {
	Root     string
	Children []child
}

type child struct {
	Name string
	Dir  bool
}

var _ sort.Interface = Dir{}

func (c Dir) Len() int { return len(c.Children) }

func (c Dir) Swap(i, j int) {
	c.Children[i], c.Children[j] = c.Children[j], c.Children[i]
}

func (c Dir) Less(i, j int) bool {
	id := c.Children[i].Dir
	jd := c.Children[j].Dir

	if id != jd {
		// if they're not both dirs or !dirs then list then order the directory
		// first
		return id
	}

	// assuming both have the same dir status order by lowercased alpha
	in := strings.ToLower(c.Children[i].Name)
	jn := strings.ToLower(c.Children[j].Name)

	return in < jn
}
