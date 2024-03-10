package eval

import (
	"fmt"

	"github.com/samasno/little-compiler/pkg/ast"
	"github.com/samasno/little-compiler/pkg/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.LetStatement:
    println(node.String())
    val := Eval(node.Value, env)
    if isError(val) {
      return val
    }

    _, ok := env.Get(node.Name.Value) 
    if ok {
      return newError("binding for %s already exists", node.Name.Value)
    }
    env.Set(node.Name.Value, val)
    return val
       
	case *ast.Identifier:
    val, ok := env.Get(node.Value)
    if !ok {
      return newError("variable %s does not exist", node.Value)
    }

    return val

	case *ast.ReturnStatement:
    val := Eval(node.Value, env)
    if isError(val) {
      return val
    }

		return &object.Return{Value: val}

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BoolLiteral:
		return returnNativeBool(node.Value)

	case *ast.FnLiteral:
    params := node.Params
    body := node.Body
    return &object.Function{Params: params, Body: body, Env: env}

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
    if isError(left) {
      return left
    }
    
		right := Eval(node.Right, env)
    if isError(right) {
      return right
    }

		return evalInfix(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
      return right
    }

    return evalPrefix(node.Operator, right)

  case *ast.CallExpression:
    fn := Eval(node.Function, env)

    if isError(fn) {
      return fn
    }

    args := evalExpressions(node.Arguments, env)
    
    if len(args) == 1 && isError(args[0]) {
      return args[0]
    }
    
    return applyFunction(fn, args)
    
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	}

	return NULL
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
  args := []object.Object{}

  for _, exp := range exps {
    evaluated := Eval(exp, env)
    if isError(evaluated) {
      return []object.Object{evaluated}
    }

    args = append(args, evaluated)
  }

  return args
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
  function, ok := fn.(*object.Function)

  if !ok {
    return newError("expected function, got %s", fn.Type())
  }

  extendedEnv := extendedFunctionEnv(function, args)

  result := Eval(function.Body, extendedEnv)

  return unwrapEvaluatedReturn(result)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
  env := object.NewEnclosedEnvironment(fn.Env)

  for pid, p := range fn.Params {
    env.Set(p.Value, args[pid])
  }

  return env
}

func unwrapEvaluatedReturn(obj object.Object) object.Object {
  if rv, ok := obj.(*object.Return); ok {
    return rv.Value
  }
  return obj
}

func isError(obj object.Object) bool {
  if obj.Type() == object.ERROR_OBJ {
    return true
  }

  return false
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
		
    switch r := result.(type) {
      case *object.Return:
        return r.Value
      
      case *object.Error:
        return r
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
  case left.Type() != right.Type():
    return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
    return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	con := Eval(node.Condition, env)
	res := isTruthy(con)
	if res && node.Consequence != nil {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
    switch result.Type() {
    case object.RETURN_OBJ, object.ERROR_OBJ:
      return result
    }
	}

	return result
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
    return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
  default:
    return newError("unknown operator: %s %s", operator, right.Type())
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

func newError(format string, a ...interface{}) *object.Error {
  return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalMinusOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
    return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}


var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)
