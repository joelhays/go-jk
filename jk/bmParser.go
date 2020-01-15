package jk

import "github.com/joelhays/go-jk/jk/jktypes"

type BmParser struct {
}

func NewBmParser() *BmParser {
	return &BmParser{}
}

func (p *BmParser) ParseFromBytes(data []byte) jktypes.BMFile {
	result := jktypes.BMFile{}

	cursor := 0
	var header jktypes.TBMHeader
	cursor += readBytes(data, cursor, &header)

	result.Header = header
	result.Images = make([]jktypes.TImage, header.NumImages)

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

		var image jktypes.TImage
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
		var palette jktypes.TPalette
		readBytes(data, cursor, &palette)
		result.Palette = palette
	}

	if header.PaletteIncluded != 2 {
		cmp := GetLoader().LoadCMP("dflt.cmp")
		result.Palette.Palette = cmp.Palette
	}

	return result
}
