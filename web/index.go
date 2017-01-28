package web

import (
	"sort"
	"strings"
)

type IndexView struct {
	Root        string
	Children    []ChildView
	IndexToken  string
	BrowseToken string
}

type ChildView struct {
	Name  string
	IsDir bool
}

var _ sort.Interface = IndexView{}

func NewIndexView(rootDir, indexToken, browseToken string) IndexView {
	return IndexView{rootDir, nil, indexToken, browseToken}
}

func (iv *IndexView) AddChild(name string, isDir bool) {
	iv.Children = append(iv.Children, ChildView{name, isDir})
}

func (iv *IndexView) Sort() {
	sort.Sort(iv)
}

func (c IndexView) Len() int { return len(c.Children) }

func (c IndexView) Swap(i, j int) {
	c.Children[i], c.Children[j] = c.Children[j], c.Children[i]
}

func (c IndexView) Less(i, j int) bool {
	id := c.Children[i].IsDir
	jd := c.Children[j].IsDir

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
