package golden

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type HttpResponse struct {
	impl *resty.Response
}

// Body method returns HTTP response as []byte array for the executed request.
// Note: `Response.Body` might be nil, if `Request.SetOutput` is used.
func (r *HttpResponse) Body() []byte {
	return r.impl.Body()
}

// Status method returns the HTTP status string for the executed request.
//	Example: 200 OK
func (r *HttpResponse) Status() string {
	return r.impl.Status()
}

// StatusCode method returns the HTTP status code for the executed request.
//	Example: 200
func (r *HttpResponse) StatusCode() int {
	return r.impl.StatusCode()
}

// Result method returns the response value as an object if it has one
func (r *HttpResponse) Result() interface{} {
	return r.impl.Result()
}

// Error method returns the error object if it has one
func (r *HttpResponse) Error() interface{} {
	return r.impl.Error()
}

// Header method returns the response headers
func (r *HttpResponse) Header() http.Header {
	return r.impl.Header()
}

// Cookies method to access all the response cookies
func (r *HttpResponse) Cookies() []*http.Cookie {
	return r.impl.Cookies()
}

// String method returns the body of the server response as String.
func (r *HttpResponse) String() string {
	return r.impl.String()
}

// Time method returns the time of HTTP response time that from request we sent and received a request.
// See `response.ReceivedAt` to know when client recevied response and see `response.Request.Time` to know
// when client sent a request.
func (r *HttpResponse) Time() time.Duration {
	return r.impl.Time()
}

// ReceivedAt method returns when response got recevied from server for the request.
func (r *HttpResponse) ReceivedAt() time.Time {
	return r.impl.ReceivedAt()
}

// Size method returns the HTTP response size in bytes. Ya, you can relay on HTTP `Content-Length` header,
// however it won't be good for chucked transfer/compressed response. Since Resty calculates response size
// at the client end. You will get actual size of the http response.
func (r *HttpResponse) Size() int64 {
	return r.impl.Size()
}
