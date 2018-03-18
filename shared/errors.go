package shared

import "fmt"

// ------------------------------------------- Debugging Errors ----------------------------------- //
type CodeError string

func (e CodeError) Error() string {
	return fmt.Sprintf("Code is wrong in [%s]", string(e))
}