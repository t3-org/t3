package model

import "github.com/sony/sonyflake"

var sf *sonyflake.Sonyflake

func SetIDGenerator(flake *sonyflake.Sonyflake) {
	sf = flake
}

func genId() int64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return int64(id) // id is 63 bits, so we can convert it to int64 safely.
}
