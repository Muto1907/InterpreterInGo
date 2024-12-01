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
	"head": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of arguments for `head` need=%d got=%d", 1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return newError("invalid argument for `head` expected ARRAY got %s", arg.Type())
			}
		},
	},
	"last": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of arguments for `last` need=%d got=%d", 1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL
			default:
				return newError("invalid argument for `last` expected ARRAY got %s", arg.Type())
			}
		},
	},
	"tail": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of arguments for `tail` need=%d got=%d", 1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					return &object.Array{Elements: arg.Elements[1:]}
				}
				return NULL
			default:
				return newError("invalid argument for `tail` expected ARRAY got %s", arg.Type())
			}
		},
	},
	"push": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("invalid number of arguments for `push` need=%d got=%d", 2, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Array{Elements: append(arg.Elements, args[1])}
			default:
				return newError("invalid argument for `push` expected ARRAY got %s", arg.Type())
			}
		},
	},
}
