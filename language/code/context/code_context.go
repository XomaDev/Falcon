package context

import (
	"Falcon/code/sugar"
	"strconv"
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
	err := sugar.Format(message, args...) + "\n[line " + strconv.Itoa(column) + "]"
	code := *c.SourceCode
	beginOfLine := sugar.IndexAfterNthOccurrence(code, column-1, '\n') + 1
	endOfLine := strings.Index(code[beginOfLine:], "\n")

	if endOfLine == -1 {
		endOfLine = len(code) - beginOfLine
	}

	line := code[beginOfLine : beginOfLine+endOfLine]

	var builder strings.Builder
	boxTop := strings.Repeat(".", max(len(line), len(err)))

	builder.WriteByte('\n')
	builder.WriteString(boxTop)
	builder.WriteByte('\n')
	builder.WriteString(line)
	builder.WriteByte('\n')
	builder.WriteString(strings.Repeat(" ", row-highlightWordSize))
	builder.WriteString(strings.Repeat("^", highlightWordSize))
	builder.WriteByte('\n')
	builder.WriteString(err)
	builder.WriteByte('\n')
	builder.WriteString(boxTop)
	panic(builder.String())
}
