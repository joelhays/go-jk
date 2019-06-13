package jk

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

func parseBmFile(data []byte) BMFile {
	result := BMFile{}

	cursor := 0
	var header TBMHeader
	cursor += readBytes(data, cursor, &header)

	result.Header = header
	result.Images = make([]TImage, header.NumImages)

	for i := int32(0); i < header.NumImages; i++ {
		//fmt.Printf("reading image %d of %d\n", i+1, header.NumImages)

		imageSize := struct {
			SizeX int32
			SizeY int32
		}{
			0,
			0,
		}
		cursor += readBytes(data, cursor, &imageSize)

		var image TImage
		image.SizeX = imageSize.SizeX
		image.SizeY = imageSize.SizeY

		if image.SizeX < 0 || image.SizeY < 0 {
			continue
		}

		if header.NumBits == 8 {
			//fmt.Printf("8-bit image %d pixels %+v\n", image.SizeX*image.SizeY, image)
			image.Data = make([]byte, image.SizeX*image.SizeY)
		} else {
			//fmt.Printf("16-bit image %d pixels %+v\n", image.SizeX*image.SizeY, image)
			image.Data = make([]byte, image.SizeX*image.SizeY*2)
		}

		cursor += readBytes(data, cursor, &image.Data)

		result.Images[i] = image
	}

	if header.PaletteIncluded == 2 {
		var palette TPalette
		readBytes(data, cursor, &palette)
		result.Palette = palette
	}

	if header.PaletteIncluded != 2 {
		cmp := GetLoader().LoadCMP("dflt.cmp")
		result.Palette.Palette = cmp.Palette
	}

	return result
}
