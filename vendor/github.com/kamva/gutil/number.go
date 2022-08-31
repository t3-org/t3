package gutil

import (
	"strconv"
	"time"
)

// ParseInt parse integer and returns value if string
// value is a valid integer or default value
func ParseInt(val string, def int) int {
	if val == "" {
		return def
	}
	if result, err := strconv.Atoi(val); err == nil {
		return result
	}
	return def
}

// ParseUint64 parse unsigned integer and returns it
// if the value is valid, otherwise returns the
// default value.
func ParseUint64(val string, def uint64) uint64 {
	if val == "" {
		return def
	}
	if result, err := strconv.ParseUint(val, 10, 64); err == nil {
		return result
	}
	return def
}

func Min(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func Min64(vars ...int64) int64 {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func MinDuration(vars ...time.Duration) time.Duration {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func Max(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func Max64(vars ...int64) int64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func MaxDuration(vars ...time.Duration) time.Duration {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}
