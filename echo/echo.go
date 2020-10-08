package echo

import (
	"strings"
)

// Run takes an input string, a number and string operation, and will return a string
// with the operation applied to the input string
func Run(input string, num int, op StringOperation) string {
	var parts []string
	for i := 0; i < num; i++ {
		parts = append(parts, op(input))
	}
	return strings.Join(parts, " ")
}
