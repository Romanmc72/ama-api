package constants

// Params for path parameter query strings in API requests.
const (
	// When you want to limit the number of returned results.
	LimitParam = "limit"
	// The default limit to limit results.
	DefaultLimit = 0
	// When paginating, the final ID that was seen in the previous request.
	// This will default to a pseudo random value if not specified.
	FinalIdParam = "finalId"
	// The parameter to use for filtering based on tags.
	TagParam = "tags"
	// The default tag parameter value (no tag filters).
	DefaultTag = ""
)
