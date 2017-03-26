package render

// Response is the return value of the endpoint.
type Response struct {
}

// DefaultResponse provides a default response object by best effort.
func DefaultResponse() *Response {
	return &Response{}
}
