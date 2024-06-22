package evaluator

import (
	"1ylang/ast"
	"1ylang/lexer"
	"1ylang/object"
	"1ylang/parser"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// Eval evaluates an AST node
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		val := Eval(node.Expression, env)
		return val

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpression(node.Operator, left)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.NewVar(node.Name.Value, val)

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.NewConst(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.Assignment:
		return evalAssignmentExpression(node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ImportExpression:
		return evalImportExpression(node, env)

	case *ast.DotExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right, ok := node.Right.(*ast.Identifier)
		if !ok {
			return newError("expected property name to be identifier, got %T", node.Right)
		}

		return evalDotExpression(left, right)

	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CNT
	}

	return nil
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	BREAK = &object.Break{}
	CNT   = &object.Continue{}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "~":
		return evalTildePrefixOperatorExpression(right)
	case "++":
		return evalIncrementExpression(right, true)
	case "--":
		return evalDecrementExpression(right, true)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		if integerObj, ok := right.(*object.Integer); ok {
			if integerObj.Value.Cmp(big.NewInt(0)) == 0 {
				return TRUE
			}
			return FALSE
		}
		if floatObj, ok := right.(*object.Float); ok {
			if floatObj.Value.Cmp(big.NewFloat(0.0)) == 0 {
				return TRUE
			}
			return FALSE
		}
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: new(big.Int).Neg(right.Value)}
	case *object.Float:
		return &object.Float{Value: new(big.Float).Neg(right.Value)}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "&&":
		return evalLogicalAndExpression(left, right)
	case operator == "||":
		return evalLogicalOrExpression(left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.HASH_OBJ && right.Type() == object.HASH_OBJ:
		return evalHashInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ || left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ:
		if operator == "*" {
			return &object.String{Value: strings.Repeat(left.(*object.String).Value, int(right.(*object.Integer).Value.Int64()))}
		}
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right)
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.INTEGER_OBJ || left.Type() == object.INTEGER_OBJ && right.Type() == object.ARRAY_OBJ:
		if operator == "*" {
			elements := make([]object.Object, 0)
			for i := 0; i < int(right.(*object.Integer).Value.Int64()); i++ {
				elements = append(elements, left.(*object.Array).Elements...)
			}
			return &object.Array{Elements: elements}
		}
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}



func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: new(big.Int).Add(leftVal, rightVal)}
	case "+=":
		result := new(big.Int).Add(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "-":
		return &object.Integer{Value: new(big.Int).Sub(leftVal, rightVal)}
	case "-=":
		result := new(big.Int).Sub(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "*":
		return &object.Integer{Value: new(big.Int).Mul(leftVal, rightVal)}
	case "*=":
		result := new(big.Int).Mul(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "/":
		if rightVal.Cmp(big.NewInt(0)) == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: new(big.Int).Div(leftVal, rightVal)}
	case "/=":
		if rightVal.Cmp(big.NewInt(0)) == 0 {
			return newError("division by zero")
		}
		result := new(big.Int).Div(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "%":
		if rightVal.Cmp(big.NewInt(0)) == 0 {
			return newError("modulus by zero")
		}
		return &object.Integer{Value: new(big.Int).Mod(leftVal, rightVal)}
	case "%=":
		if rightVal.Cmp(big.NewInt(0)) == 0 {
			return newError("modulus by zero")
		}
		result := new(big.Int).Mod(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "**":
		return &object.Float{Value: bigFloatPow(new(big.Float).SetInt(leftVal), new(big.Float).SetInt(rightVal))}
	case "**=":
		result := bigFloatPow(new(big.Float).SetInt(leftVal), new(big.Float).SetInt(rightVal))
		left = &object.Float{Value: result}
		return &object.Float{Value: result}
	case "<":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) < 0)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) > 0)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) == 0)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) != 0)
	case ">=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) >= 0)
	case "<=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) <= 0)
	case "&":
		return &object.Integer{Value: new(big.Int).And(leftVal, rightVal)}
	case "&=":
		result := new(big.Int).And(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "|":
		return &object.Integer{Value: new(big.Int).Or(leftVal, rightVal)}
	case "|=":
		result := new(big.Int).Or(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "^":
		return &object.Integer{Value: new(big.Int).Xor(leftVal, rightVal)}
	case "^=":
		result := new(big.Int).Xor(leftVal, rightVal)
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case ">>":
		return &object.Integer{Value: new(big.Int).Rsh(leftVal, uint(rightVal.Int64()))}
	case ">>=":
		result := new(big.Int).Rsh(leftVal, uint(rightVal.Int64()))
		leftVal.Set(result)
		return &object.Integer{Value: result}
	case "<<":
		return &object.Integer{Value: new(big.Int).Lsh(leftVal, uint(rightVal.Int64()))}
	case "<<=":
		result := new(big.Int).Lsh(leftVal, uint(rightVal.Int64()))
		leftVal.Set(result)
		return &object.Integer{Value: result}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func bigFloatPow(x, y *big.Float) *big.Float {
	xVal, _ := x.Float64()
	yVal, _ := y.Float64()
	return new(big.Float).SetFloat64(math.Pow(xVal, yVal))
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := toFloat(left)
	rightVal := toFloat(right)

	result := new(big.Float)

	switch operator {
	case "+":
		result.Add(leftVal, rightVal)
	case "-":
		result.Sub(leftVal, rightVal)
	case "*":
		result.Mul(leftVal, rightVal)
	case "/":
		if rightVal.Cmp(big.NewFloat(0)) == 0 {
			return newError("division by zero")
		}
		result.Quo(leftVal, rightVal)
	case "**":
		result = bigFloatPow(leftVal, rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) < 0)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) > 0)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) == 0)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) != 0)
	case ">=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) >= 0)
	case "<=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) <= 0)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	return &object.Float{Value: result}
}

