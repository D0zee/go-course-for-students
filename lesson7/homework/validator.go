package homework

import (
	"fmt"
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
	//Type  reflect.StructField
	Value    reflect.Value
	TagValue string
}

func getValueOfTag(f Field, prefix string) (string, error) {
	tagValue := f.TagValue
	after := strings.TrimPrefix(tagValue, prefix)
	if after == "" {
		return "", ErrInvalidValidatorSyntax
	}
	return after, nil
}

func validateLen(errs ValidationErrors, f Field) ValidationErrors {
	fmt.Println(f.Value.String())
	after, err := getValueOfTag(f, "len:")
	if err != nil {
		return append(errs, ValidationError{err})
	}

	l, err := strconv.Atoi(after)
	if err != nil || l < 0 {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if l != len(f.Value.String()) {
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

func validateIn(errs ValidationErrors, f Field) ValidationErrors {
	after, err := getValueOfTag(f, "in:")
	if err != nil {
		return append(errs, ValidationError{err})
	}
	patterns := strings.Split(after, ",")
	contain := false
	if f.Value.Kind() == reflect.Int {
		ints, err := ConvertToInts(patterns)
		if err != nil {
			return append(errs, ValidationError{err})
		}
		contain = Contains(ints, int(f.Value.Int()))
	} else if f.Value.Kind() == reflect.String {
		contain = Contains(patterns, f.Value.String())
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

func validateMinMax(errs ValidationErrors, f Field, p PredicateWithInfo) ValidationErrors {
	after, err := getValueOfTag(f, p.name+":")
	if err != nil {
		return append(errs, ValidationError{err})
	}

	bound, err := strconv.Atoi(after)
	if err != nil {
		return append(errs, ValidationError{ErrInvalidValidatorSyntax})
	}
	if f.Value.Kind() == reflect.Int {
		if !p.predicate(int(f.Value.Int()), bound) {
			return append(errs, p.getValidationError())
		}
	} else if f.Value.Kind() == reflect.String {
		if !p.predicate(len(f.Value.String()), bound) {
			return append(errs, p.getValidationError())
		}
	}
	return errs
}

func ValidateField(errs ValidationErrors, f Field) ValidationErrors {
	tv := f.TagValue
	fmt.Println("validate:" + tv)
	if tv == "" {
		return errs
	}
	if strings.HasPrefix(tv, "len:") {
		errs = validateLen(errs, f)
	}
	if strings.HasPrefix(tv, "in:") {
		errs = validateIn(errs, f)
	}
	if strings.HasPrefix(tv, "min:") {
		errs = validateMinMax(errs, f,
			PredicateWithInfo{name: "min", predicate: func(a int, b int) bool {
				return a >= b
			}})
	}
	if strings.HasPrefix(tv, "max:") {
		errs = validateMinMax(errs, f,
			PredicateWithInfo{name: "max", predicate: func(a int, b int) bool {
				return a <= b
			}})
	}
	return errs
}

func Validate(val any) error {
	typeV := reflect.TypeOf(val)
	valueV := reflect.ValueOf(val)
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
		if vt.Kind() == reflect.Slice {
			for i := 0; i < vt.Len(); i++ {
				errs = ValidateField(errs, Field{TagValue: tv, Value: vt.Index(i)})

				if len(errs) != 0 && errs[len(errs)-1].Err == ErrInvalidValidatorSyntax { // if we get InvalidSyntax on first element, then no sense to continue iterate
					break
				}
			}
			continue
		}

		errs = ValidateField(errs, Field{Value: vt, TagValue: tv})
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
