// Package log provides context-aware and structured logging capabilities.
package log

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strconv"
)

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})
	// Warn uses fmt.Sprint to construct and log a message at WARN level
	Warn(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
	// Warnf uses fmt.Sprintf to construct and log a message at WARN level
	Warnf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
	cacheStatusIDKey
)

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

// With returns a logger based off the root logger and decorates it with the given context and arguments.
//
// If the context contains request ID and/or correlation ID information (recorded via WithRequestID()
// and WithCorrelationID()), they will be added to every log message generated by the new logger.
//
// The arguments should be specified as a sequence of name, value pairs with names being strings.
// The arguments will also be added to every log message generated by the logger.
func (l *logger) With(ctx context.Context, args ...interface{}) Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
		if status, ok := ctx.Value(cacheStatusIDKey).(string); ok {
			args = append(args, zap.String("cache_status", status))
		}
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}

func (l *logger) GetZap() *zap.SugaredLogger {
	return l.SugaredLogger
}

// WithRequest returns a context which knows the request ID and correlation ID in the given request.
func WithRequest(ctx context.Context, req *fasthttp.Request) context.Context {
	id := getRequestID(req)
	if id == "" {
		id = uuid.New().String()
	}

	ctx = context.WithValue(ctx, requestIDKey, id)
	if id = getCorrelationID(req); id != "" {
		ctx = context.WithValue(ctx, correlationIDKey, id)
	}

	return ctx
}

// getCorrelationID extracts the correlation ID from the HTTP request
func getCorrelationID(req *fasthttp.Request) string {
	return string(req.Header.Peek("X-Correlation-ID"))
}

// getRequestID extracts the correlation ID from the HTTP request
func getRequestID(req *fasthttp.Request) string {
	return string(req.Header.Peek("X-Request-ID"))
}

func LoggerMiddleware(l Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var body string
		contentType := utils.UnsafeString(c.Request().Header.ContentType())
		if contentType == "multipart/form-data; boundary=boundary" {
			body = " ------ FILE " + strconv.Itoa(c.Request().Header.ContentLength()) + " size ------"
		} else {
			body = utils.UnsafeString(c.Body())
		}

		l.With(
			WithRequest(c.Context(), c.Request()),
		).
			Info(c.String(), body)

		return c.Next()
	}
}
