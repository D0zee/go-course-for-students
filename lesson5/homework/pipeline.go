package executor

import (
	"context"
	"sync"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	resCh := make(chan any)
	done := sync.Once{}
	closeCh := func() { close(resCh) }

	go func() {
		<-ctx.Done()
		done.Do(closeCh)
	}()

	out := in
	for _, stage := range stages {
		out = stage(out)
	}

	go func() {
		defer done.Do(closeCh)
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
