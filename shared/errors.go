package shared

import "fmt"

// ------------------------------------------- Debugging Errors ----------------------------------- //
type CodeError string

func (e CodeError) Error() string {
	return fmt.Sprintf("Code is wrong in [%s]", string(e))
}
type SameTimeStampError string

func (e SameTimeStampError) Error() string {
	return fmt.Sprintf("Same time stamp detected on same point - Code is wrong in [%s]", string(e))
}