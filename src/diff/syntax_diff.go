package diff

import (
	"Falcon/context"
	"Falcon/lex"
	"Falcon/sugar"
	"strings"
)

type SyntaxDiff struct {
	humanSyntax   string
	machineSyntax string

	humanWords   []string
	machineWords []string

	hSize int
	mSize int

	hIndex int
	mIndex int

	mergedWords strings.Builder
}

func MakeSyntaxDiff(humanSyntax string, machineSyntax string) *SyntaxDiff {
	humanWords := makeSimpleWords(humanSyntax, lex.NewLexer(&context.CodeContext{SourceCode: &humanSyntax, FileName: "human.diff"}).Lex())
	machineWords := makeSimpleWords(machineSyntax, lex.NewLexer(&context.CodeContext{SourceCode: &machineSyntax, FileName: "machine.diff"}).Lex())

	return &SyntaxDiff{
		humanSyntax:   humanSyntax,
		machineSyntax: machineSyntax,
		humanWords:    humanWords,
		machineWords:  machineWords,
		hSize:         len(humanWords),
		mSize:         len(machineWords),
		hIndex:        0,
		mIndex:        0,
		mergedWords:   strings.Builder{},
	}
}

func makeSimpleWords(sourceCode string, tokens []*lex.Token) []string {
	var words []string

	var lastTokenColumn = -1
	var lastTokenRow = 1

	for _, token := range tokens {
		var tokenContent string
		if token.Type == lex.Text {
			tokenContent = "\"" + *token.Content + "\""
		} else {
			tokenContent = *token.Content
		}
		if lastTokenColumn != -1 {
			if lastTokenColumn != token.Column {
				separationLines := token.Column - lastTokenColumn
				words = append(words, strings.Repeat("\n", separationLines))

				afterNthLine := sugar.IndexAfterNthOccurrence(sourceCode, token.Column-1, '\n') + 1
				beforeNthSpace := afterNthLine + token.Row - len(tokenContent)
				words = append(words, sourceCode[afterNthLine:beforeNthSpace])
			} else {
				afterNthLine := sugar.IndexAfterNthOccurrence(sourceCode, token.Column-1, '\n') + 1
				afterNthRow := afterNthLine + lastTokenRow
				beforeNthRow := afterNthLine + token.Row - len(tokenContent)
				words = append(words, sourceCode[afterNthRow:beforeNthRow])
			}
		}
		words = append(words, tokenContent)

		lastTokenColumn = token.Column
		lastTokenRow = token.Row
	}
	return words
}

func (s *SyntaxDiff) Merge() string {
	for s.mIndex < s.mSize && s.hIndex < s.hSize {
		// Keep appending until mismatch
		for {
			if !s.skipSpaces(true) {
				break
			}
			hToken := s.humanWords[s.hIndex]
			mToken := s.machineWords[s.mIndex]
			if hToken != mToken {
				break
			}
			s.mergedWords.WriteString(hToken)

			s.hIndex++
			s.mIndex++
		}

		// Keep appending machineWords until match
		for s.mIndex < s.mSize {
			if !s.skipSpaces(false) {
				break
			}
			hToken := s.humanWords[s.hIndex]
			mToken := s.machineWords[s.mIndex]

			if hToken == mToken {
				break
			}
			s.mergedWords.WriteString(mToken)

			s.mIndex++
		}
	}
	for s.mIndex < s.mSize {
		mToken := s.machineWords[s.mIndex]
		s.mergedWords.WriteString(mToken)
		s.mIndex++
	}
	return s.mergedWords.String()
}

func (s *SyntaxDiff) skipSpaces(appendHumanSpace bool) bool {
	// Skip all spaces in humanWords
	for s.hIndex < s.hSize {
		hWord := s.humanWords[s.hIndex]
		if strings.TrimSpace(hWord) == "" {
			if appendHumanSpace {
				s.mergedWords.WriteString(hWord)
			}
			s.hIndex++
		} else {
			break
		}
	}
	// Skip all spaces in machineWords
	for s.mIndex < s.mSize {
		mWord := s.machineWords[s.mIndex]
		if strings.TrimSpace(mWord) == "" {
			if !appendHumanSpace {
				s.mergedWords.WriteString(mWord)
			}
			s.mIndex++
		} else {
			break
		}
	}
	return s.hIndex < s.hSize && s.mIndex < s.mSize
}
