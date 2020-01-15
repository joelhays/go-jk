package jk

import "github.com/joelhays/go-jk/jk/jktypes"

type SftParser struct {
}

func NewSftParser() *SftParser {
	return &SftParser{}
}

func (p *SftParser) ParseFromBytes(data []byte) jktypes.SFTFile {
	result := jktypes.SFTFile{}

	cursor := 0
	var header jktypes.TSFTHeader
	cursor += readBytes(data, cursor, &header)

	result.Header = header
	result.CharacterTables = make([]jktypes.TCharacterTable, header.NumTables)

	bmParser := NewBmParser()

	for i := int32(0); i < header.NumTables; i++ {
		//fmt.Println("reading table", i+1, "of", header.NumTables)

		tableInfo := struct {
			FirstChar int16
			LastChar  int16
		}{
			0,
			0,
		}
		cursor += readBytes(data, cursor, &tableInfo)

		var table jktypes.TCharacterTable
		table.FirstChar = tableInfo.FirstChar
		table.LastChar = tableInfo.LastChar
		table.CharDefs = make([]jktypes.TCharDef, table.LastChar-table.FirstChar+1)

		for j := 0; j < len(table.CharDefs); j++ {
			var def jktypes.TCharDef
			cursor += readBytes(data, cursor, &def)
			table.CharDefs[j] = def
		}

		result.CharacterTables[i] = table
	}

	bm := bmParser.ParseFromBytes(data[cursor:])
	if bm.Header.PaletteIncluded != 2 {
		cmp := GetLoader().LoadCMP("uicolormap.cmp")
		bm.Palette.Palette = cmp.Palette
	}

	result.BMFile = bm

	return result
}
