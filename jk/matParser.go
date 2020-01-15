package jk

import (
	"encoding/binary"
	"github.com/joelhays/go-jk/jk/jktypes"
)

type MatParser struct{}

func NewMatParser() *MatParser {
	return &MatParser{}
}

func (p *MatParser) ParseFromBytes(data []byte) jktypes.Material {
	cursor := 0
	var header jktypes.MtlHeader
	cursor += readBytes(data, cursor, &header)

	if header.MatType == 0 {
		// TODO: handle color-only materials
		var colHeader jktypes.ColorHeader
		cursor += readBytes(data, cursor, &colHeader)

		texture := make([]byte, 4)
		binary.LittleEndian.PutUint32(texture, uint32(colHeader.ColorNum))

		return jktypes.Material{Texture: texture, SizeX: 1, SizeY: 1, Transparent: false}
	}

	if header.MatType == 2 {
		// get the first texture header (full-sized image)
		var texHeader jktypes.TextureHeader
		cursor += readBytes(data, cursor, &texHeader) * int(header.NumTextures)

		var texData jktypes.TextureData
		cursor += readBytes(data, cursor, &texData)

		textureBytes := data[cursor : cursor+int(texData.SizeX*texData.SizeY)]

		var transparent bool
		for i := 0; i < len(textureBytes); i++ {
			if textureBytes[i] == 0 {
				transparent = true
				break
			}
		}

		return jktypes.Material{Texture: textureBytes, SizeX: texData.SizeX, SizeY: texData.SizeY, Transparent: transparent}
	}

	return jktypes.Material{}
}
