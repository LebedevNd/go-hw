package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	message := ""
	for _, validationError := range v {
		message += validationError.Err.Error()
	}
	return message
}

func Validate(v interface{}) error {
	errorsPack := ValidationErrors{}

	refValueStruct := reflect.ValueOf(v)
	refTypeStruct := refValueStruct.Type()
	numFields := refTypeStruct.NumField()

	for i := 0; i < numFields; i++ {
		refFieldType := refTypeStruct.Field(i)
		refFieldValue := refValueStruct.Field(i)
		fieldTag, found := refFieldType.Tag.Lookup("validate")
		fieldName := refFieldType.Name

		if found {
			validators := strings.Split(fieldTag, "|")
			for _, validator := range validators {
				errs := validateField(fieldName, validator, refFieldValue)
				errorsPack = append(errorsPack, errs...)
			}
		}
	}

	return errorsPack
}

func validateField(fieldName string, validator string, refValue reflect.Value) ValidationErrors {
	var errs ValidationErrors

	fieldTypeBase := refValue.Type().Kind().String()
	fieldType := refValue.Type().String()

	switch {
	case fieldType == "string" || fieldTypeBase == "string":
		err, validationFailed := validateStringField(fieldName, validator, refValue.String())
		if validationFailed {
			errs = append(errs, err)
		}
	case fieldType == "int" || fieldTypeBase == "int":
		err, validationFailed := validateIntField(fieldName, validator, refValue.Int())
		if validationFailed {
			errs = append(errs, err)
		}
	case fieldType == "[]string" || fieldTypeBase == "[]string":
		stringSlice := refValue.Interface().([]string)
		for _, value := range stringSlice {
			err, validationFailed := validateStringField(fieldName, validator, value)
			if validationFailed {
				errs = append(errs, err)
			}
		}
	case fieldType == "[]int" || fieldTypeBase == "[]int":
		intSlice := refValue.Interface().([]int64)
		for _, value := range intSlice {
			err, validationFailed := validateIntField(fieldName, validator, value)
			if validationFailed {
				errs = append(errs, err)
			}
		}
	default:
	}

	return errs
}

func validateIntField(name string, validator string, value int64) (ValidationError, bool) {
	validationArray := strings.Split(validator, ":")
	if len(validationArray) != 2 {
		return ValidationError{}, false
	}

	validatorName := validationArray[0]
	validatorValue := validationArray[1]

	if validatorName == "min" {
		validLen, err := strconv.Atoi(validatorValue)
		if err == nil && value < int64(validLen) {
			return ValidationError{
				name,
				errors.New("wrong int length, must be more or equal to " + (validatorValue)),
			}, true
		}
	}

	if validatorName == "max" {
		validLen, err := strconv.Atoi(validatorValue)
		if err == nil && value > int64(validLen) {
			return ValidationError{
				name,
				errors.New("wrong int length, must be less than " + validatorValue),
			}, true
		}
	}

	if validatorName == "in" {
		validStrings := strings.Split(validatorValue, ",")
		if !contains(validStrings, string(rune(int(value)))) {
			return ValidationError{
				name,
				errors.New("wrong int, not mach to: " + validatorValue),
			}, true
		}
	}

	return ValidationError{}, false
}

func validateStringField(name string, validator string, value string) (ValidationError, bool) {
	validationArray := strings.Split(validator, ":")
	if len(validationArray) != 2 {
		return ValidationError{}, false
	}

	validatorName := validationArray[0]
	validatorValue := validationArray[1]

	if validatorName == "len" {
		validLen, err := strconv.Atoi(validatorValue)
		if err == nil && len(value) != validLen {
			return ValidationError{
				name,
				errors.New("wrong string length"),
			}, true
		}
	}

	if validatorName == "regexp" {
		regexpComp, _ := regexp.Compile(validatorValue)
		isValid := regexpComp.MatchString(value)
		if !isValid {
			return ValidationError{
				name,
				errors.New("wrong string, regexp error"),
			}, true
		}
	}

	if validatorName == "in" {
		validStrings := strings.Split(validatorValue, ",")

		if !contains(validStrings, value) {
			return ValidationError{
				name,
				errors.New("wrong string, not mach to: " + validatorValue),
			}, true
		}
	}

	return ValidationError{}, false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
