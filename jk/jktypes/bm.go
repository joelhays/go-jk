package jktypes

type BMFile struct {
	Header  TBMHeader
	Images  []TImage
	Palette TPalette
}

type TBMHeader struct {
	FileType        [3]byte
	Ver             byte
	Unknown1        int32
	Unknown2        int32
	PaletteIncluded int32
	NumImages       int32
	XOffset         int32
	YOffset         int32
	Transparent     int32
	Unknown3        int32
	NumBits         int32
	BlueBits        int32
	GreenBits       int32
	RedBits         int32
	Unknown4        int32
	Unknown5        int32
	Unknown6        int32
	Unknown7        int32
	Unknown8        int32
	Unknown9        int32
	Padding         [13]int32
}

type TImage struct {
	SizeX int32
	SizeY int32
	Data  []byte
}

type TPalette struct {
	Palette [256]Vec3Byte
}
