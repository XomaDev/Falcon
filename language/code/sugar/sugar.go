package sugar

import (
	"strings"
)

// TODO:
// Be careful about the use of %, we have got a remainder operator too
// In future we need to replace with a more robust thingy
func Format(format string, args ...string) string {
	var builder strings.Builder
	formatLen := len(format)
	argsLen := len(args)
	argIndex := 0
	for i := 0; i < formatLen; i++ {
		if format[i] == '%' && argIndex < argsLen {
			builder.WriteString(args[argIndex])
			argIndex++
		} else {
			builder.WriteByte(format[i])
		}
	}
	return builder.String()
}

func IndexAfterNthOccurrence(s string, n int, r rune) int {
	count := 0
	for i, ch := range s {
		if ch == r {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}
