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
		{"6", 6},
		{"1907", 1907},
		{"-3", -3},
		{"--3", 3},
		{"3 + 3 + 3 + 3- 6", 6},
		{"3 * 2 * 4 * 3 * 2", 144},
		{"-30 + 60 +-30", 0},
		{"4 * 4 + 10", 26},
		{"2 + 6 * 10", 62},
		{"30 + 2 *-15", 0},
		{"30 / 2 * 2 + 10", 40},
		{"3 * (8 + 12)", 60},
		{"3 * 3 * 3 + 12", 39},
		{"3 * (3 * 3) + 12", 39},
		{"(3 + 8 * 2 + 15 / 3) * 2 +-10", 38},
	}

	for _, tcase := range tests {
		val := testEval(tcase.input)
		testIntegerObject(t, val, tcase.expectedVal)
	}
}

func TestEvalBoolExpr(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal bool
	}{
		{
			"true", true,
		},
		{
			"false", false,
		},
	}

	for _, tcase := range tests {
		val := testEval(tcase.input)
		testBooleanObject(t, val, tcase.expectedVal)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!2", false},
		{"!!2", true},
	}
	for _, tcase := range tests {
		val := testEval(tcase.input)
		testBooleanObject(t, val, tcase.expectedVal)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()

	return Eval(program)
}

func testBooleanObject(t *testing.T, booleanObject object.Object, expectedbool bool) bool {
	res, ok := booleanObject.(*object.Boolean)
	if !ok {
		t.Fatalf("Object is not booleanObject. got=%T", booleanObject)
	}
	if res.Value != expectedbool {
		t.Errorf("Unexpected Value of booleanObject. Expected=%t, got=%t", expectedbool, res.Value)
		return false
	}
	return true
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
