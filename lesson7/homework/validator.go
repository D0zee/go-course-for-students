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
var ErrValidatorLen = errors.New("validator LEN: wrong length of string")
var ErrValidatorIn = errors.New("validator IN: there isn't the value")

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

type Field struct {
	Type  reflect.StructField
	Value reflect.Value
}

func getValueOfTag(ft reflect.StructField, prefix string) (string, error) {
	tagValue := ft.Tag.Get("validate")
	after := strings.TrimPrefix(tagValue, prefix)
	if after == "" {
		return "", ErrInvalidValidatorSyntax
	}
	return after, nil
}

func validateLen(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	after, err := getValueOfTag(ft, "len:")
	if err != nil {
		return append(errs, ValidationError{err})
	}

	l, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if l != len(vt.String()) {
		errs = append(errs, ValidationError{ErrValidatorLen})
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
	arr := make([]int, len(s))
	for idx, pattern := range s {
		i, err := strconv.Atoi(pattern)
		if err != nil {
			return nil, ErrInvalidValidatorSyntax
		}
		arr[idx] = i
	}
	return arr, nil
}

func validateIn(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	after, err := getValueOfTag(ft, "in:")
	if err != nil {
		return append(errs, ValidationError{err})
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
		errs = append(errs, ValidationError{ErrValidatorIn})
	}
	return errs
}

type PredicateWithInfo struct {
	name      string
	predicate func(int, int) bool
}

func (p PredicateWithInfo) getValidationError() ValidationError {
	return ValidationError{errors.New("field isn't validated by " + p.name + " function")}
}

func validateMinMax(errs ValidationErrors, ft reflect.StructField, vt reflect.Value, p PredicateWithInfo) ValidationErrors {
	after, err := getValueOfTag(ft, p.name+":")
	if err != nil {
		return append(errs, ValidationError{err})
	}

	bound, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if vt.Kind() == reflect.Int {
		if !p.predicate(int(vt.Int()), bound) {
			return append(errs, p.getValidationError())
		}
	} else if vt.Kind() == reflect.String {
		if !p.predicate(len(vt.String()), bound) {
			return append(errs, p.getValidationError())
		}
	}
	return errs
}

func ValidateField(errs ValidationErrors, ft reflect.StructField, vt reflect.Value) ValidationErrors {
	tv := ft.Tag.Get("validate")
	if tv == "" {
		return errs
	}
	if strings.HasPrefix(tv, "len:") {
		errs = validateLen(errs, ft, vt)
	}
	if strings.HasPrefix(tv, "in:") {
		errs = validateIn(errs, ft, vt)
	}
	if strings.HasPrefix(tv, "min:") {
		errs = validateMinMax(errs, ft, vt,
			PredicateWithInfo{name: "min", predicate: func(a int, b int) bool {
				return a >= b
			}})
	}
	if strings.HasPrefix(tv, "max:") {
		errs = validateMinMax(errs, ft, vt,
			PredicateWithInfo{name: "max", predicate: func(a int, b int) bool {
				return a <= b
			}})
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
		errs = ValidateField(errs, ft, vt)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
