package validate

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TransErr(e error) error {
	var ve validator.ValidationErrors
	if errors.As(e, &ve) {
		errMessages := make([]string, 0)
		for _, value := range ve {
			errMessages = append(errMessages, value.Translate(trans))
		}

		return errors.New(strings.Join(errMessages, ", "))
	}

	return e
}
