package lib

import (
	"1ylang/object"
	"strings"
)

var arrayFuncs = map[string]interface{}{
	"len": func(arr []interface{}) float64 {
		return float64(len(arr))
	},
	"push": func(arr []interface{}, elem interface{}) []interface{} {
		return append(arr, elem)
	},
	"pop": func(arr []interface{}) ([]interface{}, interface{}) {
		if len(arr) == 0 {
			return arr, nil
		}
		elem := arr[len(arr)-1]
		return arr[:len(arr)-1], elem
	},
	"shift": func(arr []interface{}) ([]interface{}, interface{}) {
		if len(arr) == 0 {
			return arr, nil
		}
		elem := arr[0]
		return arr[1:], elem
	},
	"unshift": func(arr []interface{}, elem interface{}) []interface{} {
		return append([]interface{}{elem}, arr...)
	},
	"indexOf": func(arr []interface{}, elem interface{}) float64 {
		for i, v := range arr {
			if object.IsEqual(v.(object.Object), elem.(object.Object)) {
				return float64(i)
			}
		}
		return -1
	},
	"contains": func(arr []interface{}, elem interface{}) bool {
		for _, v := range arr {
			if object.IsEqual(v.(object.Object), elem.(object.Object)) {
				return true
			}
		}
		return false
	},
	"slice": func(arr []interface{}, start, end float64) []interface{} {
		return arr[int(start):int(end)]
	},
	"join": func(arr []interface{}, sep string) string {
		elements := make([]string, len(arr))
		for i, v := range arr {
			elements[i] = v.(object.Object).Inspect()
		}
		return strings.Join(elements, sep)
	},
}

func RegisterArrayFuncs(env *object.Environment) {
	object.RegisterFunctions(env, "Array", arrayFuncs)
}
