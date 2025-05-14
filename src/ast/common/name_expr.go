package common

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
)

type Name struct {
	Where  lex.Token
	Name   string
	Global bool
}

func (ne *Name) String() string {
	if ne.Global {
		return "glob." + ne.Name
	}
	return ne.Name
}

func (ne *Name) Blockly() blockly.Block {
	panic("unimplemented")
}
