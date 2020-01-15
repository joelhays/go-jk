package jkparsers

import (
	"github.com/joelhays/go-jk/jk/jktypes"
)

type CmpParser struct {
}

func NewCmpParser() *CmpParser {
	return &CmpParser{}
}

func (p *CmpParser) ParseFromBytes(data []byte) jktypes.ColorMap {
	cursor := 0
	var header jktypes.TCMPHeader
	cursor += readBytes(data, cursor, &header)

	return jktypes.ColorMap{Palette: header.Palette}
}
