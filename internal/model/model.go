package model

import (
	"github.com/kamva/gutil"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func SetIDGenerator(flake *sonyflake.Sonyflake) {
	sf = flake
}

// We've disabled it because when we send these ids to the browser, it converts
// them to another number (we should use BigInt to fix it I think)
//func genId() int64 {
//id, err := sf.NextID()
//if err != nil {
//	panic(err)
//}
//return int64(id) // id is 63 bits, so we can convert it to int64 safely.
//}

func genId() string {
	return gutil.UUID()
}
