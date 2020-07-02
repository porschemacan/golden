package golden

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"time"
)

const KEY_SPAN = "SPAN"
const ERROR_NO = "ERRORNO"
const PROCTIME = "PROCTIME"
const TRACEID = "TRACEID"
const SPANID = "SPANID"

type HttpContext struct {
	context *gin.Context
	options *ServerOptions
	client  *resty.Client
}

func NewEmptyHttpContext() *HttpContext {
	return &HttpContext{
		context: &gin.Context{},
		client:  resty.New(),
	}
}

func (c *HttpContext) NewRequest(ctx context.Context) *HttpRequest {
	return NewHttpRequestWithHttpContext(ctx, c)
}

func (c *HttpContext) RawRequest() *http.Request {
	return c.context.Request
}

func (c *HttpContext) SetRawRequest(r *http.Request) *http.Request {
	old := c.context.Request
	c.context.Request = r
	return old
}

func (c *HttpContext) ResponseStatus() int {
	return c.context.Writer.Status()
}

func (c *HttpContext) ResponseSize() int {
	return c.context.Writer.Size()
}

func (c *HttpContext) Next() {
	c.context.Next()
}

func (c *HttpContext) JSON(code int, obj interface{}) {
	c.context.JSON(code, obj)
}

func (c *HttpContext) String(code int, format string, values ...interface{}) {
	c.context.String(code, format, values...)
}

func (c *HttpContext) ClientIP() string {
	return c.context.ClientIP()
}

// It parses the request's body as JSON
func (c *HttpContext) TryBindJSON(obj interface{}) error {
	return c.context.ShouldBindWith(obj, binding.JSON)
}

// It parses the request's body as protobuf
func (c *HttpContext) TryBindProtobuf(obj interface{}) error {
	return c.context.ShouldBindWith(obj, binding.ProtoBuf)
}

func (c *HttpContext) GetRawData() ([]byte, error) {
	return c.context.GetRawData()
}

func (c *HttpContext) LogKV(k string, v interface{}) {
	span := c.GetSpan()
	span.LogKV(k, v)
}

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (c *HttpContext) Param(key string) string {
	return c.context.Param(key)
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func (c *HttpContext) Query(key string) (string, bool) {
	return c.context.GetQuery(key)
}

func (c *HttpContext) DefaultQuery(key, defaultValue string) string {
	return c.context.DefaultQuery(key, defaultValue)
}

func (c *HttpContext) RequestHeader(key string) string {
	return c.context.Request.Header.Get(key)
}

func (c *HttpContext) Set(key string, value interface{}) {
	c.context.Set(key, value)
}

func (c *HttpContext) GetSpan() opentracing.Span {
	v, exist := c.context.Get(KEY_SPAN)
	if !exist {
		return nil
	}
	return v.(opentracing.Span)
}

func (c *HttpContext) Get(key string) (value interface{}, exists bool) {
	return c.context.Get(key)
}

func (c *HttpContext) GetTraceId() string {
	return c.GetString(TRACEID)
}

func (c *HttpContext) GetSpanId() string {
	return c.GetString(SPANID)
}

func (c *HttpContext) GetString(key string) (s string) {
	return c.context.GetString(key)
}

func (c *HttpContext) GetBool(key string) (b bool) {
	return c.context.GetBool(key)
}

func (c *HttpContext) GetInt(key string) (i int) {
	return c.context.GetInt(key)
}

func (c *HttpContext) GetInt64(key string) (i64 int64) {
	return c.context.GetInt64(key)
}

func (c *HttpContext) GetFloat64(key string) (f64 float64) {
	return c.context.GetFloat64(key)
}

func (c *HttpContext) GetTime(key string) (t time.Time) {
	return c.context.GetTime(key)
}

func (c *HttpContext) GetDuration(key string) (d time.Duration) {
	return c.context.GetDuration(key)
}

func (c *HttpContext) GetStringSlice(key string) (ss []string) {
	return c.context.GetStringSlice(key)
}

func (c *HttpContext) GetStringMap(key string) (sm map[string]interface{}) {
	return c.context.GetStringMap(key)
}

func (c *HttpContext) GetStringMapString(key string) (sms map[string]string) {
	return c.context.GetStringMapString(key)
}

func (c *HttpContext) GetStringMapStringSlice(key string) (smss map[string][]string) {
	return c.context.GetStringMapStringSlice(key)
}
