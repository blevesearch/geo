module github.com/blevesearch/geo

// This declares what language features we use and may be updated freely up to "oldstable".
// Update .github/workflows/go.yml when bumping this version.
go 1.24

toolchain go1.24.0

require (
	github.com/blevesearch/bleve_index_api v1.4.1-0.20260729060817-8e56340f2a7e
	github.com/google/go-cmp v0.7.0
	github.com/json-iterator/go v0.0.0-20171115153421-f7279a603ede
)

require (
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
)
