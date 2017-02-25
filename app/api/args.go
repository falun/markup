package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Args are extracted from the query parameters of a request and allow access
// to relevant parameters through convenience methods.
type Args map[string][]string

// NewArgs extracts query parameters into an Args struct.
func NewArgs(r *http.Request) (Args, error) {
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

// IsSet returns whether a given parameter is in the query string. I parameter
// present with no value (e.g. ?some-name) is considered set.
func (a Args) IsSet(p string) bool {
	_, ok := a[p]
	return ok
}

// Get returns the first value of a parameter or "" if no value was set.
func (a Args) Get(p string) string {
	if a.IsSet(p) {
		return a[p][0]
	}
	return ""
}

// GetAll returns all values for a parameter if it was specified more than
// once.
func (a Args) GetAll(p string) []string {
	return a[p]
}

// GetInt returns a parameter as an int or some default value if it is not set
// or if it fails to convert cleanly.
func (a Args) GetInt(p string, def int) int {
	if !a.IsSet(p) {
		return def
	}

	i, err := strconv.Atoi(a.Get(p))
	if err != nil {
		return def
	}
	return i
}
