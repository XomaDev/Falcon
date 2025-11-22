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
	ScopeLoop
	ScopeIfBody
	ScopeSmartBody
)

type ScopeCursor struct {
	currScopes []Scope
}

func (s *ScopeCursor) Enter(where *lex.Token, t Scope) {
	s.checkScope(where, t)
	s.currScopes = append(s.currScopes, t)
}

func (s *ScopeCursor) checkScope(where *lex.Token, t Scope) {
	depth := len(s.currScopes)
	if t == ScopeRetProc || t == ScopeProc {
		if depth != 0 {
			where.Error("Functions can only be defined at the root.")
		}
	} else if t == ScopeGenericEvent || t == ScopeEvent {
		if depth != 0 {
			where.Error("Events can only be defined at the root.")
		}
	}
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

func (s *ScopeCursor) AtRoot() bool {
	return len(s.currScopes) == 0
}
