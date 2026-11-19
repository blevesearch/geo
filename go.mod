module github.com/blevesearch/geo

// This declares what language features we use and may be updated freely up to "oldstable".
// Update .github/workflows/go.yml when bumping this version.
go 1.24

toolchain go1.24.0

require (
	github.com/blevesearch/bleve_index_api v1.4.1
	github.com/google/go-cmp v0.7.0
)
