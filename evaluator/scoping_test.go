package evaluator

import (
	"testing"

	"github.com/Muto1907/interpreterInGo/object"
)

func TestScoping_BlockScope(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int64
	}{{
		input: `
			let x = 5;
			if(true) {
				let y = 99;
				x = x + y;
			};
			x;
		
		`,
		expectedVal: 104,
	},
		{
			input: `
				let x = 2 * 2;
				if (true) {
					let y = x + 8;
				};
				x;
		`,
			expectedVal: 4,
		},
	}
	for _, tcase := range tests {
		val := testEval(tcase.input)
		errObj, ok := val.(*object.Error)
		if ok {
			t.Fatalf("Got Error: %s", errObj.Message)
		}
		testIntegerObject(t, val, tcase.expectedVal)
	}
}

func TestScoping_Shadowing(t *testing.T) {
	input := `
		let x = 10;
		if (true) {
		let x = 99;
		}
		x;
	`
	val := testEval(input)
	testIntegerObject(t, val, 10)
}

func TestScoping_FunctionClosure(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int64
	}{
		{
			input: `
				let outerVal = 50;
				let makeAdder = fnc(){
					return fnc(x) { x + outerVal };
				};
				let addOuter = makeAdder();
				addOuter(10);
			`,
			expectedVal: 60,
		},
		{
			input: `
				let x = 5
				let fn = fnc(){
					let x = 999;
					return x;
				};
				fn()
			`,
			expectedVal: 999,
		},
		{
			input: `
				let x = 10;
				let fn = fnc(){
					x = x + 1;
				};
				fn();
				x;	
			`,
			expectedVal: 11,
		},
	}
	for _, tcase := range tests {
		val := testEval(tcase.input)
		testIntegerObject(t, val, tcase.expectedVal)
	}
}

func TestScoping_BlockInFunction(t *testing.T) {
	input := `
		let fn = fnc(){
			let a = 10;
			if (true) {
				let b = a + 5;
				return b;
				a = 999;
			}
		};
		fn();
	`
	val := testEval(input)
	testIntegerObject(t, val, 15)
}

/*
func TestScoping_WhileLoopBlocks(t *testing.T) {
	input := `
		let i = 0;
		while (i < 5) {
			if (true){
				let temp = i;
				i = temp + 2;
			}
		}
		i;
	`
	val := testEval(input)
	testIntegerObject(t, val, 6)
}
*/
