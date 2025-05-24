package lex

import (
	"Falcon/context"
	"Falcon/sugar"
	"strconv"
	"strings"
)

type Lexer struct {
	ctx        *context.CodeContext
	source     string
	sourceLen  int
	currIndex  int
	currColumn int
	currRow    int
	Tokens     []*Token
}

func NewLexer(ctx *context.CodeContext) *Lexer {
	return &Lexer{
		ctx:        ctx,
		source:     *ctx.SourceCode,
		sourceLen:  len(*ctx.SourceCode),
		currIndex:  0,
		currColumn: 1, // current line
		currRow:    0, // nth character of current line
		Tokens:     []*Token{},
	}
}

func (l *Lexer) Lex() []*Token {
	for l.notEOF() {
		l.parse()
	}
	return l.Tokens
}

func (l *Lexer) parse() {
	c := l.next()
	if c == '\n' {
		l.currColumn++
		l.currRow = 0
		return
	} else if c == ' ' || c == '\t' {
		return
	}
	switch c {
	case '+':
		l.createOp("+")
	case '-':
		if l.consume('>') {
			l.createOp("->")
		} else {
			l.createOp("-")
		}
	case '*':
		l.createOp("*")
	case '/':
		l.createOp("/")
	case '|':
		if l.consume('|') {
			l.createOp("||")
		} else {
			l.createOp("|")
		}
	case '&':
		if l.consume('&') {
			l.createOp("&&")
		} else {
			l.createOp("&")
		}
	case '~':
		l.createOp("~")
	case '<':
		if l.consume('=') {
			l.createOp("<=")
		} else if l.consume('<') {
			l.createOp("<<")
		} else {
			l.createOp("<")
		}
	case '>':
		if l.consume('=') {
			l.createOp(">=")
		} else if l.consume('>') {
			l.createOp(">>")
		} else {
			l.createOp(">")
		}
	case '(':
		l.createOp("(")
	case ')':
		l.createOp(")")
	case '[':
		l.createOp("[")
	case ']':
		l.createOp("]")
	case '{':
		l.createOp("{")
	case '}':
		l.createOp("}")
	case '=':
		if l.consume('=') {
			if l.consume('=') {
				l.createOp("===")
			} else {
				l.createOp("==")
			}
		} else {
			l.createOp("=")
		}
	case '.':
		l.createOp(".")
	case ',':
		l.createOp(",")
	case '?':
		l.createOp("?")
	case '!':
		if l.consume('=') {
			if l.consume('=') {
				l.createOp("!==")
			} else {
				l.createOp("!=")
			}
		} else {
			l.createOp("!")
		}
	case ':':
		if l.consume(':') {
			l.createOp("::")
		} else {
			l.createOp(":")
		}
	case '_':
		l.createOp("_")
	case '"':
		l.text()
	default:
		l.back()
		if l.isAlpha() {
			l.alpha()
		} else if l.isDigit() {
			l.numeric()
		} else {
			l.error("Unexpected character '%'", string(c))
		}
	}
}

func (l *Lexer) createOp(op string) {
	sToken, ok := Symbols[op]
	if !ok {
		l.error("Bad createOp('%')", op)
	} else {
		l.appendToken(sToken.normal(l.currColumn, l.currRow, l.ctx, op))
	}
}

func (l *Lexer) text() {
	startIndex := l.currIndex
	for l.notEOF() && l.peek() != '"' {
		l.skip()
	}
	l.eat('"')
	content := l.source[startIndex : l.currIndex-1]
	l.appendToken(&Token{Type: Text, Content: &content, Flags: []Flag{Value, ConstantValue}, Context: l.ctx})
}

func (l *Lexer) alpha() {
	startIndex := l.currIndex
	l.skip()
	for l.notEOF() && l.isAlphaNumeric() {
		l.skip()
	}
	content := l.source[startIndex:l.currIndex]
	sToken, ok := Keywords[content]
	if ok {
		l.appendToken(sToken.normal(l.currColumn, l.currRow, l.ctx))
	} else {
		l.appendToken(&Token{Type: Name, Content: &content, Flags: []Flag{Value}, Context: l.ctx})
	}
}

func (l *Lexer) numeric() {
	var numb strings.Builder
	l.writeNumeric(&numb)
	if l.notEOF() && l.peek() == '.' {
		l.skip()
		numb.WriteByte('.')
		l.writeNumeric(&numb)
	}
	content := numb.String()
	l.appendToken(&Token{Type: Number, Content: &content, Flags: []Flag{Value, ConstantValue}, Context: l.ctx})
}

func (l *Lexer) appendToken(token *Token) {
	l.Tokens = append(l.Tokens, token)
}

func (l *Lexer) writeNumeric(builder *strings.Builder) {
	startIndex := l.currIndex
	for l.notEOF() && l.isDigit() {
		l.skip()
	}
	builder.WriteString(l.source[startIndex:l.currIndex])
}

func (l *Lexer) isAlphaNumeric() bool {
	return l.isAlpha() || l.isDigit()
}

func (l *Lexer) isAlpha() bool {
	c := l.peek()
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}

func (l *Lexer) isDigit() bool {
	c := l.peek()
	return c >= '0' && c <= '9'
}

func (l *Lexer) eat(expect uint8) {
	got := l.next()
	if got != expect {
		l.error("Expected '%', but got '%'", string(expect), string(got))
	}
}

func (l *Lexer) error(message string, args ...string) {
	panic("[line " + strconv.Itoa(l.currColumn) + "] " + sugar.Format(message, args...))
}

func (l *Lexer) consume(expect uint8) bool {
	if l.peek() == expect {
		l.currIndex++
		return true
	}
	return false
}

func (l *Lexer) back() {
	l.currIndex--
	l.currRow--
}

func (l *Lexer) skip() {
	l.currIndex++
	l.currRow++
}

func (l *Lexer) peek() uint8 {
	return l.source[l.currIndex]
}

func (l *Lexer) next() uint8 {
	c := l.source[l.currIndex]
	l.currIndex++
	l.currRow++
	return c
}

func (l *Lexer) isEOF() bool {
	return l.currIndex >= l.sourceLen
}

func (l *Lexer) notEOF() bool {
	return l.currIndex < l.sourceLen
}
