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
	Offset    int
	Limit     int
	BlockSize int
	Conv      []string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "count of bytes from begin of file, which program must pass. by default - 0")
	flag.IntVar(&opts.Limit, "limit", math.MaxInt32, "max count of copy bytes. by default value is more than size of file - MaxInt32")
	flag.IntVar(&opts.BlockSize, "block-size", opts.Limit, "size of one block for copying. By default it equals opts.Limit")
	var names string
	flag.StringVar(&names, "name", "", "conversations under text")
	opts.Conv = strings.Split(names, ",")
	flag.Parse()
	//fmt.Println(opts)

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
	var reader io.Reader
	stream := os.Stdin
	bufferSize := opts.Limit + opts.Offset

	if len(opts.From) != 0 {
		var err error
		stream, err = os.Open(opts.From)
		if err != nil {
			return nil, err
		}
		fileInfo, err := stream.Stat()
		if err != nil {
			return nil, err
		}
		bufferSize = int(fileInfo.Size())
	}

	buf := make([]byte, bufferSize)

	reader = io.Reader(stream)
	var cntReadBytes int
	for {
		cntBytes, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		cntReadBytes += cntBytes
		if cntReadBytes > opts.Limit+opts.Offset {
			break
		}
	}
	if opts.Offset >= cntReadBytes {
		return nil, errors.New("offset must be less than input.size")
	}
	left := int(math.Min(float64(cntReadBytes), float64(opts.Limit))) + opts.Offset
	buf = buf[opts.Offset:left]
	return buf, nil
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

	err = WriteBytes(opts, bufFromReader)
	checkError(err)

	// todo: implement the functional requirements described in read.me
}
