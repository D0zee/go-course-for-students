package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	resCh := make(chan any)

	go func() {
		<-ctx.Done()
		close(resCh)
	}()

	out := in
	for _, stage := range stages {
		out = stage(out)
	}

	go func() {
		defer close(resCh)
		for val := range out {
			select {
			case <-ctx.Done():
				return
			case resCh <- val:
			}
		}
	}()

	return resCh
}
