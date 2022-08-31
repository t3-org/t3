package gutil

import "sort"

// Contains check that needle exists in the haystack or not.
func Contains(haystack []string, needle string) bool {
	sort.Strings(haystack)
	i := sort.SearchStrings(haystack, needle)
	return i < len(haystack) && haystack[i] == needle
}

// Other methods which we can implement: ContainsInt