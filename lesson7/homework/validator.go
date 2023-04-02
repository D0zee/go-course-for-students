package homework

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	for _, err := range v {
		s += err.Err.Error()
	}
	return s
}

func validateLen(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	tagValue := ft.Tag.Get("validate")
	after := strings.TrimPrefix(tagValue, "len:")
	l, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if l != len(vt.String()) {
		errs = append(errs, ValidationError{errors.New("wrong length")})
	}
	return errs
}

func Contains[T comparable](t []T, needle T) bool {
	for _, v := range t {
		if v == needle {
			return true
		}
	}
	return false
}

func ConvertToInts(s []string) ([]int, error) {
	ints := make([]int, len(s))
	for idx, pattern := range s {
		i, err := strconv.Atoi(pattern)
		if err != nil {
			return nil, ErrInvalidValidatorSyntax
		}
		ints[idx] = i
	}
	return ints, nil
}

func validateIn(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	tagValue := ft.Tag.Get("validate")
	after := strings.TrimPrefix(tagValue, "in:")
	if after == "" {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	patterns := strings.Split(after, ",")
	contain := false
	if vt.Kind() == reflect.Int {
		ints, err := ConvertToInts(patterns)
		if err != nil {
			return append(errs, ValidationError{err})
		}
		contain = Contains(ints, int(vt.Int()))
	} else if vt.Kind() == reflect.String {
		contain = Contains(patterns, vt.String())
	}
	if !contain {
		errs = append(errs, ValidationError{errors.New("don't contain")})
	}
	return errs
}

func validateMin(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	tagValue := ft.Tag.Get("validate")
	after := strings.TrimPrefix(tagValue, "min:")
	if after == "" {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	min, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if vt.Kind() == reflect.Int {
		if int(vt.Int()) < min {
			return append(errs, ValidationError{errors.New("min isn't correct")})
		}
	} else if vt.Kind() == reflect.String {
		if len(vt.String()) < min {
			return append(errs, ValidationError{errors.New("min isn't correct")})
		}
	}
	return errs
}

func validateMax(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	tagValue := ft.Tag.Get("validate")
	after := strings.TrimPrefix(tagValue, "max:")
	if after == "" {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	max, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if vt.Kind() == reflect.Int {
		if int(vt.Int()) > max {
			return append(errs, ValidationError{errors.New("min isn't correct")})
		}
	} else if vt.Kind() == reflect.String {
		if len(vt.String()) > max {
			return append(errs, ValidationError{errors.New("min isn't correct")})
		}
	}
	return errs
}

func Validate(v any) error {
	typeV := reflect.TypeOf(v)
	valueV := reflect.ValueOf(v)
	if typeV.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	n := valueV.NumField()
	errs := ValidationErrors{}
	for i := 0; i < n; i++ {
		ft := typeV.Field(i)
		vt := valueV.Field(i)

		tv := ft.Tag.Get("validate")
		if tv == "" {
			continue
		}
		if !ft.IsExported() {
			errs = append(errs, ValidationError{ErrValidateForUnexportedFields})
			continue
		}
		if strings.HasPrefix(tv, "len:") {
			errs = validateLen(errs, ft, vt)
		}
		if strings.HasPrefix(tv, "in:") {
			errs = validateIn(errs, ft, vt)
		}
		if strings.HasPrefix(tv, "min:") {
			errs = validateMin(errs, ft, vt)
		}
		if strings.HasPrefix(tv, "max:") {
			errs = validateMax(errs, ft, vt)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
