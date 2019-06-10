package jk

type SFTFile struct {
	Header          TSFTHeader
	CharacterTables []TCharacterTable
	BMFile          BMFile
}

type TSFTHeader struct {
	FileType  [4]byte
	_         [4]int32
	NumTables int32
	Padding   [4]int32
}

type TCharacterTable struct {
	FirstChar int16
	LastChar  int16
	CharDefs  []TCharDef
}

type TCharDef struct {
	XOffset int32
	Width   int32
}

func parseSFTFile(data []byte) SFTFile {
	result := SFTFile{}

	cursor := 0
	var header TSFTHeader
	cursor += readBytes(data, cursor, &header)

	result.Header = header
	result.CharacterTables = make([]TCharacterTable, header.NumTables)

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

		var table TCharacterTable
		table.FirstChar = tableInfo.FirstChar
		table.LastChar = tableInfo.LastChar
		table.CharDefs = make([]TCharDef, table.LastChar-table.FirstChar+1)

		for j := 0; j < len(table.CharDefs); j++ {
			var def TCharDef
			cursor += readBytes(data, cursor, &def)
			table.CharDefs[j] = def
		}

		result.CharacterTables[i] = table
	}

	bm := parseBmFile(data[cursor:])
	if bm.Header.PaletteIncluded != 2 {
		cmp := GetLoader().LoadCMP("uicolormap.cmp")
		bm.Palette.Palette = cmp.Palette
	}

	result.BMFile = bm

	return result
}
