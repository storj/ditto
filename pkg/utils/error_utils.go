package utils

import (
	"errors"
	"strings"
)

//CombineErrors will combine all error messages from error array in a single error
//Params:
//erray     - slice of errors to combine
//Returns   - new error
func CombineErrors(erray []error) (err error) {
	length := len(erray)

	if erray == nil || length == 0 {
		return nil
	}

	var errorStrings []string

	for i := 0; i < length; i++ {
		if erray[i] != nil {
			errorStrings = append(errorStrings, erray[i].Error())
		}
	}

	if len(errorStrings) > 0 {
		err = errors.New(strings.Join(errorStrings, "\n"))
	}

	return
}

//NewError  - extends error with additional message
//Params:
//err       - error to extend
//message   - message to extend error
//Returns   - new error
func NewError(err error, message string) error {

	if err == nil {
		return nil
	}

	return errors.New(message + err.Error())
}