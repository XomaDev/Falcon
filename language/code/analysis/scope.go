package analysis

import (
	"Falcon/code/lex"
)

//go:generate stringer -type=Scope
type Scope int

const (
	ScopeRetProc Scope = iota
	ScopeProc
	ScopeGenericEvent
	ScopeEvent
	Loop
	IfBody
)

type ScopeCursor struct {
	currScopes []Scope
}

func (s *ScopeCursor) Enter(where *lex.Token, t Scope) {
	// TODO: check for bad scoping
	s.currScopes = append(s.currScopes, t)
}

func (s *ScopeCursor) Exit(t Scope) {
	topIndex := len(s.currScopes) - 1
	current := s.currScopes[topIndex]
	s.currScopes = s.currScopes[:topIndex]
	if current != t {
		panic("Bad scope exit! Expected " + current.String() + " but got " + t.String())
	}
}

func (s *ScopeCursor) In(t Scope) bool {
	for _, scope := range s.currScopes {
		if scope == t {
			return true
		}
	}
	return false
}
