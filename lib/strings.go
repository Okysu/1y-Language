package lib

import (
	"1ylang/object"
	"strings"
	"unicode"
)

var stringFuncs = map[string]interface{}{
	"concat": func(a, b string) string {
		return a + b
	},
	"len": func(s string) float64 {
		return float64(len(s))
	},
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	"trim": func(s string) string {
		return strings.TrimSpace(s)
	},
	"contains": func(s, substr string) bool {
		return strings.Contains(s, substr)
	},
	"replace": func(s, old, new string, n float64) string {
		return strings.Replace(s, old, new, int(n))
	},
	"split": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"join": func(elems []string, sep string) string {
		return strings.Join(elems, sep)
	},
	"index": func(s, substr string) float64 {
		return float64(strings.Index(s, substr))
	},
	"lastIndex": func(s, substr string) float64 {
		return float64(strings.LastIndex(s, substr))
	},
	"hasPrefix": func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	},
	"hasSuffix": func(s, suffix string) bool {
		return strings.HasSuffix(s, suffix)
	},
	"repeat": func(s string, count float64) string {
		return strings.Repeat(s, int(count))
	},
	"toTitle": func(s string) string {
		return strings.ToTitle(s)
	},
	"toTitleSpecial": func(s string, c *unicode.SpecialCase) string {
		return strings.ToTitleSpecial(*c, s)
	},
	"map": func(mapping func(rune) rune, s string) string {
		return strings.Map(mapping, s)
	},
	"fields": func(s string) []string {
		return strings.Fields(s)
	},
	"fieldsFunc": func(s string, f func(rune) bool) []string {
		return strings.FieldsFunc(s, f)
	},
	"trimPrefix": func(s, prefix string) string {
		return strings.TrimPrefix(s, prefix)
	},
	"trimSuffix": func(s, suffix string) string {
		return strings.TrimSuffix(s, suffix)
	},
	"trimSpace": func(s string) string {
		return strings.TrimSpace(s)
	},
	"trimLeft": func(s, cutset string) string {
		return strings.TrimLeft(s, cutset)
	},
	"trimRight": func(s, cutset string) string {
		return strings.TrimRight(s, cutset)
	},
	"trimFunc": func(s string, f func(rune) bool) string {
		return strings.TrimFunc(s, f)
	},
	"compare": func(a, b string) float64 {
		return float64(strings.Compare(a, b))
	},
	"count": func(s, substr string) float64 {
		return float64(strings.Count(s, substr))
	},
}

func RegisterStringFuncs(env *object.Environment) {
	object.RegisterFunctions(env, "String", stringFuncs)
}
