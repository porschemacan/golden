package golden

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
)

type HttpRequest struct {
	impl *resty.Request
}

func newHttpRequestFromContextAndSpan(ctx context.Context, span opentracing.Span, request *resty.Request) *HttpRequest {
	ctxWithSpan := opentracing.ContextWithSpan(ctx, span)
	request.SetContext(ctxWithSpan)
	return &HttpRequest{
		impl: request,
	}
}

func NewHttpRequestWithHttpContext(ctx context.Context, c *HttpContext) *HttpRequest {
	return newHttpRequestFromContextAndSpan(ctx, c.GetSpan(), c.client.R())
}

func NewHttpRequest(ctx context.Context) *HttpRequest {
	return newHttpRequestFromContextAndSpan(ctx, opentracing.StartSpan("clientSpan"), resty.New().R())
}

func (r *HttpRequest) SetImpl(impl *resty.Request) {
	r.impl = impl
}

//			SetHeader("Content-Type", "application/json").
//			SetHeader("Accept", "application/json")
func (r *HttpRequest) SetHeader(header, value string) *HttpRequest {
	r.impl.SetHeader(header, value)
	return r
}

//			SetHeaders(map[string]string{
//				"Content-Type": "application/json",
//				"Accept": "application/json",
//			})
func (r *HttpRequest) SetHeaders(headers map[string]string) *HttpRequest {
	r.impl.SetHeaders(headers)
	return r
}

//			SetQueryParam("search", "kitchen papers").
//			SetQueryParam("size", "large")
func (r *HttpRequest) SetQueryParam(param, value string) *HttpRequest {
	r.impl.SetQueryParam(param, value)
	return r
}

//			SetQueryParams(map[string]string{
//				"search": "kitchen papers",
//				"size": "large",
//			})
func (r *HttpRequest) SetQueryParams(params map[string]string) *HttpRequest {
	r.impl.SetQueryParams(params)
	return r
}

// 			SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")
func (r *HttpRequest) SetQueryString(query string) *HttpRequest {
	r.impl.SetQueryString(query)
	return r
}

// 			SetFormData(map[string]string{
//				"access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
//				"user_id": "3455454545",
//			})
func (r *HttpRequest) SetFormData(data map[string]string) *HttpRequest {
	r.impl.SetFormData(data)
	return r
}

// SetBody method sets the request body for the request. It supports various realtime needs as easy.
// We can say its quite handy or powerful. Supported request body data types is `string`,
// `[]byte`, `struct`, `map`, `slice` and `io.Reader`. Body value can be pointer or non-pointer.
func (r *HttpRequest) SetBody(body interface{}) *HttpRequest {
	r.impl.SetBody(body)
	return r
}

// SetResult method is to register the response `Result` object for automatic unmarshalling in the RESTful mode
// if response status code is between 200 and 299 and content type either JSON or XML.
//
// Note: Result object can be pointer or non-pointer.
//		request.SetResult(&AuthToken{})
//		// OR
//		request.SetResult(AuthToken{})
//
// Accessing a result value
//		response.Result().(*AuthToken)
//
func (r *HttpRequest) SetResult(res interface{}) *HttpRequest {
	r.impl.SetResult(res)
	return r
}

// SetError method is to register the request `Error` object for automatic unmarshalling in the RESTful mode
// if response status code is greater than 399 and content type either JSON or XML.
//
// Note: Error object can be pointer or non-pointer.
// 		request.SetError(&AuthError{})
//		// OR
//		request.SetError(AuthError{})
//
// Accessing a error value
//		response.Error().(*AuthError)
//
func (r *HttpRequest) SetError(err interface{}) *HttpRequest {
	r.impl.SetError(err)
	return r
}

// SetBasicAuth method sets the basic authentication header in the current HTTP request.
// For Header example:
//		Authorization: Basic <base64-encoded-value>
//
// To set the header for username "go-resty" and password "welcome"
// 		request.SetBasicAuth("go-resty", "welcome")
//
// This method overrides the credentials set by method `resty.SetBasicAuth`.
//
func (r *HttpRequest) SetBasicAuth(username, password string) *HttpRequest {
	r.impl.SetBasicAuth(username, password)
	return r
}

// SetAuthToken method sets bearer auth token header in the current HTTP request. Header example:
// 		Authorization: Bearer <auth-token-value-comes-here>
//
// Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
//
// 		request.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
//
// This method overrides the Auth token set by method `resty.SetAuthToken`.
//
func (r *HttpRequest) SetAuthToken(token string) *HttpRequest {
	r.impl.SetAuthToken(token)
	return r
}

// Get method does GET HTTP request. It's defined in section 4.3.1 of RFC7231.
func (r *HttpRequest) Get(url string) (*HttpResponse, error) {
	resp, err := r.impl.Get(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Head method does HEAD HTTP request. It's defined in section 4.3.2 of RFC7231.
func (r *HttpRequest) Head(url string) (*HttpResponse, error) {
	resp, err := r.impl.Head(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Post method does POST HTTP request. It's defined in section 4.3.3 of RFC7231.
func (r *HttpRequest) Post(url string) (*HttpResponse, error) {
	resp, err := r.impl.Post(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Put method does PUT HTTP request. It's defined in section 4.3.4 of RFC7231.
func (r *HttpRequest) Put(url string) (*HttpResponse, error) {
	resp, err := r.impl.Put(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Delete method does DELETE HTTP request. It's defined in section 4.3.5 of RFC7231.
func (r *HttpRequest) Delete(url string) (*HttpResponse, error) {
	resp, err := r.impl.Delete(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Options method does OPTIONS HTTP request. It's defined in section 4.3.7 of RFC7231.
func (r *HttpRequest) Options(url string) (*HttpResponse, error) {
	resp, err := r.impl.Options(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}

// Patch method does PATCH HTTP request. It's defined in section 2 of RFC5789.
func (r *HttpRequest) Patch(url string) (*HttpResponse, error) {
	resp, err := r.impl.Patch(url)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		impl: resp,
	}, nil
}
