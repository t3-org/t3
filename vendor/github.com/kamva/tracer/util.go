package tracer

import "fmt"

// StackAsString returns error's stack(if exists) as string.
func StackAsString(err error) string {
	if te, ok := err.(StackTracer); err != nil && ok {
		return fmt.Sprintf("%+v", te)
	}
	return ""
}
