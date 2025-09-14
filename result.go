package hapi

// Result represents the parsed query parameters including filters, sorting, and pagination.
type Result struct {
	Filters Filters // Collection of filter conditions
	Sorts   Sorts   // Sorting configuration
	Page    int     // Current page number (1-based)
	PerPage int     // Number of items per page
}
