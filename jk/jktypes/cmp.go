package jktypes

type TCMPHeader struct {
	Name         [4]byte
	Ver          int32
	Transparency int32
	Padding      [52]byte
	Palette      [256]Vec3Byte
}

type Vec3Byte struct {
	R byte
	G byte
	B byte
}

type ColorMap struct {
	Palette [256]Vec3Byte
}
