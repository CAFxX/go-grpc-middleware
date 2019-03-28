package grpc_timeout_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	grpc_timeout "github.com/grpc-ecosystem/go-grpc-middleware/timeout"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithDefault(t *testing.T) {
	const noTimeout = time.Duration(-1)

	cases := []struct {
		defaultTimeout  time.Duration
		explicitTimeout time.Duration
		wantTimeout     time.Duration
	}{
		{defaultTimeout: 2 * time.Second, explicitTimeout: noTimeout, wantTimeout: 2 * time.Second},
		{defaultTimeout: 2 * time.Second, explicitTimeout: 0 * time.Second, wantTimeout: 0 * time.Second},
		{defaultTimeout: 2 * time.Second, explicitTimeout: 1 * time.Second, wantTimeout: 1 * time.Second},
		{defaultTimeout: 2 * time.Second, explicitTimeout: 2 * time.Second, wantTimeout: 2 * time.Second},
		{defaultTimeout: 2 * time.Second, explicitTimeout: 3 * time.Second, wantTimeout: 3 * time.Second},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			wd := grpc_timeout.WithDefault(c.defaultTimeout)

			ctx := context.Background()
			if c.explicitTimeout != noTimeout {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, c.explicitTimeout)
				defer cancel()
			}

			method := "method"
			req := &struct{ _ float64 }{666}
			rep := &struct{ _ int }{42}
			cc := &grpc.ClientConn{}
			opts := []grpc.CallOption{grpc.WaitForReady(true)}
			_err := errors.New("some error")

			invoker := func(_ctx context.Context, _method string, _req, _rep interface{}, _cc *grpc.ClientConn, _opts ...grpc.CallOption) error {
				deadline, exists := _ctx.Deadline()
				assert.True(t, exists)
				assert.InDelta(t, c.wantTimeout.Seconds(), time.Until(deadline).Seconds(), 0.01)
				assert.Equal(t, method, _method)
				assert.Equal(t, req, _req.(*struct{ _ float64 }))
				assert.Equal(t, rep, _rep.(*struct{ _ int }))
				assert.Equal(t, cc, _cc)
				assert.Equal(t, opts, _opts)
				return _err
			}

			err := wd(ctx, method, req, rep, cc, invoker, opts...)
			assert.Equal(t, _err, err)
		})
	}
}
