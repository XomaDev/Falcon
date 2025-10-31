package common

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/ast/variables"
	"Falcon/code/lex"
	"Falcon/code/sugar"
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
	Args  []blockly2.Expr
}

func (f *FuncCall) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (f *FuncCall) String() string {
	return sugar.Format("%(%)", f.Name, blockly2.JoinExprs(", ", f.Args))
}

func (f *FuncCall) Blockly() blockly2.Block {
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
		return f.ctrlSimpleBlock("controls_getStartValue")
	case "closeScreen":
		return f.ctrlSimpleBlock("controls_closeScreen")
	case "closeApp":
		return f.ctrlSimpleBlock("controls_closeApplication")
	case "getPlainStartText":
		return f.ctrlSimpleBlock("controls_getPlainStartText")
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

func (f *FuncCall) Continuous() bool {
	return true
}

func (f *FuncCall) Consumable() bool {
	// TODO: in the long run, use a signature based model, maybe like that of the call.go
	if f.Name == "setRandSeed" || f.Name == "println" ||
		f.Name == "openScreen" || f.Name == "openScreenWithValue" ||
		f.Name == "closeScreen" || f.Name == "closeScreenWithValue" ||
		f.Name == "closeApp" || f.Name == "closeScreenWithPlainText" || f.Name == "set" {
		return false
	}
	return true
}

func (f *FuncCall) everyComponent() blockly2.Block {
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for every() 1st argument!")
	}
	return blockly2.Block{
		Type:     "component_all_component_block",
		Mutation: &blockly2.Mutation{ComponentType: compType.Name},
		Fields:   []blockly2.Field{{Name: "COMPONENT_SELECTOR", Value: compType.Name}},
	}
}

func (f *FuncCall) genericCall() blockly2.Block {
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
	return blockly2.Block{
		Type: "component_method",
		Mutation: &blockly2.Mutation{
			MethodName:    vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Values: blockly2.ValueArgsByPrefix(f.Args[1], "COMPONENT", "ARG", f.Args[3:]),
	}
}

func (f *FuncCall) genericGet() blockly2.Block {
	f.assertArgLen(3)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for get() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for get() 3rd argument!")
	}
	return blockly2.Block{
		Type: "component_set_get",
		Mutation: &blockly2.Mutation{
			SetOrGet:      "get",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []blockly2.Field{{Name: "PROP", Value: vGet.Name}},
		Values: []blockly2.Value{{Name: "COMPONENT", Block: f.Args[1].Blockly()}},
	}
}

func (f *FuncCall) genericSet() blockly2.Block {
	f.assertArgLen(4)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for set() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for set() 3rd argument!")
	}
	return blockly2.Block{
		Type: "component_set_get",
		Mutation: &blockly2.Mutation{
			SetOrGet:      "set",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []blockly2.Field{{Name: "PROP", Value: vGet.Name}},
		Values: blockly2.MakeValues([]blockly2.Expr{f.Args[1], f.Args[3]}, "COMPONENT", "VALUE"),
	}
}

func (f *FuncCall) splitColor() blockly2.Block {
	return blockly2.Block{
		Type:   "color_make_color",
		Values: blockly2.MakeValues(f.Args, "COLOR"),
	}
}

func (f *FuncCall) makeColor() blockly2.Block {
	return blockly2.Block{
		Type:   "color_make_color",
		Values: blockly2.MakeValues(f.Args, "COLORLIST"),
	}
}

func (f *FuncCall) copyDict() blockly2.Block {
	return blockly2.Block{
		Type:   "dictionaries_copy",
		Values: blockly2.MakeValues(f.Args, "DICT"),
	}
}

func (f *FuncCall) copyList() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_copy",
		Values: blockly2.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) ctrlSimpleBlock(blockType string) blockly2.Block {
	return blockly2.Block{Type: blockType}
}

