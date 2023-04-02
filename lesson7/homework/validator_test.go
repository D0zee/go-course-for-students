package homework

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "LEN: all elements of slices is valid",
			args: args{
				v: struct {
					MaxA []string `validate:"len:1"`
					MaxB []string `validate:"len:2"`
				}{
					MaxA: []string{"e"},
					MaxB: []string{"aa"},
				},
			},
			wantErr: false,
		},
		{
			name: "LEN: wrong len and wrong annotation",
			args: args{
				v: struct {
					MaxA []string `validate:"len:1"`
					MaxB []string `validate:"len:2"`

					wrongAn1 []string `validate:"len:-2"` // error
					wrongAn2 []string `validate:"len:"`   // error
					wrongAn3 []string `validate:"len:fd"` // error
					empty    []string `validate:"len:1"`  // error

				}{
					MaxA:     []string{"e", "12", "asdf"}, // 2 errors
					MaxB:     []string{"aa", "1"},         // error
					wrongAn1: []string{},
					wrongAn2: []string{},
					wrongAn3: []string{},
					empty:    []string{},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 7)
				return true
			},
		},

		{
			name: "MIN/MAX: all is right",
			args: args{
				v: struct {
					MaxA    []string `validate:"max:3"`
					MaxB    []string `validate:"min:2"`
					MaxInt  []int    `validate:"min:-5"`
					MaxInt2 []int    `validate:"max:5"`
					//MaxEmpty []int    `validate:"max:0"`
				}{
					MaxA:    []string{"e", "12", "asd"},
					MaxB:    []string{"aa", "122"},
					MaxInt:  []int{-5, 1, 0, 2, 65, 100},
					MaxInt2: []int{1, -1, 0, 5, 4},
					//MaxEmpty: []int{},
				},
			},
			wantErr: false,
		},

		{
			name: "MIN/MAX: errors",
			args: args{
				v: struct {
					MaxA     []string `validate:"max:abadsf"` // error
					MaxB     []string `validate:"min:"`       // error
					MaxInt   []int    `validate:"min:-5"`
					wrongMax []int    `validate:"max:5"`
				}{
					MaxA:     []string{"e", "12", "asd"}, //
					MaxB:     []string{"aa", "122"},
					MaxInt:   []int{},
					wrongMax: []int{100}, // error
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 3)
				return true
			},
		},

		{
			name: "SLICE: all is right IN",
			args: args{
				v: struct {
					InA     []string `validate:"in:ab,cd"`
					InB     []string `validate:"in:aa,bb,cd,ee"`
					InC     []int    `validate:"in:1,2,3"`
					InD     []int    `validate:"in:5"`
					InEmpty []int    `validate:"in:5"`
				}{
					InA:     []string{"ab", "ab"},
					InB:     []string{"aa", "bb", "cd", "ee"},
					InC:     []int{1, 2, 3},
					InD:     []int{5, 5, 5},
					InEmpty: []int{},
				},
			},
			wantErr: false,
		},

		{
			name: "SLICE: wrong in",
			args: args{
				v: struct {
					InA []string `validate:"in:ab,cd"`
					InB []string `validate:"in:,,,"`
					InC []int    `validate:"in:1,2,3"`
					InD []int    `validate:"in:fdsa"` // error
				}{
					InA: []string{"ab", "ba"}, // error because of "ba"
					InB: []string{"", ""},
					InC: []int{1, 2, 4, 6}, // 2 errors because of 4 and 6
					InD: []int{5, 5, 5},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 4)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
