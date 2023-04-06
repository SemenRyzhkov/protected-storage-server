package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// LogInterceptor log interceptor type
type LogInterceptor struct {
}

// NewLogInterceptor constructor
func NewLogInterceptor() *LogInterceptor {
	return &LogInterceptor{}
}

// Unary is a gRPC logging interceptor
func (l *LogInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		status := "error"
		if err == nil {
			status = "ok"
		}
		log.Printf("[gRPC] %s requestTime: %s responseTime: %s status: %s",
			info.FullMethod, start.Format(time.RFC3339), time.Since(start), status)
		return resp, err
	}
}
