package parser

import (
	"fmt"
	"testing"

	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let i = 5;
	let j = 7;
	let testval = 27;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		input              string
		expectedItentifier string
		expectedVal        interface{}
	}{
		{"let g = 19;", "g", 19},
		{"let f = false;", "f", false},
		{"let fenerbahce = goat", "fenerbahce", "goat"},
	}

	for _, tcase := range tests {
		lex := lexer.New(tcase.input)
		parser := New(lex)
		program = parser.ParseProgram()
		checkParserErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("Number of Program statements not equal to 1. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tcase.expectedItentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpr(t, val, tcase.expectedVal) {
			return
		}
	}
}

func TestAssignmentStatements(t *testing.T) {

	input := `let x = 5; x = 10;`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 2 {
		t.Fatalf("Program does not contain 2 Statements. got = %d", len(program.Statements))
	}

	assStmt, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("Statement is not AssignmentStatement. got=%T", program.Statements[1])
	}
	assName, ok := assStmt.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("Left of Assignment is not Identifier. got=%T", assStmt.Left)
	}
	if assName.Value != "x" {
		t.Fatalf("Assignment Name not 'x'. got = %s", assName.Value)
	}
	if !testLiteralExpr(t, assStmt.Value, 10) {
		return
	}
}

func TestAssignmentStatementsWithPrefixExpression(t *testing.T) {
	tests := []struct {
		input        string
		expectedLeft string
		expectedVal  interface{}
	}{
		{
			input:        "*x = 5;",
			expectedLeft: "(*x)",
			expectedVal:  5,
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("[%d] program.Statements does not contain 1 statement. got=%d",
				i, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("[%d] stmt is not *ast.AssignmentStatement. got=%T",
				i, program.Statements[0])
		}

		if stmt.Left.String() != tt.expectedLeft {
			t.Errorf("[%d] left side String() wrong. expected=%q, got=%q",
				i, tt.expectedLeft, stmt.Left.String())
		}

		testLiteralExpr(t, stmt.Value, tt.expectedVal)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal interface{}
	}{
		{input: "return 1907", expectedVal: 1907},
		{input: "return fenerbahce", expectedVal: "fenerbahce"},
		{input: "return false", expectedVal: false},
	}

	for _, tcase := range tests {
		lex := lexer.New(tcase.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Number of Program statements not 1. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		ret, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt is not a return Statement but a %T", stmt)
		}
		if ret.TokenLiteral() != "return" {
			t.Fatalf("Wrong Token Literal for return statement. got=%q", ret.TokenLiteral())
		}
		if !testLiteralExpr(t, ret.ReturnValue, tcase.expectedVal) {
			return
		}
	}
}

func TestIdentifierExpr(t *testing.T) {
	input := "thingy"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkParserErrors(t, pars)
	if len(program.Statements) != 1 {
		t.Fatalf("Program has more or less than 1 Statement. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement. got=%T", program.Statements[0])
	}
	id, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Statement Expression is not of Type Identifier. got=%T", stmt.Expression)
	}
	if id.Value != "thingy" {
		t.Fatalf("Identifier Value is not thingy. got=%s", id.Value)
	}
	if id.TokenLiteral() != "thingy" {
		t.Fatalf("Token Literal is not thingy. got=%s", id.TokenLiteral())
	}

}

func TestIntegerExpr(t *testing.T) {
	input := "1907"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkParserErrors(t, pars)
	if len(program.Statements) != 1 {
		t.Fatalf("Program has more or less than 1 Statement. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement. got=%T", program.Statements[0])
	}
	inte, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement Expression is not of Type Integer. got=%T", stmt.Expression)
	}
	if inte.Value != 1907 {
		t.Fatalf("Integer Value is not thingy. got=%d", inte.Value)
	}
	if inte.TokenLiteral() != "1907" {
		t.Fatalf("Token Literal is not %d. got=%s", 5, inte.TokenLiteral())
	}

}

