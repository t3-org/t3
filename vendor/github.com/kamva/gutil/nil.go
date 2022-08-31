package gutil

// AnyNil check if any value of provided values is nil, returns true, otherwise returns false.
func AnyNil(val ...interface{}) bool {
	for _, v := range val {
		if IsNil(v) {
			return true
		}
	}

	return false
}
