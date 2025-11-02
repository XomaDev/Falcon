package components

import (
	"Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type PropertySet struct {
	ComponentName string
	ComponentType string
	Property      string
	Value         ast.Expr
}

func (p *PropertySet) Yail() string {
	yail := "(set-and-coerce-property! '"
	yail += p.ComponentType
	yail += " '"
	yail += p.Property
	yail += " "
	yail += p.Value.Yail()
	yail += " '"
	yail += ")"
	return yail
}

func (p *PropertySet) String() string {
	return sugar.Format("%.% = %", p.ComponentName, p.Property, p.Value.String())
}

func (p *PropertySet) Blockly() ast.Block {
	return ast.Block{
		Type: "component_set_get",
		Mutation: &ast.Mutation{
			SetOrGet:      "set",
			PropertyName:  p.Property,
			IsGeneric:     false,
			InstanceName:  p.ComponentName,
			ComponentType: p.ComponentType,
		},
		Fields: ast.FieldsFromMap(map[string]string{
			"COMPONENT_SELECTOR": p.ComponentName,
			"PROP":               p.Property,
		}),
		Values: []ast.Value{{Name: "VALUE", Block: p.Value.Blockly()}},
	}
}

func (p *PropertySet) Continuous() bool {
	return false
}

func (p *PropertySet) Consumable() bool {
	return false
}
