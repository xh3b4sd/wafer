package search

// Request is the configuration for the endpoint.
type Request struct {
}

// DefaultRequest provides a default request object by best effort.
func DefaultRequest() Request {
	return Request{}
}
