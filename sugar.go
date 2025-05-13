package main

import "strings"

func Format(format string, args ...string) string {
	var builder strings.Builder
	formatLen := len(format)
	argIndex := 0
	for i := 0; i < formatLen; i++ {
		if format[i] == '%' {
			builder.WriteString(args[argIndex])
			argIndex++
		} else {
			builder.WriteByte(format[i])
		}
	}
	return builder.String()
}
