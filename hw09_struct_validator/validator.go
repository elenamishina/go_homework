package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidateLength = errors.New("length validate failed")
	ErrValidateRegexp = errors.New("regexp validate failed")
	ErrValidateIn     = errors.New("entry into the list validate failed")
	ErrValidateMin    = errors.New("min validate failed")
	ErrValidateMax    = errors.New("max validate failed")

	ErrUnsupportedType = errors.New("unsupported validate type")
	ErrUnsupportedRule = errors.New("unsuported validate rule")
	ErrNotCorrectRule  = errors.New("not correct validate rule ([tag]:[value])")
	ErrInputType       = errors.New("expect Struct as input")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	if v.Field != "" {
		return fmt.Sprintf("%s: %v", v.Field, v.Err)
	}
	return v.Err.Error()
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return sb.String()
}

func (v ValidationErrors) Unwrap() []error {
	errs := make([]error, len(v))
	for i, err := range v {
		errs[i] = err
	}
	return errs
}

func Validate(v interface{}) error {
	var errs ValidationErrors
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrInputType
	}

	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := valType.Field(i)
		fieldValue := val.Field(i)

		tagValue := field.Tag.Get("validate")
		if tagValue == "" {
			continue
		}

		rules := strings.Split(tagValue, "|")
		err := CheckRule(rules, field.Name, fieldValue, &errs)
		if err != nil {
			return err
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func CheckRule(rules []string, fieldName string, fieldValue reflect.Value, errs *ValidationErrors) error {
	for _, rule := range rules {
		partsRule := strings.Split(rule, ":")
		if len(partsRule) != 2 {
			return ErrNotCorrectRule
		}
		ruleName := partsRule[0]
		ruleValue := partsRule[1]

		switch ruleName {
		case "len":
			err := ValidateLength(fieldName, fieldValue, ruleValue, errs)
			if err != nil {
				return err
			}
		case "regexp":
			err := ValidateRegexp(fieldName, fieldValue, ruleValue, errs)
			if err != nil {
				return err
			}
		case "in":
			err := ValidateIn(fieldName, fieldValue, ruleValue, errs)
			if err != nil {
				return err
			}
		case "min":
			err := ValidateMin(fieldName, fieldValue, ruleValue, errs)
			if err != nil {
				return err
			}
		case "max":
			err := ValidateMax(fieldName, fieldValue, ruleValue, errs)
			if err != nil {
				return err
			}
		default:
			return ErrUnsupportedRule
		}
	}
	return nil
}

func ValidateLength(fieldName string, fieldVal reflect.Value, ruleVal string, errs *ValidationErrors) error {
	length, err := strconv.Atoi(ruleVal)
	if err != nil {
		return err
	}

	//exhaustive:ignore
	switch fieldVal.Kind() {
	case reflect.String:
		if len(fieldVal.String()) != length {
			*errs = append(*errs, ValidationError{
				Field: fieldName,
				Err:   ErrValidateLength,
			})
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < fieldVal.Len(); i++ {
			ValidateLength(fieldName, fieldVal.Index(i), ruleVal, errs)
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}

func ValidateRegexp(fieldName string, fieldVal reflect.Value, ruleVal string, errs *ValidationErrors) error {
	regex, err := regexp.Compile(ruleVal)
	if err != nil {
		return err
	}
	//exhaustive:ignore
	switch fieldVal.Kind() {
	case reflect.String:
		if !regex.MatchString(fieldVal.String()) {
			*errs = append(*errs, ValidationError{
				Field: fieldName,
				Err:   ErrValidateRegexp,
			})
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < fieldVal.Len(); i++ {
			ValidateRegexp(fieldName, fieldVal.Index(i), ruleVal, errs)
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}

func ValidateIn(fieldName string, fieldVal reflect.Value, ruleVal string, errs *ValidationErrors) error {
	listIn := strings.Split(ruleVal, ",")
	if len(listIn) == 0 {
		return ErrNotCorrectRule
	}
	//exhaustive:ignore
	switch fieldVal.Kind() {
	case reflect.String:
		err := ValidateInGeneric(fieldVal.String(), listIn)
		if err != nil {
			*errs = append(*errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		val, err := extractNumber(fieldVal)
		if err != nil {
			return err
		}
		nums, err := convertStrToNums(listIn)
		if err != nil {
			return err
		}
		err = ValidateInGeneric(val, nums)
		if err != nil {
			*errs = append(*errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < fieldVal.Len(); i++ {
			ValidateIn(fieldName, fieldVal.Index(i), ruleVal, errs)
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}

func ValidateInGeneric[T comparable](value T, allowedValues []T) error {
	found := false
	for _, v := range allowedValues {
		if value == v {
			found = true
		}
	}
	if !found {
		return ErrValidateIn
	}
	return nil
}

func convertStrToNums(listVal []string) ([]float64, error) {
	nums := make([]float64, len(listVal))
	for i := range listVal {
		num, err := strconv.ParseFloat(listVal[i], 64)
		if err != nil {
			return nums, err
		}
		nums[i] = num
	}
	return nums, nil
}

func ValidateMin(fieldName string, fieldVal reflect.Value, ruleVal string, errs *ValidationErrors) error {
	minValue, err := strconv.ParseFloat(ruleVal, 64)
	if err != nil {
		return err
	}

	val, err := extractNumber(fieldVal)
	if err != nil {
		return err
	}

	if val < minValue {
		*errs = append(*errs, ValidationError{
			Field: fieldName,
			Err:   ErrValidateMin,
		})
	}
	return nil
}

func ValidateMax(fieldName string, fieldVal reflect.Value, ruleVal string, errs *ValidationErrors) error {
	maxValue, err := strconv.ParseFloat(ruleVal, 64)
	if err != nil {
		return err
	}

	val, err := extractNumber(fieldVal)
	if err != nil {
		return err
	}
	if val > maxValue {
		*errs = append(*errs, ValidationError{
			Field: fieldName,
			Err:   ErrValidateMax,
		})
	}
	return nil
}

func extractNumber(fieldVal reflect.Value) (float64, error) {
	//exhaustive:ignore
	switch fieldVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(fieldVal.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(fieldVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fieldVal.Float(), nil
	default:
		return 0, ErrUnsupportedType
	}
}
