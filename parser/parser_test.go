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

func TestParsingPrefixExpr(t *testing.T) {
	Tests := []struct {
		inp      string
		operator string
		intVal   int64
	}{
		{"!8", "!", 8},
		{"-1907", "-", 1907},
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
		if !testIntegerLiteral(t, prf.Right, tcase.intVal) {
			return
		}

	}
}

func TestParsingInfixExpr(t *testing.T) {
	Tests := []struct {
		inp      string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"8 + 7;", 8, "+", 7},
		{"8 - 7;", 8, "-", 7},
		{"8 * 7;", 8, "*", 7},
		{"8 / 7;", 8, "/", 7},
		{"8 > 7;", 8, ">", 7},
		{"8 < 7;", 8, "<", 7},
		{"8 == 7;", 8, "==", 7},
		{"8 != 7;", 8, "!=", 7},
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
		if !testIntegerLiteral(t, inf.Left, tcase.leftVal) {
			return
		}
		if inf.Operator != tcase.operator {
			t.Fatalf("Operator is not %s. got=%s", tcase.operator, inf.Operator)
		}
		if !testIntegerLiteral(t, inf.Right, tcase.rightVal) {
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

	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
