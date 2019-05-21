package jk

import (
	"encoding/binary"
)

type mtlHeader struct {
	Name         [4]byte
	Ver          int32
	MatType      int32
	NumTextures  int32
	NumTextures1 int32
	Unk0         int32
	Unk1         int32
	Unk2         [12]int32
}

type textureHeader struct {
	TexType      int32
	ColorNum     int32
	Unk0         [4]int32
	Unk1         [2]int32
	Unk2         int32
	CurrentTXNum int32
}

type colorHeader struct {
	TexType  int32
	ColorNum int32
	Unk0     [4]int32
}

type textureData struct {
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

func ParseMatFile(data []byte) Material {
	cursor := 0
	var header mtlHeader
	cursor += readBytes(data, cursor, &header)

	if header.MatType == 0 {
		// TODO: handle color-only materials
		var colHeader colorHeader
		cursor += readBytes(data, cursor, &colHeader)

		texture := make([]byte, 4)
		binary.LittleEndian.PutUint32(texture, uint32(colHeader.ColorNum))

		return Material{Texture: texture, SizeX: 1, SizeY: 1, Transparent: false}
	}

	if header.MatType == 2 {
		// get the first texture header (full-sized image)
		var texHeader textureHeader
		cursor += readBytes(data, cursor, &texHeader) * int(header.NumTextures)

		var texData textureData
		cursor += readBytes(data, cursor, &texData)

		textureBytes := data[cursor : cursor+int(texData.SizeX*texData.SizeY)]

		var transparent bool
		for i := 0; i < len(textureBytes); i++ {
			if textureBytes[i] == 0 {
				transparent = true
				break
			}
		}

		return Material{Texture: textureBytes, SizeX: texData.SizeX, SizeY: texData.SizeY, Transparent: transparent}
	}

	return Material{}
}