func TestStringExpr(t *testing.T) {
	input := `"whats up";`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("Number of ProgramStatements != 1. got=%d", len(program.Statements))
	}
	strStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not Expressionstatement. got=%T(%v)", program.Statements[0], program.Statements[0])
	}
	strLit, ok := strStmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expression isn't string literal. got=%T(%v)", strStmt, strStmt)
	}
	if strLit.Value != "whats up" {
		t.Fatalf("Wrong String value. Expected= whats up, got=%q", strLit.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	Tests := []struct {
		inp          string
		expectedBool bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range Tests {
		lex := lexer.New(tt.inp)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBool {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBool,
				boolean.Value)
		}
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (3 < 4) { 5 * 5; }`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	while, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpr(t, while.Condition, 3, "<", 4) {
		return
	}

	if len(while.Body.Statements) != 1 {
		t.Fatalf("Wrong number of Statements in Loop Body. got=%d", len(while.Body.Statements))
	}

	body, ok := while.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong Statementtype for Statements[0] expected=ExpressionStatement. got=%T", while.Body.Statements[0])
	}
	if !testInfixExpr(t, body.Expression, 5, "*", 5) {
		t.Fatalf("Wrong Expressiontype in Loop Body. expected 5 * 5 got=%s", body.Expression.String())
	}

}

func TestIfExpression(t *testing.T) {
	input := `if (f > b) { g }`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	iff, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Wrong Expression type. Expected IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpr(t, iff.Condition, "f", ">", "b") {
		return
	}

	if len(iff.Then.Statements) != 1 {
		t.Fatalf("Wrong number of Then Statements. got=%d", len(iff.Then.Statements))
	}

	then, ok := iff.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong Statementtype for Statements[0] expected=ExpressionStatement. got=%T", iff.Then.Statements[0])
	}

	if !testIdentifier(t, then.Expression, "g") {
		return
	}

	if iff.Alt != nil {
		t.Errorf("iff.Alt.Statements is not nil. got=%v", iff.Alt)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (f > b) { g } else { h }`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	iff, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Wrong Expression type. Expected IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpr(t, iff.Condition, "f", ">", "b") {
		return
	}

	if len(iff.Then.Statements) != 1 {
		t.Fatalf("Wrong number of Then Statements. got=%d", len(iff.Then.Statements))
	}

	then, ok := iff.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong Statementtype for iff.then.Statements[0] expected=ExpressionStatement. got=%T", iff.Then.Statements[0])
	}

	if !testIdentifier(t, then.Expression, "g") {
		return
	}

	alt, ok := iff.Alt.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong Statementtype for iff.alt.Statements[0] expected=ExpressionStatement. got=%T", iff.Alt.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "h") {
		return
	}
}

func TestFunctionExpr(t *testing.T) {
	input := `fnc (f, b) { f * b; }`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("number of statements in program body not 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	fnc, ok := stmt.Expression.(*ast.FuncLiteral)
	if !ok {
		t.Fatalf("Wrong Expression type. Expected FuncLiteral. got=%T", stmt.Expression)
	}

	if len(fnc.Parameters) != 2 {
		t.Fatalf("Wrong number of Parameters. Expected 2 got=%d", len(fnc.Parameters))
	}
	testLiteralExpr(t, fnc.Parameters[0], "f")
	testLiteralExpr(t, fnc.Parameters[1], "b")

	if len(fnc.Body.Statements) != 1 {
		t.Fatalf("Wrong number of Statements in function body. Expected 1 got=%d", len(fnc.Body.Statements))
	}
	body, ok := fnc.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Function Body is not an ExpressionStatement. got=%T", fnc.Body.Statements[0])
	}

	testInfixExpr(t, body.Expression, "f", "*", "b")
}

