package golden

import (
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/porschemacan/golden/libs"
	"strings"
)

var logTagMappingFunc = func(method string, path string) string {
	return "_undef"
}

func SetLogTagMapping(mappingFunc func(method string, path string) string) {
	logTagMappingFunc = mappingFunc
}

var LogTraceIntegrator = func() func(basictracer.SpanEvent) {
	var fields = map[string]interface{}{}
	var operationName string
	return func(e basictracer.SpanEvent) {
		switch t := e.(type) {
		case basictracer.EventCreate:
			operationName = t.OperationName
		case basictracer.EventFinish:
			tagName := operationName
			items := strings.Split(operationName, "__")
			if len(items) >= 2 {
				tagName = logTagMappingFunc(items[0], items[1])
			}
			entry := libs.GetDLog(tagName)
			for k, v := range fields {
				entry = entry.WithField(k, v)
			}
			entry.Info(operationName)
		case basictracer.EventTag:
			fields[t.Key] = t.Value
		case basictracer.EventLogFields:
			for _, f := range t.Fields {
				fields[f.Key()] = f.Value()
			}
		case basictracer.EventLog:
			fields[t.Event] = t.Payload
		}
	}
}

// record request
func RecordInputSpan(spanContext opentracing.SpanContext, c *HttpContext) {
	var traceId uint64
	var spanId uint64
	if spanContext != nil {
		basicSpanContext := spanContext.(basictracer.SpanContext)
		traceId = basicSpanContext.TraceID
		spanId = basicSpanContext.SpanID
	}

	libs.GetDLog("_com_request_in").
		WithField("traceid", traceId).
		WithField("parentSpanid", spanId).
		WithField("spanid", c.GetSpanId()).
		WithField("method", c.RawRequest().Method).
		WithField("uri", c.RawRequest().URL.Path).
		WithField("contentlength", c.RawRequest().ContentLength).
		Infoln()
}

// record response
func RecordOutputSpan(spanContext opentracing.SpanContext, c *HttpContext) {
	basicSpanContext := spanContext.(basictracer.SpanContext)

	libs.GetDLog("_com_request_out").
		WithField("traceid", basicSpanContext.TraceID).
		WithField("spanid", basicSpanContext.SpanID).
		WithField("method", c.RawRequest().Method).
		WithField("uri", c.RawRequest().URL.Path).
		WithField("proc_time", c.GetInt64(PROCTIME)).
		WithField("errno", c.GetInt64(ERROR_NO)).
		Infoln()
}

func RecordSubCallResponse(c *resty.Client, r *resty.Response) error {
	if r.Request.Context() == nil {
		return nil
	}
	span := opentracing.SpanFromContext(r.Request.Context())
	if span == nil {
		return nil
	}

	spanContext := span.Context()
	if spanContext == nil {
		return nil
	}

	basicSpanContext := spanContext.(basictracer.SpanContext)
	procTime := int64(r.Time().Nanoseconds() / 1000000)

	keyTag := "_com_http_success"
	errno := 0
	if r.IsError() {
		keyTag = "_com_http_failure"
		errno = r.StatusCode()
	}

	libs.GetDLog(keyTag).
		WithField("traceid", basicSpanContext.TraceID).
		WithField("spanid", basicSpanContext.SpanID).
		WithField("method", r.Request.Method).
		WithField("url", r.Request.URL).
		WithField("proc_time", procTime).
		WithField("errno", errno).
		Infoln()

	return nil
}

type SpanFinishRecorder struct {
}

func (*SpanFinishRecorder) RecordSpan(span basictracer.RawSpan) {
	libs.GetDLog("_undef").
		WithField("traceid", span.Context.TraceID).
		WithField("parentSpanid", span.ParentSpanID).
		WithField("spanid", span.Context.SpanID).
		WithField("starttime", span.Start.Unix()).
		WithField("duration", int64(span.Duration.Nanoseconds()/1000000)).
		Info(span.Context.Baggage)
}
