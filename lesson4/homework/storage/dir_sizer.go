package storage

import (
	"context"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	// maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

type workerOutput struct {
	err error
	res Result
}


func collectFromChannel(channel <-chan workerOutput, resInDir *Result, length int) error {
	for i := 0; i < length; i++ {
		resFromChannel := <-channel
		if resFromChannel.err != nil {
			return resFromChannel.err
		}
		resInDir.Size += resFromChannel.res.Size
		resInDir.Count += resFromChannel.res.Count
	}
	return nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return Result{}, err
	}

	outputChannel := make(chan workerOutput, len(files)+len(dirs))

	go func(channel chan<- workerOutput) {
		res := Result{}
		for _, file := range files {
			size, err := file.Stat(ctx)
			if err != nil {
				channel <- workerOutput{err: err, res: res}
			} else {
				channel <- workerOutput{err: err, res: Result{size, 1}}
			}
		}
	}(outputChannel)

	for _, dir := range dirs {
		go func(dir Dir, channel chan<- workerOutput) {
			resLocal, err := a.Size(ctx, dir)
			if err != nil {
				channel <- workerOutput{res: resLocal, err: err}
			} else {
				channel <- workerOutput{res: resLocal, err: nil}
			}
		}(dir, outputChannel)
	}

	resInDir := Result{}
	err = collectFromChannel(outputChannel, &resInDir, len(files)+len(dirs))
	if err != nil {
		return Result{}, err
	}

	return resInDir, nil
}
