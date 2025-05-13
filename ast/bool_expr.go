package ast

type BoolExpr struct {
	Value *string
}

func (b *BoolExpr) String() string {
	return *b.Value
}

func (b *BoolExpr) Blockly() Block {
	var bText string
	if *b.Value == "true" {
		bText = "TRUE"
	} else {
		bText = "FALSE"
	}
	return Block{
		Type:   "logic_boolean",
		Fields: ToFields(map[string]string{"BOOL": bText}),
	}
}
