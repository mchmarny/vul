package metric

import (
	"context"
	"time"
)

type metrixContext struct {
	ctx context.Context
}

func (c metrixContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c metrixContext) Done() <-chan struct{}             { return nil }
func (c metrixContext) Err() error                        { return nil }
func (c metrixContext) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// ToMetricContext returns a context that is never canceled.
func ToMetricContext(ctx context.Context) context.Context {
	return metrixContext{ctx: ctx}
}
