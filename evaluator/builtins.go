package evaluator

import (
	"github.com/Muto1907/interpreterInGo/object"
)

var builtIns = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of arguments for `len need=%d got=%d", 1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("invalid argument for `len` got %s", args[0].Type())
			}
		},
	},
}
