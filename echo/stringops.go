package echo

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// StringOperation is a function that maps a string to another string
type StringOperation func(string) string

// NoOp returns a no-op string operation
func NoOp() StringOperation {
	return func(in string) string {
		return in
	}
}

// Reverse returns a StringOperation function that reverses the string
func Reverse() StringOperation {
	return func(in string) string {
		runes := []rune(in)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}
}

// Randomize returns a StringOperation which randomizes the string
func Randomize() StringOperation {
	return func(in string) string {
		runes := make([]rune, len(in))
		for i := range runes {
			runes[i] = rune(in[i])
		}
		for i := range runes {
			dst := rand.Intn(len(runes))
			tmp := runes[dst]
			runes[dst] = runes[i]
			runes[i] = tmp
		}
		return string(runes)
	}
}
