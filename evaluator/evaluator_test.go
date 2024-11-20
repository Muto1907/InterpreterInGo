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
		{"true", true},
		{"false", false},
		{"2 > 3", false},
		{"3 < 5", true},
		{"2 < 2", false},
		{"3 > 3", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"false != true", true},
		{"true != false", true},
		{"(2 < 3) == true", true},
		{"(2 < 3) == false", false},
		{"(3 > 4) == true", false},
		{"(3 > 4) == false", true},
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int64
	}{
		{"return 3", 3},
		{"return 3; 4", 3},
		{"return 3 + 8; 25", 11},
		{"89; return 3 + 8; 25", 11},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 7;
			}`, 10},
	}
	for _, tcase := range tests {
		val := testEval(tcase.input)
		testIntegerObject(t, val, tcase.expectedVal)
	}
}
func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal interface{}
	}{
		{"if (true) { 12 }", 12},
		{"if (false) { 12 }", nil},
		{"if (6) { 12 }", 12},
		{"if (3 < 4) { 12 }", 12},
		{"if (3 > 4) { 12 }", nil},
		{"if (3 > 4) { 12 } else { 22 }", 22},
		{"if (3 < 4) { 12 } else { 22 }", 12},
	}

	for _, tcase := range tests {
		val := testEval(tcase.input)
		integer, ok := tcase.expectedVal.(int)
		if ok {
			testIntegerObject(t, val, int64(integer))
		} else {
			testNullObject(t, val)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		Input            string
		ExpectedErrorMsg string
	}{
		{"3 + false", "type mismatch: INTEGER + BOOLEAN"},
		{"false + true", "unknown operator: BOOLEAN + BOOLEAN"},
		{"3 + true; 3", "type mismatch: INTEGER + BOOLEAN"},
		{"-false", "unknown operator: -BOOLEAN"},
		{"false + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"3; false + false; 3;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (2 > 1) {true + false;}", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
			if (2 < 4) {
				if (2 < 4){
					return false * false;
				}
				return 3;
			}
		`, "unknown operator: BOOLEAN * BOOLEAN"},
		{"stuff", "identifier not found: stuff"},
	}

	for _, tcase := range tests {
		val := testEval(tcase.Input)
		errorObj, ok := val.(*object.Error)
		if !ok {
			t.Errorf("Object is not of Type Error. got=%T(%+v)", val, val)
			continue
		}

		if errorObj.Message != tcase.ExpectedErrorMsg {
			t.Errorf("Wrong Error message. Expected=%q, got=%q", tcase.ExpectedErrorMsg, errorObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		Input       string
		ExpectedVal int64
	}{
		{"let x = 5; x", 5},
		{"let x = 5 * 5; x", 25},
		{"let x = 5; let y = x; y", 5},
		{"let x = 5; let y = x; let z = x + y + 4; z", 14},
	}

	for _, tcase := range tests {
		val := testEval(tcase.Input)
		testIntegerObject(t, val, tcase.ExpectedVal)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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
		t.Errorf("Object is not integerObject. got=%T", integerObject)
		return false
	}
	if res.Value != expectedInt {
		t.Errorf("Unexpected Value of IntegerObject. Expected=%d, got=%d", expectedInt, res.Value)
		return false
	}
	return true
}

func testNullObject(t *testing.T, nullObject object.Object) bool {
	if nullObject != NULL {
		t.Errorf("Nullobject is not NULL. got=%T (%v)", nullObject, nullObject)
		return false
	}
	return true
}
