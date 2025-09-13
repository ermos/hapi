package hapi

// Result represents the parsed query parameters including filters, sorting, and pagination.
type Result struct {
	Filters Filters // Collection of filter conditions
	Sort    Sort    // Sorting configuration
	Limit   int     // Maximum number of results to return (0 means unlimited)
	Offset  int     // Number of results to skip for pagination (1-based)
}
