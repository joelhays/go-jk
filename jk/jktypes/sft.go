package jktypes

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
