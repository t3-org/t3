package gutil

// PanicNil will panic error param if value is nil.
func PanicNil(val interface{}, err error) {
	if IsNil(val) {
		panic(err)
	}
}

// PanicErr panic if error value is not nil.
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
