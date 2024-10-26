package parser

import (
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
		t.Fatalf("Statement is not ExpreseionStatement. got=%T", program.Statements[0])
	}
	id, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Statement Expression is not of Type Identifiert. got=%T", stmt.Expression)
	}
	if id.Value != "thingy" {
		t.Fatalf("Identifier Value is not thingy. got=%s", id.Value)
	}
	if id.TokenLiteral() != "thingy" {
		t.Fatalf("Token Literal is not thingy. got=%s", id.TokenLiteral())
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
