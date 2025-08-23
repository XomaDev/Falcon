package common

import (
	"Falcon/ast/blockly"
	"Falcon/ast/fundamentals"
	"Falcon/ast/variables"
	"Falcon/lex"
	"Falcon/sugar"
	"strconv"
)

var mathConversions = map[string]string{
	"root":     "ROOT",
	"abs":      "ABS",
	"neg":      "NEG",
	"log":      "LN",
	"exp":      "EXP",
	"round":    "ROUND",
	"ceil":     "CEILING",
	"floor":    "FLOOR",
	"sin":      "SIN",
	"cos":      "COS",
	"tan":      "TAN",
	"asin":     "ASIN",
	"acos":     "ACOS",
	"atan":     "ATAN",
	"degrees":  "RADIANS_TO_DEGREES",
	"radians":  "DEGREES_TO_RADIANS",
	"decToHex": "DEC_TO_HEX",
	"decToBin": "DEC_TO_BIN",
	"hexToDec": "HEX_TO_DEC",
	"binToDec": "BIN_TO_DEC",
}

type FuncCall struct {
	Where *lex.Token
	Name  string
	Args  []blockly.Expr
}

func (f *FuncCall) String() string {
	return sugar.Format("%(%)", f.Name, blockly.JoinExprs(", ", f.Args))
}

func (f *FuncCall) Blockly() blockly.Block {
	switch f.Name {
	case "root", "abs", "neg", "log", "exp", "round", "ceil", "floor",
		"sin", "cos", "tan", "asin", "acos", "atan", "degrees", "radians",
		"decToHex", "decToBin", "hexToDec", "binToDec":
		return f.mathConversions()

	case "dec", "bin", "octal", "hexa":
		return f.mathRadix()
	case "randInt":
		return f.randInt()
	case "randFloat":
		return f.randFloat()
	case "setRandSeed":
		return f.setRandSeed()
	case "min", "max":
		return f.minOrMax()
	case "avgOf", "maxOf", "minOf", "geoMeanOf", "stdDevOf", "stdErrOf":
		return f.mathOnList()
	case "modeOf":
		return f.modeOf()
	case "mod", "rem", "quot":
		return f.mathDivide()
	case "aTan2":
		return f.atan2()
	case "formatDecimal":
		return f.formatDecimal()

	case "println":
		return f.println()
	case "openScreen":
		return f.openScreen()
	case "openScreenWithValue":
		return f.openScreenWithValue()
	case "closeScreenWithValue":
		return f.closeScreenWithValue()
	case "getStartValue":
		return f.ctrlSimpleBlock("controls_getStartValue", true)
	case "closeScreen":
		return f.ctrlSimpleBlock("controls_closeScreen", false)
	case "closeApp":
		return f.ctrlSimpleBlock("controls_closeApplication", false)
	case "getPlainStartText":
		return f.ctrlSimpleBlock("controls_getPlainStartText", true)
	case "closeScreenWithPlainText":
		return f.closeScreenWithPlainText()
	case "copyList":
		return f.copyList()
	case "copyDict":
		return f.copyDict()

	case "makeColor":
		return f.makeColor()
	case "splitColor":
		return f.splitColor()

	case "set":
		return f.genericSet()
	case "get":
		return f.genericGet()
	case "call":
		return f.genericCall()
	case "every":
		return f.everyComponent()
	default:
		panic("Unknown function " + f.Name)
	}
}

func (f *FuncCall) everyComponent() blockly.Block {
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for every() 1st argument!")
	}
	return blockly.Block{
		Type:     "component_all_component_block",
		Mutation: &blockly.Mutation{ComponentType: compType.Name},
		Fields:   []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: compType.Name}},
	}
}

func (f *FuncCall) genericCall() blockly.Block {
	// arg[0] 	 compType
	// arg[1] 	 component (any object)
	// arg[2] 	 method name
	// arg[4->n] invoke args
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for call() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a method name for call() 3rd argument!")
	}
	return blockly.Block{
		Type: "component_method",
		Mutation: &blockly.Mutation{
			MethodName:    vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Values: blockly.ValueArgsByPrefix(f.Args[1], "COMPONENT", "ARG", f.Args[3:]),
	}
}

