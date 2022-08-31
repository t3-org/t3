package hexa

import "fmt"

// ID is the id of entities across the hexa packages.
// This is because hexa does not want to be dependent
// to specific type of id. (e.g mongo ObjectID, mysql integer,...)
type ID interface {
	fmt.Stringer

	// Validate specify provided id is valid or not.
	Validate(id any) error

	// From convert provided value to its id.
	// From will returns error if provided value
	// can not convert to an native id.
	From(id any) error

	// MustFrom Same as FromString but on occur error, it will panic.
	MustFrom(id any)

	// Val returns the native id value (e.g ObjectID in mongo, ...).
	Val() any

	// IsEqual say that two hexa id are equal or not.
	IsEqual(ID) bool
}
