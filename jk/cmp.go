package jk

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

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
	Pallette [256]Vec3Byte
}

func ParseCmpFile(data []byte) ColorMap {
	cursor := 0
	var header TCMPHeader
	headerSize := int(unsafe.Sizeof(header))
	headerBuf := bytes.NewBuffer(data[cursor:headerSize])
	binary.Read(headerBuf, binary.LittleEndian, &header)

	return ColorMap{Pallette: header.Palette}
}
