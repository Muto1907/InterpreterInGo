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
		expectedIdent string
	}{
		{"i"},
		{"j"},
		{"testval"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdent) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 0;
	return 2;
	return 1907;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
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

	fnc, ok := stmt.Expression.(*ast.FuncExpression)
	if !ok {
		t.Fatalf("Wrong Expression type. Expected FuncExpression. got=%T", stmt.Expression)
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
