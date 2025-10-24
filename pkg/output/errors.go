package output

import (
	"errors"
	"fmt"
)

func NewError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return errors.New(ColorRed + message + ColorReset)
}
