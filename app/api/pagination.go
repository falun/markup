package api

const (
	DefaultPageSize = 25
	FlagPageSize    = "page-size"
	FlagPageNumber  = "page"
)

type Pagination struct {
	PageNumber int
	PageSize   int
}

// NewPagination extracts a Pagination object from query args. If no args are
// present it uses the package defaults.
func NewPagination(a Args) Pagination {
	return Pagination{
		PageNumber: a.GetInt(FlagPageNumber, 0),
		PageSize:   a.GetInt(FlagPageSize, DefaultPageSize),
	}
}