func toFloat(obj object.Object) *big.Float {
	switch obj := obj.(type) {
	case *object.Integer:
		return new(big.Float).SetInt(obj.Value)
	case *object.Float:
		return obj.Value
	default:
		return new(big.Float)
	}
}



func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	for _, elif := range ie.Elifs {
		condition := Eval(elif.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(elif.Consequence, env)
		}
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok, _ := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := big.NewInt(int64(len(arrayObj.Elements) - 1))

	if idx.Cmp(big.NewInt(0)) < 0 || idx.Cmp(max) > 0 {
		return NULL
	}

	return arrayObj.Elements[idx.Int64()]
}

func evalAssignmentExpression(node *ast.Assignment, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	switch name := node.Name.(type) {
	case *ast.Identifier:
		_, ok, readOnly := env.Get(name.Value)
		if !ok {
			return newError("identifier not found: " + name.Value)
		}

		if readOnly {
			return newError("cannot assign to constant '%s'", name.Value)
		}

		env.Set(name.Value, val)
		return val

	case *ast.DotExpression:
		left := Eval(name.Left, env)
		if isError(left) {
			return left
		}

		right, ok := name.Right.(*ast.Identifier)
		if !ok {
			return newError("expected property name to be identifier, got %T", name.Right)
		}

		return evalDotAssignment(left, right, val)

	default:
		return newError("invalid assignment target: %T", node.Name)
	}
}

