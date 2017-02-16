// Package scanner runs buildindex.OfDir asynchronously until killed with a
// provided config. The results and a timestamp of the last run are maintained
// so we know liveness of the scan.
//
// TODO: currently this uses buildindex to walk a tree and collect a list of
// files. Ideally we would use fsnotify or similar to just watch for changes
// but there isn't an OOTB solution for a recursive cross platform directory
// watch via FS notifications and I don't feel like dealing with it now.
package scanner

import (
	"fmt"
	"sync"
	"time"

	"github.com/falun/markup/app/buildindex"
)

var (
	mtx          = sync.RWMutex{}
	lastResults  buildindex.FileSet
	scanComplete time.Time
	timer        *time.Timer
)

// Start starts reading a directory with the specified frequency and
// excludes a set of subdirectories.
//
// Will return an error if called more than once.
func Start(dir string, freq time.Duration, exclude []string) error {
	cfg := buildindex.Config{
		RootDir:        dir,
		MaxDepth:       -1,
		FollowSymlinks: true,
		Exclude:        exclude,
	}

	return StartCfg(cfg, freq, updateResults)
}

// StartCfg allows a custom configuration for buildindex and the callback
// upon scan completion to be specified. The scan will be initiated after
// freq time has passed since the completaion of the last scan. If 0 is
// passed as freq then an index will be built once and not rescheduled.
//
// Will return an error if called more than once.
func StartCfg(
	cfg buildindex.Config,
	freq time.Duration,
	updateFn func(buildindex.FileSet),
) error {
	mtx.Lock()
	defer mtx.Unlock()

	if timer == nil {
		// create a new timer and set it to fire immediately
		timer = time.NewTimer(0)
		go scan(cfg, freq, updateFn)
	} else {
		return fmt.Errorf("scanner.StartCfg may not be called more than once")
	}

	return nil
}

func Results() (buildindex.FileSet, time.Time) {
	var result buildindex.FileSet
	var ts time.Time

	mtx.RLock()
	defer mtx.RUnlock()

	if lastResults != nil {
		result = lastResults
		ts = scanComplete
	}

	return result, ts
}

func updateResults(fs buildindex.FileSet) {
	mtx.Lock()
	defer mtx.Unlock()

	lastResults = fs
	scanComplete = time.Now()
}

func scan(
	cfg buildindex.Config,
	freq time.Duration,
	updateFn func(buildindex.FileSet),
) {
	// wait for timer to expire
	<-timer.C

	// do a thing
	updateFn(buildindex.OfDir(cfg))

	if freq == 0 {
		return
	}

	// ...and then do it again
	timer.Reset(freq)
	go scan(cfg, freq, updateFn)
}
