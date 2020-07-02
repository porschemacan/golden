package golden

import (
	gocontext "context"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
	"time"
)

var ContextKey = struct{}{}

func init() {
	opts := basictracer.DefaultOptions()
	opts.Recorder = &SpanFinishRecorder{}
	opts.NewSpanEventListener = LogTraceIntegrator
	opentracing.SetGlobalTracer(basictracer.NewWithOptions(opts))
}

func ServerOpenTracing(c *HttpContext) {
	textCarrier := opentracing.HTTPHeadersCarrier(c.RawRequest().Header)
	wireSpanContext, _ := opentracing.GlobalTracer().Extract(
		opentracing.TextMap, textCarrier)

	serverSpan := opentracing.GlobalTracer().StartSpan(
		c.RawRequest().Method+"__"+c.RawRequest().URL.Path,
		ext.RPCServerOption(wireSpanContext))

	defer serverSpan.Finish()
	ext.Component.Set(serverSpan, c.options.ServiceName)
	serverBasicSpanContext := serverSpan.Context().(basictracer.SpanContext)
	c.Set(TRACEID, serverBasicSpanContext.TraceID)
	c.Set(SPANID, serverBasicSpanContext.SpanID)
	c.Set(KEY_SPAN, serverSpan)
	RecordInputSpan(wireSpanContext, c)
	c.SetRawRequest(c.RawRequest().WithContext(gocontext.WithValue(gocontext.Background(), ContextKey, c)))
	stime := time.Now()

	c.Next()

	c.Set(PROCTIME, int64(time.Since(stime).Nanoseconds()/1000000))
	ext.HTTPStatusCode.Set(serverSpan, uint16(c.ResponseStatus()))
	RecordOutputSpan(serverSpan.Context(), c)

	return
}

func ClientOpenTracing(c *resty.Client, r *resty.Request) error {
	span := opentracing.SpanFromContext(r.Context())
	if span != nil {
		textCarrier := opentracing.HTTPHeadersCarrier(r.Header)
		_ = span.Tracer().Inject(span.Context(), opentracing.TextMap, textCarrier)
	}

	return nil
}

func SubCallResponseTracing(c *resty.Client, r *resty.Response) error {
	return RecordSubCallResponse(c, r)
}

func GetHttpContextFromRequest(r *http.Request) *HttpContext {
	if nil == r {
		return nil
	}
	ctx := r.Context()
	if nil == ctx {
		return nil
	}
	httpCtx := ctx.Value(ContextKey)
	if nil == httpCtx {
		return nil
	}

	return httpCtx.(*HttpContext)
}

func GetTraceIdFromRequest(r *http.Request) string {
	if nil == r {
		return ""
	}
	ctx := r.Context()
	if nil == ctx {
		return ""
	}
	httpCtx := ctx.Value(ContextKey)
	if nil == httpCtx {
		return ""
	}

	return httpCtx.(*HttpContext).GetString(TRACEID)
}

func GetSpanIdFromRequest(r *http.Request) string {
	if nil == r {
		return ""
	}
	ctx := r.Context()
	if nil == ctx {
		return ""
	}
	httpCtx := ctx.Value(ContextKey)
	if nil == httpCtx {
		return ""
	}

	return httpCtx.(*HttpContext).GetString(SPANID)
}