func TestParsingFuncParameters(t *testing.T) {
	Tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fnc() {};", expectedParams: []string{}},
		{input: "fnc(f) {};", expectedParams: []string{"f"}},
		{input: "fnc(g, f, b) {};", expectedParams: []string{"g", "f", "b"}},
	}

	for _, tcase := range Tests {
		lex := lexer.New(tcase.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fnc := stmt.Expression.(*ast.FuncLiteral)

		if len(fnc.Parameters) != len(tcase.expectedParams) {
			t.Fatalf("wrong number of Functionparameters expected=%d got=%d", len(tcase.expectedParams), len(fnc.Parameters))
		}

		for i, identifier := range tcase.expectedParams {
			testLiteralExpr(t, fnc.Parameters[i], identifier)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 1 + 2, 2 * 3)"
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("number of statements in program body not 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Wrong Expression type. Expected CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("Wrong number of Arguments. Expected 3 got=%d", len(callExp.Arguments))
	}
	testLiteralExpr(t, callExp.Arguments[0], 1)
	testInfixExpr(t, callExp.Arguments[1], 1, "+", 2)
	testInfixExpr(t, callExp.Arguments[2], 2, "*", 3)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	Test := []struct {
		input              string
		expectedIdentifier string
		expectedArguments  []string
	}{
		{
			input:              "print();",
			expectedIdentifier: "print",
			expectedArguments:  []string{},
		},
		{
			input:              "print();",
			expectedIdentifier: "print",
			expectedArguments:  []string{},
		},
		{
			input:              "print();",
			expectedIdentifier: "print",
			expectedArguments:  []string{},
		},
	}
	for _, tcase := range Test {
		lex := lexer.New(tcase.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		callExpr, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt is not CallExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, callExpr.Function, tcase.expectedIdentifier) {
			return
		}

		if len(callExpr.Arguments) != len(tcase.expectedArguments) {
			t.Fatalf("Wrong number of Arguments. Expected=%d, got=%d",
				len(tcase.expectedArguments), len(callExpr.Arguments))
		}

		for i, arg := range tcase.expectedArguments {
			if callExpr.Arguments[i].String() != arg {
				t.Errorf("argument number %d is wrong. Expected=%q, got=%q", i, arg, callExpr.Arguments[i].String())
			}
		}
	}
}

func TestParsingPrefixExpr(t *testing.T) {
	Tests := []struct {
		inp      string
		operator string
		val      interface{}
	}{
		{"!8", "!", 8},
		{"-1907", "-", 1907},
		{"!false", "!", false},
		{"!true", "!", true},
		{"-false", "-", false},
		{"&3", "&", 3},
		{"*32", "*", 32},
	}

	for _, tcase := range Tests {
		lex := lexer.New(tcase.inp)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("Program has more or less than 1 Statement. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not ExpressionStatement. got=%T", program.Statements[0])
		}
		prf, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Statement Expression is not of Type PrefixExpression. got=%T", stmt.Expression)
		}
		if prf.Operator != tcase.operator {
			t.Fatalf("Operator is not %s. got=%s", tcase.operator, prf.Operator)
		}
		if !testLiteralExpr(t, prf.Right, tcase.val) {
			return
		}

	}
}

func TestParsingInfixExpr(t *testing.T) {
	Tests := []struct {
		inp      string
		leftVal  interface{}
		operator string
		rightVal interface{}
	}{
		{"8 + 7;", 8, "+", 7},
		{"8 - 7;", 8, "-", 7},
		{"8 * 7;", 8, "*", 7},
		{"8 / 7;", 8, "/", 7},
		{"8 > 7;", 8, ">", 7},
		{"8 < 7;", 8, "<", 7},
		{"8 == 7;", 8, "==", 7},
		{"8 != 7;", 8, "!=", 7},
		{"true == false", true, "==", false},
		{"true != false", true, "!=", false},
	}

	for _, tcase := range Tests {
		lex := lexer.New(tcase.inp)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("Program has more or less than 1 Statement. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not ExpressionStatement. got=%T", program.Statements[0])
		}
		inf, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Statement Expression is not of Type InfixExpression. got=%T", stmt.Expression)
		}
		if !testLiteralExpr(t, inf.Left, tcase.leftVal) {
			return
		}
		if inf.Operator != tcase.operator {
			t.Fatalf("Operator is not %s. got=%s", tcase.operator, inf.Operator)
		}
		if !testLiteralExpr(t, inf.Right, tcase.rightVal) {
			return
		}

	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := `[1, 3, 5, 7 * 9, 25 - 5]`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Expression is not ArrayLiteral got=%T", stmt.Expression)
	}

	if len(arr.Elements) != 5 {
		t.Fatalf("wron number of elements for arr. Expected 5 got=%d", len(arr.Elements))
	}

	testIntegerLiteral(t, arr.Elements[0], 1)
	testIntegerLiteral(t, arr.Elements[1], 3)
	testIntegerLiteral(t, arr.Elements[2], 5)
	testInfixExpr(t, arr.Elements[3], 7, "*", 9)
	testInfixExpr(t, arr.Elements[4], 25, "-", 5)
}

