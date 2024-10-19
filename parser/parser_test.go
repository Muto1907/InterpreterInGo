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
