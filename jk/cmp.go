package jk

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

func parseCmpFile(data []byte) ColorMap {
	cursor := 0
	var header TCMPHeader
	cursor += readBytes(data, cursor, &header)

	return ColorMap{Palette: header.Palette}
}