func TestEmptyArrayLiteralParsing(t *testing.T) {
	input := `[]`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Expression is not ArrayLiteral got=%T", stmt.Expression)
	}

	if len(arr.Elements) != 0 {
		t.Fatalf("wron number of elements for arr. Expected 0 got=%d", len(arr.Elements))
	}

}

func TestIndexExpressionParsing(t *testing.T) {
	input := `arr [8 + 2]`
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	idexp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("Wrong Expressiontype expected Indexexpresion got %T", stmt.Expression)
	}
	if !testIdentifier(t, idexp.Left, "arr") {
		return
	}
	if !testInfixExpr(t, idexp.Index, 8, "+", 2) {
		return
	}

}

func TestEmptyHashLiteralParsing(t *testing.T) {
	input := "{}"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	hsh, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expression is not Hasliteral got=%T", stmt.Expression)
	}
	if len(hsh.Pairs) != 0 {
		t.Errorf("Incorrect length of Hash.pairs expected 0 got=%d", len(hsh.Pairs))
	}
}

func TestHashLiteralParsingStringKeys(t *testing.T) {
	input := `{"age": 23, "year": 2024, "month": 12}`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	hsh, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expression is not Hasliteral got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"age":   23,
		"year":  2024,
		"month": 12,
	}

	if len(hsh.Pairs) != len(expected) {
		t.Errorf("Incorrect length of Hash.pairs expected %d got=%d", len(expected), len(hsh.Pairs))
	}

	for k, v := range hsh.Pairs {
		key, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("Wrong key Type. Expected string got=%T", k)
			continue
		}

		expectedVal := expected[key.String()]
		testIntegerLiteral(t, v, expectedVal)
	}
}

func TestHashLiteralParsingBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 0}`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	hsh, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expression is not Hasliteral got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 0,
	}

	if len(hsh.Pairs) != len(expected) {
		t.Errorf("Incorrect length of Hash.pairs expected %d got=%d", len(expected), len(hsh.Pairs))
	}

	for k, v := range hsh.Pairs {
		key, ok := k.(*ast.Boolean)
		if !ok {
			t.Errorf("Wrong key Type. Expected Boolean got=%T", k)
			continue
		}

		expectedVal := expected[key.String()]
		testIntegerLiteral(t, v, expectedVal)
	}
}

func TestHashLiteralParsingIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	hsh, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expression is not Hasliteral got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hsh.Pairs) != len(expected) {
		t.Errorf("Incorrect length of Hash.pairs expected %d got=%d", len(expected), len(hsh.Pairs))
	}

	for k, v := range hsh.Pairs {
		key, ok := k.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("Wrong key Type. Expected Integer got=%T", k)
			continue
		}

		expectedVal := expected[key.String()]
		testIntegerLiteral(t, v, expectedVal)
	}
}

func TestHashLiteralParsingExpressionValues(t *testing.T) {
	input := `{"first": 1 + 2, "second": 2 * 3, "third": 3 - 8}`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	hsh, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expression is not Hasliteral got=%T", stmt.Expression)
	}

	if len(hsh.Pairs) != 3 {
		t.Errorf("Incorrect length of Hash.pairs expected %d got=%d", 3, len(hsh.Pairs))
	}

	expected := map[string]func(ast.Expression){
		"first": func(exp ast.Expression) {
			testInfixExpr(t, exp, 1, "+", 2)
		},
		"second": func(exp ast.Expression) {
			testInfixExpr(t, exp, 2, "*", 3)
		},
		"third": func(exp ast.Expression) {
			testInfixExpr(t, exp, 3, "-", 8)
		},
	}

	for k, v := range hsh.Pairs {
		key, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("Wrong key Type. Expected String got=%T", k)
			continue
		}

		test, ok := expected[key.String()]
		if !ok {
			t.Errorf("Couldnt find test Function for key %q", key.String())
			continue
		}
		test(v)
	}
}

func TestPrecedenceParsing(t *testing.T) {
	test := []struct {
		inp   string
		expct string
	}{
		{
			"f * -b",
			"(f * (-b))",
		},
		{
			"-!f",
			"(-(!f))",
		},
		{
			"g + f + b",
			"((g + f) + b)",
		},
		{
			"g - f + b",
			"((g - f) + b)",
		},
		{
			"g * f * b",
			"((g * f) * b)",
		},
		{
			"g / f * b",
			"((g / f) * b)",
		},
		{
			"g - f / b",
			"(g - (f / b))",
		},
		{
			"g + f / b + a * e - c",
			"(((g + (f / b)) + (a * e)) - c)",
		},
		{
			"12 + 43; -2 * 64",
			"(12 + 43)((-2) * 64)",
		},
		{
			"2 < 42 == 32 < 4",
			"((2 < 42) == (32 < 4))",
		},
		{
			"532 < 42 != 332 > 41",
			"((532 < 42) != (332 > 41))",
		},
		{
			"17 - 2 * 4 == 6 * 2 + 23 * 5",
			"((17 - (2 * 4)) == ((6 * 2) + (23 * 5)))",
		},

		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"2 < 23 == true",
			"((2 < 23) == true)",
		},
		{
			"6 < 2 == false",
			"((6 < 2) == false)",
		},
		{
			"3 + (2 + 5) + 1",
			"((3 + (2 + 5)) + 1)",
		},
		{
			"3 * (2 + 5) ",
			"(3 * (2 + 5))",
		},
		{
			"(3 + 2) / 5 ",
			"((3 + 2) / 5)",
		},
		{
			"- (3 + 2) ",
			"(-(3 + 2))",
		},
		{
			"!(true != false)",
			"(!(true != false))",
		},
		{
			"g + print(f * b) + e",
			"((g + print((f * b))) + e)",
		},
		{
			"print(g, f, 1, 9 * 0, 7 + 3, show(6, 7 * 8))",
			"print(g, f, 1, (9 * 0), (7 + 3), show(6, (7 * 8)))",
		},
		{
			"print(g + f + b * e / f + g)",
			"print((((g + f) + ((b * e) / f)) + g))",
		},
		{
			"g * [1, 4, 5, 6, 8][3 + 2] * f",
			"((g * ([1, 4, 5, 6, 8][(3 + 2)])) * f)",
		},
		{
			"print(3 * g[5], f[1], 9 * [3, 5] [1])",
			"print((3 * (g[5])), (f[1]), (9 * ([3, 5][1])))",
		},
	}

	for _, tcase := range test {
		lex := lexer.New(tcase.inp)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		str := program.String()
		if str != tcase.expct {
			t.Errorf("expected=%q, received=%q", tcase.expct, str)
		}
	}
}

func testLetStatement(t *testing.T, st ast.Statement, ident string) bool {
	if st.TokenLiteral() != "let" {
		t.Errorf("st.Tokenliteral not 'let'. got=%q", st.TokenLiteral())
		return false
	}
	letStmt, ok := st.(*ast.LetStatement)
	if !ok {
		t.Errorf("st not *ast.LetStatement. got=%T", st)
		return false
	}
	if letStmt.Name.Value != ident {
		t.Errorf("Error Parsing letStmt.Name. expected=%s, got=%s", ident, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != ident {
		t.Errorf("Error Parsing st.Name. expected=%s, got=%s", ident, letStmt.Name)
		return false
	}
	return true
}
func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d Errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, val int64) bool {
	inte, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expression not IntegerLiteral. got=%T", expr)
		return false
	}

	if inte.Value != val {
		t.Errorf("inte.Value is not %d. got=%d", val, inte.Value)
		return false
	}

	if inte.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("inte.TokenLiteral not %d. got=%s", val, inte.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, val string) bool {
	id, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression not Identifier. got=%T", expr)
		return false
	}

	if id.Value != val {
		t.Errorf("id.Value is not %s. got=%s", val, id.Value)
		return false
	}

	if id.TokenLiteral() != val {
		t.Errorf("id.TokenLiteral not %s. got=%s", val, id.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpr(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpr(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpr(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBoolLiteral(t *testing.T, exp ast.Expression, val bool) bool {
	boo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("Expression not BoolLiteral. got=%T", exp)
		return false
	}

	if boo.Value != val {
		t.Errorf("boo.Value is not %t. got=%t", val, boo.Value)
		return false
	}

	if boo.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("boo.TokenLiteral not %t. got=%s", val, boo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpr(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
