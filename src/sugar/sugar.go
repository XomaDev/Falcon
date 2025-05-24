package sugar

import (
	"strings"
)

func Format(format string, args ...string) string {
	var builder strings.Builder
	formatLen := len(format)
	argsLen := len(args)
	argIndex := 0
	for i := 0; i < formatLen; i++ {
		if format[i] == '%' {
			builder.WriteString(args[argIndex])
			argIndex++
			if argIndex == argsLen {
				break
			}
		} else {
			builder.WriteByte(format[i])
		}
	}
	return builder.String()
}