func (f *FuncCall) closeScreenWithPlainText() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "controls_closeScreenWithPlainText",
		Values: blockly2.MakeValues(f.Args, "TEXT"),
	}
}

func (f *FuncCall) closeScreenWithValue() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "controls_closeScreenWithValue",
		Values: blockly2.MakeValues(f.Args, "SCREEN"),
	}
}

func (f *FuncCall) openScreenWithValue() blockly2.Block {
	f.assertArgLen(2)
	return blockly2.Block{
		Type:   "controls_openAnotherScreenWithStartValue",
		Values: blockly2.MakeValues(f.Args, "SCREENNAME", "STARTVALUE"),
	}
}

func (f *FuncCall) openScreen() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "controls_openAnotherScreen",
		Values: blockly2.MakeValues(f.Args, "SCREEN"),
	}
}

func (f *FuncCall) println() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "controls_eval_but_ignore",
		Values: blockly2.MakeValues(f.Args, "VALUE"),
	}
}

func (f *FuncCall) mathConversions() blockly2.Block {
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
	return blockly2.Block{
		Type:   blockType,
		Fields: []blockly2.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly2.Value{{Name: "NUM", Block: f.Args[0].Blockly()}},
	}
}

func (f *FuncCall) formatDecimal() blockly2.Block {
	f.assertArgLen(2)
	return blockly2.Block{
		Type:   "math_format_as_decimal",
		Values: blockly2.MakeValues(f.Args, "NUM", "PLACES"),
	}
}

func (f *FuncCall) atan2() blockly2.Block {
	f.assertArgLen(2)
	return blockly2.Block{
		Type:   "math_atan2",
		Values: blockly2.MakeValues(f.Args, "Y", "X"),
	}
}

func (f *FuncCall) mathDivide() blockly2.Block {
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
	return blockly2.Block{
		Type:   "math_divide",
		Fields: []blockly2.Field{{Name: "OP", Value: fieldOp}},
		Values: blockly2.MakeValues(f.Args, "DIVIDEND", "DIVISOR"),
	}
}

func (f *FuncCall) modeOf() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "math_mode_of_list",
		Values: blockly2.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) mathOnList() blockly2.Block {
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
	return blockly2.Block{
		Type:   "math_on_list2",
		Fields: []blockly2.Field{{Name: "OP", Value: fieldOp}},
		Values: blockly2.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) minOrMax() blockly2.Block {
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
	return blockly2.Block{
		Type:     "math_on_list",
		Fields:   []blockly2.Field{{Name: "OP", Value: fieldOp}},
		Mutation: &blockly2.Mutation{ItemCount: argSize},
		Values:   blockly2.ValuesByPrefix("NUM", f.Args),
	}
}

func (f *FuncCall) setRandSeed() blockly2.Block {
	f.assertArgLen(1)
	return blockly2.Block{
		Type:   "math_random_set_seed",
		Values: blockly2.MakeValues(f.Args, "NUM"),
	}
}

func (f *FuncCall) randFloat() blockly2.Block {
	f.assertArgLen(0)
	return blockly2.Block{Type: "math_random_float"}
}

func (f *FuncCall) randInt() blockly2.Block {
	f.assertArgLen(2)
	return blockly2.Block{
		Type:   "math_random_int",
		Values: blockly2.MakeValues(f.Args, "FROM", "TO"),
	}
}

func (f *FuncCall) mathRadix() blockly2.Block {
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
	return blockly2.Block{
		Type: "math_number_radix",
		Fields: []blockly2.Field{
			{Name: "OP", Value: fieldOp},
			{Name: "NUM", Value: textExpr.Content},
		},
	}
}

func (f *FuncCall) assertArgLen(expectLen int) {
	argsLen := len(f.Args)
	if argsLen != expectLen {
		f.Where.Error("Expected % argument for %() but got %", strconv.Itoa(expectLen), f.Name, strconv.Itoa(argsLen))
	}
}
