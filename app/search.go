package app

import (
	"net/http"
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
	FlagTerm = "term"

	// FlagFileType can be specified to match only a certain file type; it may
	// be passed multiple times or have multiple values comma-separated.
	FlagFileType = "ft"

	// FlagCaseSensitive is a query parameter that can be passed to make search
	// case sensitive
	FlagCaseSensitive = "case-sensitive"

	// FlagRegexp is a query parameter that can be passed which forces evaluation
	// of the search term as a regexp instead of substring match.
	FlagRegexp = "regexp"

	// FlagIncludeDir, when set, expands the search term match to include the
	// containing path as well as just the file name.
	FlagIncludeDir = "include-dir"
)

// SearchResults is the payload that will be sent back in the case of a
// successful search request.
type SearchResults struct {
	SearchTerm  string   `json:"search_term"`
	FileTypes   []string `json:"file_types"`
	TimeStamp   int      `json:"last_scan_ms"`
	Matches     []string `json:"match"`
	ResultCount int      `json:"count"`
}

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
	result := api.Results{}

	ta, err := api.NewArgs(r)
	if err != nil {
		result.Err = err
		result.RenderTo(rw)
		return
	}

	args := Args(ta)
	pg := api.NewPagination(api.Args(args))

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
		0,
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

	sr.ResultCount = len(sr.Matches)

	offset := 0
	if pg.PageSize > 0 {
		offset = pg.PageNumber * pg.PageSize
		if offset >= len(sr.Matches) {
			sr.Matches = sr.Matches[0:0]
		} else {
			sr.Matches = sr.Matches[offset:]
			if len(sr.Matches) > pg.PageSize {
				sr.Matches = sr.Matches[:pg.PageSize]
			}
		}
	}

	result.Results = sr
	result.RenderTo(rw)
}

type Args api.Args

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
			if fte == "" {
				continue
			}
			t := strings.ToLower(fte)
			if t[0] != '.' {
				t = "." + t
			}
			rm[t] = true
			ra = append(ra, t)
		}
	}

	return ra, rm, len(ra) > 0
}
