package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From              string
	To                string
	Offset            int64
	Limit             int64
	BlockSize         int64
	Conv              []string
	LimitProvided     bool
	BlockSizeProvided bool
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "count of bytes from begin of file, which program must pass. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", 0, "max count of copy bytes. by default value is more than size of file - MaxInt32")
	flag.Int64Var(&opts.BlockSize, "block-size", 0, "size of one block for copying. By default it equals opts.Limit")
	convParameters := ""
	flag.StringVar(&convParameters, "name", "", "conversations under text")

	flag.Parse()
	if convParameters != "" {

		opts.Conv = strings.Split(convParameters, ",")
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "limit" {
			opts.LimitProvided = true
		}
		if f.Name == "block-size" {
			opts.BlockSizeProvided = true
		}
	})
	return &opts, nil
}

func ValidateConv(opt Options) error {
	if len(opt.Conv) == 0 {
		return nil
	}
	allowedValues := map[string]interface{}{
		"upper_case":  nil,
		"lower_case":  nil,
		"trim_spaces": nil,
	}
	containsUpperCase := false
	containsLowerCase := false
	for _, value := range opt.Conv {
		if _, exists := allowedValues[value]; !exists {
			return fmt.Errorf("invalid value '%s' in -conv flag", value)
		}
		if value == "upper_case" {
			containsUpperCase = true
		}
		if value == "lower_case" {
			containsLowerCase = true
		}
	}
	if containsUpperCase && containsLowerCase {
		return errors.New("-conv parameters cannot contain both upper_case and lower_case")
	}
	return nil
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
	if len(opt.Conv) >= 2 {
		return errors.New("wrong count of conv opt")
	}
	return ValidateConv(*opt)
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
	if !opts.LimitProvided { // if user does not provide these parameters, it will be equal size of input
		fileInfo, _ := stream.Stat()
		opts.Limit = fileInfo.Size() + 1
	}
	if !opts.BlockSizeProvided {
		opts.BlockSize = opts.Limit
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
	stream := os.Stdout
	if opts.To != "" {
		var err error
		stream, err = os.OpenFile(opts.To, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		checkError(err)
	}
	writer := io.WriteCloser(stream)
	_, err := writer.Write(buf)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
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

}
