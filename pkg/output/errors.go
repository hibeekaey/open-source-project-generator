package output

import (
	"errors"
	"fmt"
)

func NewError(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return errors.New(ColorRed + message + ColorReset)
}
