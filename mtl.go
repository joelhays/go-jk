package main

import (
	"bytes"
	"encoding/binary"
	"unsafe"
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

type textureData struct {
	SizeX      int32
	SizeY      int32
	Pad        [3]int32
	NumMipMaps int32
}

type mtl struct {
	header mtlHeader
	textureHeader
	textureData
}

func parseMatFile(data []byte) material {
	cursor := 0
	var header mtlHeader
	headerSize := int(unsafe.Sizeof(header))
	headerBuf := bytes.NewBuffer(data[cursor:headerSize])
	binary.Read(headerBuf, binary.LittleEndian, &header)
	// fmt.Println(header, string(header.Name[:4]))

	cursor += headerSize

	if header.MatType == 2 {
		var texHeader textureHeader
		texHeaderSize := int(unsafe.Sizeof(texHeader))
		texBuf := bytes.NewBuffer(data[cursor : cursor+texHeaderSize])
		binary.Read(texBuf, binary.LittleEndian, &texHeader)
		// fmt.Println(texHeader)

		cursor += texHeaderSize

		var texData textureData
		texDataSize := int(unsafe.Sizeof(texData))
		texDataBuf := bytes.NewBuffer(data[cursor : cursor+texDataSize])
		binary.Read(texDataBuf, binary.LittleEndian, &texData)
		// fmt.Println(texData)

		cursor += texDataSize

		textureBytes := data[cursor : cursor+int(texData.SizeX)*int(texData.SizeY)]
		return material{Texture: textureBytes, SizeX: texData.SizeX, SizeY: texData.SizeY}
	}

	return material{}
}
