package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    int64
	Limit     int64
	BlockSize int64
	Conv      []string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "count of bytes from begin of file, which program must pass. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", math.MaxInt32, "max count of copy bytes. by default value is more than size of file - MaxInt32")
	flag.Int64Var(&opts.BlockSize, "block-size", opts.Limit, "size of one block for copying. By default it equals opts.Limit")
	var names string
	flag.StringVar(&names, "name", "", "conversations under text")
	opts.Conv = strings.Split(names, ",")
	flag.Parse()

	return &opts, nil
}

func ValidateOptions(opt *Options) error {
	if opt.Offset < 0 {
		return errors.New("-offset is not correct")
	}
	if opt.Limit < 0 {
		return errors.New("-limit is not correct")
	}
	if opt.BlockSize < 0 {
		return errors.New("-block-size is not correct")
	}
	return nil
}

func ReadBytes(opts *Options) ([]byte, error) {
	stream := os.Stdin
	if opts.From != "" {
		var err error
		stream, err = os.Open(opts.From)
		if err != nil {
			return nil, err
		}
		defer stream.Close()
	}
	reader := io.Reader(stream)
	builder := strings.Builder{}
	cntReadBytes, _ := io.CopyN(io.Discard, reader, opts.Offset)
	written, _ := io.CopyN(&builder, reader, opts.Limit)
	cntReadBytes += written

	if opts.Offset >= cntReadBytes {
		return nil, errors.New("offset must be less than input.size")
	}
	return []byte(builder.String()), nil
}

func WriteBytes(opts *Options, buf []byte) error {
	var writer io.WriteCloser
	if len(opts.To) == 0 {
		writer = io.WriteCloser(os.Stdout)
	} else {
		stream, err := os.OpenFile(opts.To, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		checkError(err)
		writer = io.WriteCloser(stream)
	}

	_, err := writer.Write(buf)
	if err != nil {
		return err
	}
	//println("before reading")
	err = writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
}

func main() {
	opts, err := ParseFlags()
	checkError(err)

	err = ValidateOptions(opts)
	checkError(err)

	bufFromReader, err := ReadBytes(opts)
	checkError(err)
	//println("read: bytes")

	err = WriteBytes(opts, bufFromReader)
	checkError(err)
	//println("write: bytes")

	// todo: implement the functional requirements described in read.me
}