func evalDotAssignment(left object.Object, right *ast.Identifier, val object.Object) object.Object {
	switch left := left.(type) {
	case *object.Hash:
		key := &object.String{Value: right.Value}
		hashKey := key.HashKey()
		left.Pairs[hashKey] = object.HashPair{Key: key, Value: val}
		return val
	default:
		return newError("not a hash: %s", left.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			return NULL
		}

		whileEnv := object.NewEnclosedEnvironment(env)
		result := evalLoopStatement(ws.Body, whileEnv)
		if result != nil {
			if result.Type() == object.BREAK_OBJ {
				return NULL
			}
			if result.Type() == object.CONTINUE_OBJ {
				continue
			}
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}
}

func evalLoopStatement(body *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range body.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func evalTildePrefixOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: new(big.Int).Not(right.Value)}
	default:
		return newError("unknown operator: ~%s", right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalHashInfixExpression(operator string, left, right object.Object) object.Object {
	leftHash := left.(*object.Hash)
	rightHash := right.(*object.Hash)

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(object.IsEqual(leftHash, rightHash))
	case "!=":
		return nativeBoolToBooleanObject(!object.IsEqual(leftHash, rightHash))
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalArrayInfixExpression(operator string, left, right object.Object) object.Object {
	leftArr := left.(*object.Array)
	rightArr := right.(*object.Array)

	switch operator {
	case "+":
		elements := append(leftArr.Elements, rightArr.Elements...)
		return &object.Array{Elements: elements}
	case "==":
		return nativeBoolToBooleanObject(object.IsEqual(leftArr, rightArr))
	case "!=":
		return nativeBoolToBooleanObject(!object.IsEqual(leftArr, rightArr))
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(operator string, left object.Object) object.Object {
	switch operator {
	case "++":
		return evalIncrementExpression(left, false)
	case "--":
		return evalDecrementExpression(left, false)
	default:
		return newError("unknown operator: %s%s", operator, left.Type())
	}
}

func evalIncrementExpression(operand object.Object, isPrefix bool) object.Object {
	switch operand := operand.(type) {
	case *object.Integer:
		if isPrefix {
			operand.Value.Add(operand.Value, big.NewInt(1))
			return operand
		} else {
			result := &object.Integer{Value: new(big.Int).Set(operand.Value)}
			operand.Value.Add(operand.Value, big.NewInt(1))
			return result
		}
	default:
		return newError("unknown operator: ++%s", operand.Type())
	}
}

func evalDecrementExpression(operand object.Object, isPrefix bool) object.Object {
	switch operand := operand.(type) {
	case *object.Integer:
		if isPrefix {
			operand.Value.Sub(operand.Value, big.NewInt(1))
			return operand
		} else {
			result := &object.Integer{Value: new(big.Int).Set(operand.Value)}
			operand.Value.Sub(operand.Value, big.NewInt(1))
			return result
		}
	default:
		return newError("unknown operator: --%s", operand.Type())
	}
}

func evalLogicalAndExpression(left, right object.Object) object.Object {
	if isTruthy(left) {
		return right
	}
	return left
}

func evalLogicalOrExpression(left, right object.Object) object.Object {
	if isTruthy(left) {
		return left
	}
	return right
}

func evalDotExpression(left object.Object, right *ast.Identifier) object.Object {
	switch left := left.(type) {
	case *object.Hash:
		key := &object.String{Value: right.Value}
		hashKey := key.HashKey()
		if pair, ok := left.Pairs[hashKey]; ok {
			return pair.Value
		} else {
			return NULL
		}
	default:
		return newError("not a hash: %s", left.Type())
	}
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	strObj := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := big.NewInt(int64(len(strObj.Value) - 1))

	if idx.Cmp(big.NewInt(0)) < 0 || idx.Cmp(max) > 0 {
		return NULL
	}

	return &object.String{Value: string(strObj.Value[idx.Int64()])}
}

func evalImportExpression(ie *ast.ImportExpression, env *object.Environment) object.Object {
	// Evaluate the import path
	pathObj := Eval(ie.Path, env)
	if pathObj.Type() != object.STRING_OBJ {
		return newError("import path must be a string, got %s", pathObj.Type())
	}

	path := pathObj.(*object.String).Value
	if !strings.HasSuffix(path, ".1y") {
		path += ".1y"
	}

	// Try to read the file from the current working directory or interpreter directory
	content, err := readFileFromCurrentOrInterpreterDir(path)
	if err != nil {
		return newError("could not read file: %s", path)
	}

	// Lexical and syntactical analysis
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return newError("parsing file %s failed: %s", path, strings.Join(p.Errors(), "\n"))
	}

	// Create a new environment and execute the program
	newEnv := object.NewEnvironment()
	Eval(program, newEnv)

	// Wrap the variables in the new environment into a Hash object
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	for k, v := range newEnv.Store() {
		hashKey := &object.String{Value: k}
		hash.Pairs[hashKey.HashKey()] = object.HashPair{Key: hashKey, Value: v.Value}
	}

	return hash
}

func readFileFromCurrentOrInterpreterDir(path string) ([]byte, error) {
	// Try to read the file from the current working directory
	content, err := os.ReadFile(path)
	if err == nil {
		return content, nil
	}

	// Get the interpreter directory
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	// Try to read the file from the interpreter directory
	fullPath := filepath.Join(execDir, path)
	return os.ReadFile(fullPath)
}
