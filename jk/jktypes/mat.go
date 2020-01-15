package jktypes

type MtlHeader struct {
	Name         [4]byte
	Ver          int32
	MatType      int32
	NumTextures  int32
	NumTextures1 int32
	Unk0         int32
	Unk1         int32
	Unk2         [12]int32
}

type TextureHeader struct {
	TexType      int32
	ColorNum     int32
	Unk0         [4]int32
	Unk1         [2]int32
	Unk2         int32
	CurrentTXNum int32
}

type ColorHeader struct {
	TexType  int32
	ColorNum int32
	Unk0     [4]int32
}

type TextureData struct {
	SizeX      int32
	SizeY      int32
	Pad        [3]int32
	NumMipMaps int32
}

type Material struct {
	Texture     []byte
	SizeX       int32
	SizeY       int32
	XTile       float32
	YTile       float32
	Transparent bool
}
