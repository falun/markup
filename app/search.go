package app

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/falun/markup/app/api"
	"github.com/falun/markup/app/scanner"
)

const (
	// FlagTerm is the query param that may be set to specifi what is being
	// searched for; if no term is present all files will be matched.
	FlagTerm          = "term"

	// FlagFileType can be specified to match only a certain file type; it may
	// be passed multiple times or have multiple values comma-separated.
	FlagFileType      = "ft"

	// FlagCaseSensitive is a query parameter that can be passed to make search
	// case sensitive
	FlagCaseSensitive = "case-sensitive"

	// FlagRegexp is a query parameter that can be passed which forces evaluation
	// of the search term as a regexp instead of substring match.
	FlagRegexp        = "regexp"

	// FlagIncludeDir, when set, expands the search term match to include the
	// containing path as well as just the file name.
	FlagIncludeDir    = "include-dir"
)

// SearchResults is the payload that will be sent back in the case of a
// successful search request.
type SearchResults struct {
	SearchTerm string   `json:"search_term"`
	FileTypes  []string `json:"file_types"`
	TimeStamp  int      `json:"last_scan_ms"`
	Matches    []string `json:"match"`
}

// Args are extracted from the query parameters of a request and allow access
// to relevant parameters through convenience methods.
type Args map[string][]string

type searchHandlerS struct {
	rootDir string
}

func searchHandler(root string) searchHandlerS {
	return searchHandlerS{
		rootDir: root,
	}
}

var whitespaceRE = regexp.MustCompile("\\s+")

func (sh searchHandlerS) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	args, err := getArgs(r)
	result := api.Results{}

	if err != nil {
		result.Err = err
		result.RenderTo(rw)
		return
	}

	t := args.Term()
	t = string(whitespaceRE.ReplaceAll([]byte(t), []byte(".*")))

	if !args.CaseSensitive() {
		t = strings.ToLower(t)
	}

	re, err := regexp.Compile(t)
	if err != nil {
		result.Err = err
		result.RenderTo(rw)
		return
	}

	fs, ts := scanner.Results()
	fta, ftm, checkFT := args.FileTypes()

	sr := SearchResults{
		t,
		fta,
		int(ts.UnixNano() / int64(time.Millisecond)),
		nil,
	}

	for _, f := range fs {
		if checkFT && !ftm[strings.ToLower(path.Ext(f.Name))] {
			continue
		}

		checkStr := f.Name
		if args.IncludeDir() {
			checkStr = f.Path
		} else if f.Dir {
			continue
		}

		if !args.CaseSensitive() {
			checkStr = strings.ToLower(checkStr)
		}

		if re.Match([]byte(checkStr)) {
			p, err := filepath.Rel(sh.rootDir, f.Path)
			if err != nil {
				result.Err = err
				result.RenderTo(rw)
				return
			}
			sr.Matches = append(sr.Matches, p)
		}
	}
	result.Results = sr

	result.RenderTo(rw)
}

func (a Args) CaseSensitive() bool {
	_, ok := a[FlagCaseSensitive]
	return ok
}

func (a Args) Regexp() bool {
	_, ok := a[FlagRegexp]
	return ok
}

func (a Args) IncludeDir() bool {
	_, ok := a[FlagIncludeDir]
	return ok
}

func (a Args) Term() string {
	t := a[FlagTerm]
	if len(t) == 0 {
		return ""
	}
	return t[0]
}

func (a Args) FileTypes() ([]string, map[string]bool, bool) {
	fts, ok := a[FlagFileType]
	if !ok {
		return nil, nil, false
	}

	ra := []string{}
	rm := map[string]bool{}

	for _, ft := range fts {
		fta := strings.Split(ft, ",")
		for _, fte := range fta {
			t := strings.ToLower(fte)
			if t[0] != '.' {
				t = "." + t
			}
			rm[t] = true
			ra = append(ra, t)
		}
	}
	return ra, rm, true
}

func getArgs(r *http.Request) (Args, error) {
	args := Args{}

	q := r.URL.RawQuery
	parts := strings.Split(q, "&")

	for _, p := range parts {
		k := ""
		v := ""
		var err error

		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 0 {
			continue
		}

		if len(kv) > 0 {
			k, err = url.QueryUnescape(kv[0])
			if err != nil {
				return nil, err
			}
		}

		if len(kv) > 1 {
			v, err = url.QueryUnescape(kv[1])
			if err != nil {
				return nil, err
			}
		}

		args[k] = append(args[k], v)
	}
	return args, nil
}
