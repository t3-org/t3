package gutil

import (
	"regexp"
	"strings"
	"unicode"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// StringDefault return default value if provided value is empty.
func StringDefault(val, def string) string {
	if val == "" {
		return def
	}

	return val
}

// ToSnakeCase returns snake_case of the provided value.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ReplaceAt replace provided string in specific index to another index.
// removal part is [begin,end).
// e.g., to replace char "a" in the "salam" word, call ReplaceAt("salam","b",1,2) // result will be "sblam".
func ReplaceAt(str string, replace string, begin int, end int) string {
	return str[:begin] + replace + str[end:]
}

// ReplaceRune replace single rune in specific index
func ReplaceRune(str string, new rune, index int) string {
	return ReplaceAt(str, string(new), index, index+1)
}

// AnyString returns first found non-empty string value
// in the provided values.
func AnyString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// Sub subtracts two slices. returns s1 - s2.
// e.g., [1,2,3] - [2,3,4] = [1]
func Sub(s1 []string, s2 []string) []string {
	mb := make(map[string]struct{}, len(s2))
	for _, v := range s2 {
		mb[v] = struct{}{}
	}

	diff := make([]string, 0)
	for _, v := range s1 {
		if _, ok := mb[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

// Intersect returns intersect of two slices.
// e.g., intersect of [1,2,3] & [2,3,4] = [2,3]
func Intersect(s1 []string, s2 []string) []string {
	mb := make(map[string]struct{}, len(s2))
	for _, v := range s2 {
		mb[v] = struct{}{}
	}

	intersect := make([]string, 0)
	for _, v := range s1 {
		if _, ok := mb[v]; ok {
			intersect = append(intersect, v)
		}
	}
	return intersect
}

// LowerFirst converts first letter of string to lowercase.
func LowerFirst(v string) string {
	for i, r := range v { // run loop to get first rune.
		return string(unicode.ToLower(r)) + v[i+1:]
	}

	return ""
}
