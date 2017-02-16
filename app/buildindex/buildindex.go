package buildindex

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Config allows tweaking of how a directory walk will behave.
type Config struct {
	RootDir        string   // Which directory should scanning begin at
	MaxDepth       int      // How deep should we go before giving up; -1 indicates no limit
	FollowSymlinks bool     // when we hit a symlinked directory should we walk down it
	Exclude        []string // which directories should be ignored when building the index
	Check          MatchFn  // check that will be made on each file to determine inclusion
}

// MatchFn is called or each visited path and should return whether this file is included
// in the collected FileSet or not.
type MatchFn func(string) bool

// OfDir builds a list files whose names match some function. All paths will be
// returned relative to 'dir'.
func OfDir(cfg Config) FileSet {
	prepareExcludedDirs(cfg)
	if cfg.Check == nil {
		cfg.Check = MatchAll
	}
	return ofDir(cfg.RootDir, 0, cfg)
}

func MatchAll(name string) bool {
	return true
}

// MatchExt includes any file that has an extension specified in ext.
func MatchExt(ext ...string) MatchFn {
	return func(name string) bool {
		for _, e := range ext {
			if strings.HasSuffix(name, e) {
				return true
			}
		}
		return false
	}
}

// prepareExcludedDirs cleans all excluded directories of unnecessary cruft in
// their path and modifies the provided config directly
func prepareExcludedDirs(cfg Config) {
	for i, ed := range cfg.Exclude {
		cfg.Exclude[i] = filepath.Clean(filepath.Join(cfg.RootDir, ed))
	}
}

func ofDir(
	curDir string,
	curDepth int,
	cfg Config,
) FileSet {
	// TODO: clean up so not one huge function
	// TODO: best effort to resolve symlinks and only descend the path once
	if cfg.MaxDepth != -1 && curDepth > cfg.MaxDepth {
		return FileSet{}
	}

	fptr, err := os.Open(curDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening %s: %s\n", curDir, err)
		fptr.Close()
		return FileSet{}
	}

	fstat, err := fptr.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not stat %s: %s\n", curDir, err)
		fptr.Close()
		return FileSet{}
	}

	if !fstat.IsDir() {
		fmt.Fprintf(os.Stderr, "Could not index file as directory: %s\n", curDir)
		fptr.Close()
		return FileSet{}
	}

	fileinfos, err := fptr.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hit error indexding %s: %s\n", curDir, err)
		return FileSet{}
	}

	// we're done with this dir now
	fptr.Close()

	foundfiles := FileSet{}
	for _, fi := range fileinfos {
		fipath := path.Join(curDir, fi.Name())
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

			if skipDir {
				continue
			}

			dirfiles := ofDir(fipath, curDepth+1, cfg)
			if len(dirfiles) != 0 {
				foundfiles = append(foundfiles, dirfiles...)
			}
		} else {
			if cfg.Check(fi.Name()) {
				fentry := newFile(curDir, fi.Name())
				foundfiles = append(foundfiles, fentry)
			}
		}
	}

	if len(foundfiles) != 0 {
		dentry := newDir(curDir)
		foundfiles = append([]FileEntry{dentry}, foundfiles...)
	}

	return foundfiles
}
