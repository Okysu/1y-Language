package lib

import (
	"1ylang/object"
	"math"
)

var mathFuncs = map[string]interface{}{
	"sin":  math.Sin,
	"cos":  math.Cos,
	"tan":  math.Tan,
	"asin": math.Asin,
	"acos": math.Acos,
	"atan": math.Atan,
	"exp":  math.Exp,
	"log":  math.Log,
	"sqrt": math.Sqrt,
	"pow":  math.Pow,
	"abs":  math.Abs,
	"ceil": math.Ceil,
	"floor": math.Floor,
	"round": math.Round,
	"trunc": math.Trunc,
	"mod":  math.Mod,
	"max":  math.Max,
	"min":  math.Min,
	"hypot": math.Hypot,
	"copysign": math.Copysign,
	"dim":  math.Dim,
	"remainder": math.Remainder,
	"gamma": math.Gamma,
}

func RegisterMathFuncs(env *object.Environment) {
	object.RegisterFunctions(env, "Math", mathFuncs)
}