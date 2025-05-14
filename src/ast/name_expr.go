package ast

import "Falcon/types"

type NameExpr struct {
	Where  types.Token
	Name   *string
	Global bool
}

func (ne *NameExpr) String() string {
	if ne.Global {
		return "glob." + *ne.Name
	}
	return *ne.Name
}

func (ne *NameExpr) Blockly() Block {
	panic("unimplemented")
}
