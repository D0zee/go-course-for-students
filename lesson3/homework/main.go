package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From          string
	To            string
	Offset        int64
	Limit         int64
	BlockSize     int64
	Conv          []string
	LimitProvided bool
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "count of bytes from begin of file, which program must pass. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", 0, "max count of copy bytes. by default value is more than size of file - MaxInt32")
	flag.Int64Var(&opts.BlockSize, "block-size", 1024, "size of one block for copying. By default it equals opts.Limit")
	convParameters := ""
	flag.StringVar(&convParameters, "conv", "", "conversations under text")

	flag.Parse()
	if convParameters != "" {
		opts.Conv = strings.Split(convParameters, ",")
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "limit" {
			opts.LimitProvided = true
		}
	})
	return &opts, nil
}

func ValidateConv(opt *Options) error {
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
	if len(opt.Conv) > 2 {
		return errors.New("wrong count of conv opt")
	}
	return ValidateConv(opt)
}

func ReadBytes(opts *Options) (str string, err error) {
	stream := os.Stdin
	if opts.From != "" {
		stream, err = os.Open(opts.From)
		if err != nil {
			return "", err
		}
		defer func() {
			cerr := stream.Close()
			if err == nil {
				err = cerr
			}
		}()

	}
	if !opts.LimitProvided { // if user does not provide these parameters, it will be equal size of input
		fileInfo, _ := stream.Stat()
		opts.Limit = fileInfo.Size() + 1
	}
	reader := io.Reader(stream)
	builder := strings.Builder{}
	cntReadBytes := int64(0)

	cntWrittenBytes, _ := io.CopyN(io.Discard, reader, opts.Offset)
	cntReadBytes += cntWrittenBytes

	cntWrittenBytes, _ = io.CopyN(&builder, reader, opts.Limit)
	cntReadBytes += cntWrittenBytes

	if opts.Offset >= cntReadBytes {
		return "", errors.New("offset must be less than input.size")
	}
	str = builder.String()
	return str, err
}

type ConverterWriter struct {
	writer           io.Writer
	converterOptions []string
}

func (converter *ConverterWriter) Write(p []byte) (int, error) {
	if converter.converterOptions != nil {
		for _, c := range converter.converterOptions {
			if c == "upper_case" {
				p = bytes.ToUpper(p)
			}
			if c == "lower_case" {
				p = bytes.ToLower(p)
			}
			if c == "trim_spaces" {
				p = bytes.TrimSpace(p)
			}
		}
	}
	return converter.writer.Write(p)
}

func WriteBytes(opts *Options, str string) (err error) {
	stream := os.Stdout
	if opts.To != "" {
		stream, err = os.OpenFile(opts.To, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		checkError(err)
	}
	writer := io.WriteCloser(stream)
	cw := ConverterWriter{writer: writer, converterOptions: opts.Conv}
	_, err = cw.Write([]byte(str))
	if err != nil {
		return err
	}
	return
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

	stringFromReader, err := ReadBytes(opts)
	checkError(err)

	err = WriteBytes(opts, stringFromReader)
	checkError(err)
}
