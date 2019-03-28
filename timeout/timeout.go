package grpc_timeout

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// WithDefault sets a default timeout on all requests contexts that do
// not already have a timeout or deadline.
func WithDefault(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Check if the request has an explicit timeout/deadline
		if _, exists := ctx.Deadline(); !exists {
			// The request does not have an explicit timeout/deadline:
			// use the default timeout instead.
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
