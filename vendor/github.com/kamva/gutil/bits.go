package gutil

type Bits uint32

func (b Bits) Has(f Bits) bool {
	return b&f != 0
}

func (b Bits) Set(f Bits) Bits {
	return b | f
}

func (b Bits) Toggle(f Bits) Bits {
	return b ^ f
}

func (b Bits) Clear(f Bits) Bits {
	return b &^ f
}
