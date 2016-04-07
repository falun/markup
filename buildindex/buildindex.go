package buildindex

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type FileEntry struct {
	Root string
	Name string
	Dir  bool
}

type ByFilepath []FileEntry

func (f ByFilepath) Len() int {
	return len(f)
}

func (f ByFilepath) Less(i, j int) bool {
	if f[i].Root > f[j].Root {
		return false
	}

	if (f[i].Root == f[j].Root) && (f[i].Name >= f[j].Name) {
		return false
	}

	return true
}

func (f ByFilepath) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func MatchAll(name string) bool {
	return true
}

func MatchExt(ext ...string) func(string) bool {
	return func(name string) bool {
		for _, e := range ext {
			if strings.HasSuffix(name, e) {
				return true
			}
		}
		return false
	}
}

// OfDir builds a list files whose names match some function. All paths will be
// returned relative to 'dir'.
//
// TODO: allow configuration for do/do-not follow symlinks
func OfDir(dir string, depth int, matchfn func(string) bool) ([]FileEntry, error) {
	return ofDir(dir, 0, depth, matchfn)
}

// TODO: ugh. filepath.Walk is what I wanted; more reasons that coding on the
// train is a pita. Anyway, replace with that eventually.
func ofDir(
	dir string,
	curDepth, maxDepth int,
	matchfn func(string) bool,
) ([]FileEntry, error) {
	if maxDepth != 0 && curDepth > maxDepth {
		return []FileEntry{}, nil
	}

	fptr, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening %s: %s\n", dir, err)
		fptr.Close()
		return []FileEntry{}, nil
	}

	fstat, err := fptr.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not stat %s: %s\n", dir, err)
		fptr.Close()
		return []FileEntry{}, nil
	}

	if !fstat.IsDir() {
		fmt.Fprintf(os.Stderr, "Could not index file as directory: %s\n", dir)
		fptr.Close()
		return []FileEntry{}, nil
	}

	fileinfos, err := fptr.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hit error indexding %s: %s\n", dir, err)
		return []FileEntry{}, nil
	}

	// we're done with this dir now
	fptr.Close()

	foundfiles := []FileEntry{}
	for _, fi := range fileinfos {
		fipath := path.Join(dir, fi.Name())
		if fi.IsDir() {
			dirfiles, _ := ofDir(fipath, curDepth+1, maxDepth, matchfn)
			if len(dirfiles) != 0 {
				foundfiles = append(foundfiles, dirfiles...)
			}
		} else {
			if matchfn(fi.Name()) {
				fentry := FileEntry{Root: dir, Name: fi.Name(), Dir: false}
				foundfiles = append(foundfiles, fentry)
			}
		}
	}

	if len(foundfiles) != 0 {
		dentry := FileEntry{Root: dir, Name: "", Dir: true}
		foundfiles = append([]FileEntry{dentry}, foundfiles...)
	}

	return foundfiles, nil
}