func (f *FuncCall) genericGet() blockly.Block {
	f.assertArgLen(3)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for get() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for get() 3rd argument!")
	}
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "get",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []blockly.Field{{Name: "PROP", Value: vGet.Name}},
		Values: []blockly.Value{{Name: "COMPONENT", Block: f.Args[1].Blockly()}},
	}
}

func (f *FuncCall) genericSet() blockly.Block {
	f.assertArgLen(4)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for set() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for set() 3rd argument!")
	}
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "set",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []blockly.Field{{Name: "PROP", Value: vGet.Name}},
		Values: blockly.MakeValues([]blockly.Expr{f.Args[1], f.Args[3]}, "COMPONENT", "VALUE"),
	}
}

func (f *FuncCall) Continuous() bool {
	return true
}

func (f *FuncCall) splitColor() blockly.Block {
	return blockly.Block{
		Type:       "color_make_color",
		Values:     blockly.MakeValues(f.Args, "COLOR"),
		Consumable: true,
	}
}

func (f *FuncCall) makeColor() blockly.Block {
	return blockly.Block{
		Type:       "color_make_color",
		Values:     blockly.MakeValues(f.Args, "COLORLIST"),
		Consumable: true,
	}
}

func (f *FuncCall) copyDict() blockly.Block {
	return blockly.Block{
		Type:       "dictionaries_copy",
		Values:     blockly.MakeValues(f.Args, "DICT"),
		Consumable: true,
	}
}

func (f *FuncCall) copyList() blockly.Block {
	return blockly.Block{
		Type:       "lists_copy",
		Values:     blockly.MakeValues(f.Args, "LIST"),
		Consumable: true,
	}
}

func (f *FuncCall) ctrlSimpleBlock(blockType string, consumable bool) blockly.Block {
	return blockly.Block{Type: blockType, Consumable: consumable}
}

func (f *FuncCall) closeScreenWithPlainText() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "controls_closeScreenWithPlainText",
		Values:     blockly.MakeValues(f.Args, "TEXT"),
		Consumable: false,
	}
}

func (f *FuncCall) closeScreenWithValue() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "controls_closeScreenWithValue",
		Values:     blockly.MakeValues(f.Args, "SCREEN"),
		Consumable: false,
	}
}

func (f *FuncCall) openScreenWithValue() blockly.Block {
	f.assertArgLen(2)
	return blockly.Block{
		Type:       "controls_openAnotherScreenWithStartValue",
		Values:     blockly.MakeValues(f.Args, "SCREENNAME", "STARTVALUE"),
		Consumable: false,
	}
}

func (f *FuncCall) openScreen() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "controls_openAnotherScreen",
		Values:     blockly.MakeValues(f.Args, "SCREEN"),
		Consumable: false,
	}
}

func (f *FuncCall) println() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "controls_eval_but_ignore",
		Values:     blockly.MakeValues(f.Args, "VALUE"),
		Consumable: false,
	}
}

func (f *FuncCall) mathConversions() blockly.Block {
	f.assertArgLen(1)
	fieldOp, ok := mathConversions[f.Name]
	if !ok {
		f.Where.Error("Unknown Math Conversion %()", f.Name)
	}
	var blockType string
	switch fieldOp {
	case "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN":
		blockType = "math_trig"
	case "RADIANS_TO_DEGREES", "DEGREES_TO_RADIANS":
		blockType = "math_convert_angles"
	case "DEC_TO_HEX", "HEX_TO_DEC", "DEC_TO_BIN", "BIN_TO_DEC":
		blockType = "math_convert_number"
	default:
		blockType = "math_single"
	}
	return blockly.Block{
		Type:       blockType,
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values:     []blockly.Value{{Name: "NUM", Block: f.Args[0].Blockly()}},
		Consumable: true,
	}
}

