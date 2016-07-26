package buildindex

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Config allows tweaking of how directory walk behaves.
type Config struct {
	MaxDepth       int      // How deep should we go before giving up; -1 indicates no limit
	FollowSymlinks bool     // when we hit a symlinked directory should we walk down it?
	Exclude        []string // which directories should be ignored when building the index
}

func prepareExcludedDirs(cfg Config, base string) {
	for i, ed := range cfg.Exclude {
		cfg.Exclude[i] = filepath.Clean(filepath.Join(base, ed))
	}
}

// OfDir builds a list files whose names match some function. All paths will be
// returned relative to 'dir'.
func OfDir(dir string, cfg Config, matchfn func(string) bool) []FileEntry {
	fmt.Printf("Indexing %s\n", dir)
	prepareExcludedDirs(cfg, dir)
	return ofDir(dir, 0, cfg, matchfn)
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

func ofDir(
	dir string,
	curDepth int,
	cfg Config,
	matchfn func(string) bool,
) []FileEntry {
	if cfg.MaxDepth != -1 && curDepth > cfg.MaxDepth {
		return []FileEntry{}
	}

	fptr, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening %s: %s\n", dir, err)
		fptr.Close()
		return []FileEntry{}
	}

	fstat, err := fptr.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not stat %s: %s\n", dir, err)
		fptr.Close()
		return []FileEntry{}
	}

	if !fstat.IsDir() {
		fmt.Fprintf(os.Stderr, "Could not index file as directory: %s\n", dir)
		fptr.Close()
		return []FileEntry{}
	}

	fileinfos, err := fptr.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hit error indexding %s: %s\n", dir, err)
		return []FileEntry{}
	}

	// we're done with this dir now
	fptr.Close()

	foundfiles := []FileEntry{}
	for _, fi := range fileinfos {
		fipath := path.Join(dir, fi.Name())
		if fi.Mode()&os.ModeSymlink != 0 {
			original := fipath
			fipath, err := filepath.EvalSymlinks(fipath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not resolve symlink %s: %s\n", original, err)
				continue
			}

			fi, err = os.Stat(fipath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not stat %s, resolved from symlink %s: %s\n", fipath, original, err)
				continue
			}

			if !cfg.FollowSymlinks && fi.IsDir() {
				fmt.Fprintf(os.Stderr, "Not following symlink %s\n", fipath)
				continue
			}
		}

		if fi.IsDir() {
			skipDir := false
			for _, ed := range cfg.Exclude {
				if ed == fipath {
					skipDir = true
				}
			}

			if skipDir { continue }

			dirfiles := ofDir(fipath, curDepth+1, cfg, matchfn)
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

	return foundfiles
}
