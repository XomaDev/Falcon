package context

import (
	"Falcon/sugar"
	"strings"
)

type CodeContext struct {
	SourceCode *string
	FileName   string
}

func (c *CodeContext) ReportError(
	column int,
	row int,
	highlightWordSize int,
	message string,
	args ...string,
) {
	code := *c.SourceCode
	beginOfLine := sugar.IndexAfterNthOccurrence(code, column-1, '\n') + 1
	endOfLine := strings.Index(code[beginOfLine:], "\n")
	line := code[beginOfLine:max(beginOfLine+endOfLine, len(code))]

	var builder strings.Builder
	boxTop := strings.Repeat(".", len(line))

	builder.WriteByte('\n')
	builder.WriteString(boxTop)
	builder.WriteByte('\n')
	builder.WriteString(line)
	builder.WriteByte('\n')
	builder.WriteString(strings.Repeat(" ", row-highlightWordSize))
	builder.WriteString(strings.Repeat("^", highlightWordSize))
	builder.WriteByte('\n')
	builder.WriteString(sugar.Format(message, args...))
	builder.WriteByte('\n')
	builder.WriteString(boxTop)
	panic(builder.String())
}
