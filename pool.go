package gtask

import (
	"context"
)

type Pool interface {
	Go(ctx context.Context, f func())
	Do(ctx context.Context, tasks ...*Task) <-chan *RunResult
}
