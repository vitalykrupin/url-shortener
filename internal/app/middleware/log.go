// Package middleware provides HTTP middleware functions
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData stores response data for logging
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter wraps http.ResponseWriter to capture response data
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write writes the data to the connection as part of an HTTP reply
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader sends an HTTP response header with the provided status code
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logging provides request logging middleware
func Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestLogger := zap.Must(zap.NewProduction()).Sugar()
		defer func() {
			_ = requestLogger.Sync()
		}()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		handler.ServeHTTP(&lw, req)
		requestLogger.Info(
			zap.String("Method", req.Method),
			zap.String("URI", req.RequestURI),
			zap.String("ResponseDuration", time.Since(start).String()),
			zap.Int("ResponseStatus", responseData.status),
			zap.Int("ResponseBodySize", responseData.size),
		)

	})
}
