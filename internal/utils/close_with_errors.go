package utils

import (
	"fmt"
	"io"
)

func CloseWithErrors(filename string, errorReturn error, toClose io.Closer) (errResult error) {
	if toClose != nil {
		if tempErr := toClose.Close(); tempErr != nil {
			tempErr = fmt.Errorf("error for close %v : %v", filename, tempErr.Error())
			if errorReturn == nil {
				errResult = tempErr
			} else {
				errResult = fmt.Errorf("%v. %v", errorReturn, tempErr.Error())
			}
		}
	}
	return errResult
}