func (f *FuncCall) formatDecimal() blockly.Block {
	f.assertArgLen(2)
	return blockly.Block{
		Type:       "math_format_as_decimal",
		Values:     blockly.MakeValues(f.Args, "NUM", "PLACES"),
		Consumable: true,
	}
}

func (f *FuncCall) atan2() blockly.Block {
	f.assertArgLen(2)
	return blockly.Block{
		Type:       "math_atan2",
		Values:     blockly.MakeValues(f.Args, "Y", "X"),
		Consumable: true,
	}
}

func (f *FuncCall) mathDivide() blockly.Block {
	f.assertArgLen(2)
	var fieldOp string
	switch f.Name {
	case "mod":
		fieldOp = "MODULO"
	case "rem":
		fieldOp = "REMAINDER"
	case "quot":
		fieldOp = "QUOTIENT"
	}
	return blockly.Block{
		Type:       "math_divide",
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values:     blockly.MakeValues(f.Args, "DIVIDEND", "DIVISOR"),
		Consumable: true,
	}
}

func (f *FuncCall) modeOf() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "math_mode_of_list",
		Values:     blockly.MakeValues(f.Args, "LIST"),
		Consumable: true,
	}
}

func (f *FuncCall) mathOnList() blockly.Block {
	f.assertArgLen(1)
	var fieldOp string
	switch f.Name {
	case "avgOf":
		fieldOp = "AVG"
	case "maxOf":
		fieldOp = "MAX"
	case "minOf":
		fieldOp = "MIN"
	case "geoMeanOf":
		fieldOp = "GM"
	case "stdDevOf":
		fieldOp = "SD"
	case "stdErrOf":
		fieldOp = "SE"
	}
	return blockly.Block{
		Type:       "math_on_list2",
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values:     blockly.MakeValues(f.Args, "LIST"),
		Consumable: true,
	}
}

func (f *FuncCall) minOrMax() blockly.Block {
	argSize := len(f.Args)
	if argSize == 0 {
		f.Where.Error("No arguments provided for %()", f.Name)
	}
	var fieldOp string
	switch f.Name {
	case "min":
		fieldOp = "MIN"
	case "max":
		fieldOp = "MAX"
	}
	return blockly.Block{
		Type:       "math_on_list",
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Mutation:   &blockly.Mutation{ItemCount: argSize},
		Values:     blockly.ValuesByPrefix("NUM", f.Args),
		Consumable: true,
	}
}

func (f *FuncCall) setRandSeed() blockly.Block {
	f.assertArgLen(1)
	return blockly.Block{
		Type:       "math_random_set_seed",
		Values:     blockly.MakeValues(f.Args, "NUM"),
		Consumable: false,
	}
}

func (f *FuncCall) randFloat() blockly.Block {
	f.assertArgLen(0)
	return blockly.Block{Type: "math_random_float", Consumable: true}
}

func (f *FuncCall) randInt() blockly.Block {
	f.assertArgLen(2)
	return blockly.Block{
		Type:       "math_random_int",
		Values:     blockly.MakeValues(f.Args, "FROM", "TO"),
		Consumable: true,
	}
}

func (f *FuncCall) mathRadix() blockly.Block {
	f.assertArgLen(1)
	var fieldOp string
	switch f.Name {
	case "dec":
		fieldOp = "DEC"
	case "bin":
		fieldOp = "BIN"
	case "octal":
		fieldOp = "OCT"
	case "hexa":
		fieldOp = "HEX"
	}
	textExpr, ok := f.Args[0].(*fundamentals.Text)
	if !ok {
		f.Where.Error("Expected a numeric string argument for %()", f.Name)
	}
	return blockly.Block{
		Type: "math_number_radix",
		Fields: []blockly.Field{
			{Name: "OP", Value: fieldOp},
			{Name: "NUM", Value: textExpr.Content},
		},
		Consumable: true,
	}
}

func (f *FuncCall) assertArgLen(expectLen int) {
	argsLen := len(f.Args)
	if argsLen != expectLen {
		f.Where.Error("Expected % argument for %() but got %", strconv.Itoa(expectLen), f.Name, strconv.Itoa(argsLen))
	}
}
