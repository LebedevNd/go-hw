package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

//TODO создать свой тип ошибки, реализовав тип ValidationErrors

func Validate(v interface{}) error {
	var errors = ValidationErrors{}

	refValueStruct := reflect.ValueOf(v)
	refTypeStruct := refValueStruct.Type()
	numFields := refTypeStruct.NumField()

	for i := 0; i < numFields; i++ {
		refFieldType := refTypeStruct.Field(i)
		refFieldValue := refValueStruct.Field(i)
		fieldType := refFieldValue.Type()
		fieldTag, found := refFieldType.Tag.Lookup("validate")
		fieldName := refFieldType.Name

		if found {
			fmt.Println(fieldType)

			validators := strings.Split(fieldTag, "|")

			for _, validator := range validators {
				err := validateField(fieldName, fieldType, validator)
				errors = append(errors, err)
			}
		}
	}

	return nil
}
