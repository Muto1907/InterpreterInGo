package evaluator

import (
	"testing"

	"github.com/Muto1907/interpreterInGo/lexer"
	"github.com/Muto1907/interpreterInGo/object"
	"github.com/Muto1907/interpreterInGo/parser"
)

func TestEvalIntExpr(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int64
	}{
		{
			input: "6", expectedVal: 6,
		},
		{
			input: "1907", expectedVal: 1907,
		},
	}

	for _, tcase := range tests {
		val := testEval(tcase.input)
		testIntegerObject(t, val, tcase.expectedVal)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, integerObject object.Object, expectedInt int64) bool {
	res, ok := integerObject.(*object.Integer)
	if !ok {
		t.Fatalf("Object is not integerObject. got=%T", integerObject)
	}
	if res.Value != expectedInt {
		t.Errorf("Unexpected Value of IntegerObject. Expected=%d, got=%d", expectedInt, res.Value)
		return false
	}
	return true
}