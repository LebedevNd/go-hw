package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				"adafsvsasqwefascsfww",
				"Arhechelovipopedryakhovsky",
				17,
				"boommail.ru",
				"Chief",
				[]string{
					"123123",
					"qweqds",
				},
				json.RawMessage{},
			},
			ValidationErrors{
				ValidationError{
					"ID",
					errors.New("wrong string length"),
				},
				ValidationError{
					"Age",
					errors.New("wrong int length, must be more or equal to 18"),
				},
				ValidationError{
					"Email",
					errors.New("wrong string, regexp error"),
				},
				ValidationError{
					"Role",
					errors.New("wrong string, not mach to: admin,stuff"),
				},
				ValidationError{
					"Phones",
					errors.New("wrong string length"),
				},
				ValidationError{
					"Phones",
					errors.New("wrong string length"),
				},
			},
		},
		{
			App{
				"Hello world!",
			},
			ValidationErrors{
				ValidationError{
					"Version",
					errors.New("wrong string length"),
				},
			},
		},
		{
			Token{
				[]byte("asd"),
				[]byte("asd"),
				[]byte("asd"),
			},
			ValidationErrors{},
		},
		{
			Response{
				66,
				"yo",
			},
			ValidationErrors{
				ValidationError{
					"Code",
					errors.New("wrong int, not mach to: 200,404,500"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			var validationErrors ValidationErrors

			ok := errors.As(Validate(tt.in), &validationErrors)
			require.True(t, ok)
			require.Equal(t, tt.expectedErr, validationErrors)
		})
	}
}
