package buildindex

import (
	"path/filepath"
	"sort"
)

// FileEntry represents a single item processed during a directory walk.
type FileEntry struct {
	Root string // Root is the location of a file or the current directory
	Name string // Valid only if Dir is true, the name of the file
	Path string // If not a directory this will be the full path to the file
	Dir  bool   // Set to true if this entry represents a directory
}

type FileSet []FileEntry

// ByFilepath implements a sort by file directory & file name.
type ByFilepath []FileEntry

func (fs FileSet) Sort() {
	sort.Sort(ByFilepath(fs))
}

func (f ByFilepath) Len() int {
	return len(f)
}

func (f ByFilepath) Less(i, j int) bool {
	if f[i].Root != f[j].Root {
		return f[i].Root < f[j].Root
	}

	if f[i].Name >= f[j].Name {
		return false
	}

	return true
}

func (f ByFilepath) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func newFile(dir, name string) FileEntry {
	return FileEntry{Root: dir, Name: name, Path: filepath.Join(dir, name)}
}

func newDir(dir string) FileEntry {
	return FileEntry{Root: dir, Dir: true, Path: filepath.Join(dir, "")}
}
