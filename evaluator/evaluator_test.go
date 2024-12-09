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

func TestEvalStringExpr(t *testing.T) {
	input := `"whats up"`
	val := testEval(input)
	str, ok := val.(*object.String)
	if !ok {
		t.Fatalf("Wrong Object Type. Expected String got=%T(%v)", val, val)
	}

	if str.Value != "whats up" {
		t.Fatalf("Wrong StringValue. Expected whats up got=%q", str.Value)
	}
}

func TestEvalConcatenation(t *testing.T) {
	input := `"Hey" + " " + "what's" + " " + "up";`
	val := testEval(input)
	str, ok := val.(*object.String)
	if !ok {
		t.Fatalf("Wrong Object Type. Expected String got=%T(%v)", val, val)
	}
	if str.Value != "Hey what's up" {
		t.Fatalf("Wrong String expected=%s got=%s", "Hey what's up", str.Value)
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

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int64
	}{
		{`let y = 3; while(y < 4){ let y = y + 1; } return y`, 4},
		{`let y = 3; while(y > 4){ let y = y + 1; } return y`, 3},
	}

	for _, tcase := range tests {
		val := testEval(tcase.input)
		testIntegerObject(t, val, tcase.expectedVal)
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
		{`"Hi " - "you`, "unknown operator: STRING - STRING"},
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

func TestFunctionObjects(t *testing.T) {
	input := "fnc(x) { x * 3; };"
	val := testEval(input)
	fnc, ok := val.(*object.Function)
	if !ok {
		t.Fatalf("Unexpected Object Type. Expected Function got=%T(%v)", val, val)
	}
	if len(fnc.Params) != 1 {
		t.Fatalf("Incorrect number of Parameters. Expected 1 got = %d", len(fnc.Params))
	}
	if fnc.Params[0].String() != "x" {
		t.Fatalf("incorrect Parameter Identifier. Expected = x got = %s", fnc.Params[0].String())
	}
	expectedBody := "(x * 3)"
	if fnc.Body.String() != expectedBody {
		t.Fatalf("Incorrect body. Expected= %s. got= %s", expectedBody, fnc.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	test := []struct {
		input       string
		expectedVal int64
	}{
		{"let doNothing = fnc (x){ x; }; doNothing(3)", 3},
		{"let doNothing = fnc (x){ return x; }; doNothing(3)", 3},
		{"let successor = fnc (x){ x + 1; }; successor(3)", 4},
		{"let mult = fnc (x, y){ x * y; }; mult(3, 4)", 12},
		{"let mult = fnc (x, y){ x * y; }; mult(2 * 2, mult(3, 4))", 48},
		{"fnc (x){ x; } (3)", 3},
	}
	for _, tcase := range test {
		testIntegerObject(t, testEval(tcase.input), tcase.expectedVal)
	}
}

func TestBuiltInFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hi");`, 2},
		{`len("");`, 0},
		{`len("Hi what's up");`, 12},
		{`len(4);`, "invalid argument for `len` got INTEGER"},
		{`len("Hey", "Ho")`, "invalid number of arguments for `len need=1 got=2"},
		{`len([3, 6, 9])`, 3},
		{`len([])`, 0},
		{`head([14,12,32])`, 14},
		{`head([])`, nil},
		{`head(1)`, "invalid argument for `head` expected ARRAY got INTEGER"},
		{`last([14,12,32])`, 32},
		{`last([])`, nil},
		{`tail(1)`, "invalid argument for `tail` expected ARRAY got INTEGER"},
		{`tail([14,12,32])`, []int{12, 32}},
		{`tail([])`, nil},
		{`push([], 3)`, []int{3}},
		{`push(3, 3)`, "invalid argument for `push` expected ARRAY got INTEGER"},
	}
	for _, tcase := range tests {
		val := testEval(tcase.input)

		switch expect := tcase.expected.(type) {
		case int:
			testIntegerObject(t, val, int64(expect))
		case string:
			err, ok := val.(*object.Error)
			if !ok {
				t.Errorf("Object is not Error got %T(%v)", val, val)
				continue
			}
			if err.Message != expect {
				t.Errorf("Unexpected ErrorMessage expected=%s got=%s", expect, err.Message)
			}
		}
	}
}

func TestArrayEval(t *testing.T) {
	input := "[1, 3 + 6, 0 * 5, 7 - 0]"
	val := testEval(input)
	arr, ok := val.(*object.Array)
	if !ok {
		t.Errorf("Object is not Array got=%T(%v)", val, val)
	}
	if len(arr.Elements) != 4 {
		t.Errorf("Array Length not equal 4. got = %d", len(arr.Elements))
	}
	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 9)
	testIntegerObject(t, arr.Elements[2], 0)
	testIntegerObject(t, arr.Elements[3], 7)

}

func TestArrayIndexExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[5, 9, 4][0]",
			5,
		},
		{
			"[4, 3, 8][1]",
			3,
		},
		{
			"[13, 25, 90][2]",
			90,
		},
		{
			"let m = 0; [3][m];",
			3,
		},
		{
			"[1, 5, 8][1 + 1];",
			8,
		},
		{
			"let arr = [1, 4, 8]; arr[2];",
			8,
		},
		{
			"let arr = [12, 25, 31]; arr[0] + arr[1] + arr[2];",
			68,
		},
		{
			"let arr = [12, 2, 13]; let i = arr[1]; arr[i]",
			13,
		},
		{
			"[11, 23, 53][3]",
			nil,
		},
		{
			"[11, 23, 53][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
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
