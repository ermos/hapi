package hapi

// Default pagination values applied when Options leaves them unset.
const (
	defaultPerPage    = 10
	defaultMaxPerPage = 100
)

// Options defines configuration options for parsing and validating query parameters.
type Options struct {
	DefaultPerPage int
	MaxPerPage     int
	AllowedSorts   []string
	AllowedFilters []string
}

type OptionFunc func(*Options)

// NewOptions creates a new Options instance with default values.
func NewOptions(opts ...OptionFunc) *Options {
	options := &Options{
		DefaultPerPage: defaultPerPage,
		MaxPerPage:     defaultMaxPerPage,
		AllowedSorts:   []string{},
		AllowedFilters: []string{},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// WithDefaultPerPage sets the default number of items per page.
func WithDefaultPerPage(n int) OptionFunc {
	return func(o *Options) {
		o.DefaultPerPage = n
	}
}

// WithMaxPerPage sets the maximum number of items per page.
func WithMaxPerPage(n int) OptionFunc {
	return func(o *Options) {
		o.MaxPerPage = n
	}
}

// WithAllowedSorts sets the allowed sorting fields.
func WithAllowedSorts(sorts []string) OptionFunc {
	return func(o *Options) {
		o.AllowedSorts = sorts
	}
}

// WithAllowedFilters sets the allowed filtering fields.
func WithAllowedFilters(filters []string) OptionFunc {
	return func(o *Options) {
		o.AllowedFilters = filters
	}
}
