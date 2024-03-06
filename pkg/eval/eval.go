package eval

import (
	"github.com/samasno/little-compiler/pkg/ast"
	"github.com/samasno/little-compiler/pkg/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.LetStatement:
		println("got let statement")

	case *ast.Identifier:
		println("got identifier")

	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.Value)}

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BoolLiteral:
		return returnNativeBool(node.Value)

	case *ast.FnLiteral:
		println("fnlit")

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfix(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefix(node.Operator, right)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
loop:
	for _, stmt := range stmts {
		result = Eval(stmt)
		if result.Type() == object.RETURN_OBJ {
			break loop
		}
	}
	return result
}

func returnNativeBool(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalInfix(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegers(operator, left, right)
	default:
		return NULL
	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	con := Eval(node.Condition)
	res := isTruthy(con)
	if res && node.Consequence != nil {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	} else {
		return NULL
	}
}

func evalInfixIntegers(operator string, left, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case "==":
		return returnNativeBool(l == r)
	case "<=":
		return returnNativeBool(l <= r)
	case "<":
		return returnNativeBool(l < r)
	case ">":
		return returnNativeBool(l > r)
	case "!=":
		return returnNativeBool(l != r)
	default:
		return NULL
	}
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusOperator(right)
	case "!":
		return evalBangOperator(right)
	case "--":
		println("got dec")
	case "++":
		println("got inc")
	}
	return right
}

func isTruthy(condition object.Object) bool {
	switch o := condition.(type) {
	case *object.Boolean:
		return o.Value
	case *object.Integer:
		if o.Value == 0 || o.Value == -0 {
			return false
		}
		return true
	default:
		return false
	}
}

func evalBangOperator(obj object.Object) object.Object {
	switch o := obj.(type) {
	case *object.Boolean:
		if o == TRUE {
			return FALSE
		}
		return TRUE
	case *object.Integer:
		if o.Value == 0 || o.Value == -0 {
			return TRUE
		}
		return FALSE
	}

	return FALSE
}

func evalMinusOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)
