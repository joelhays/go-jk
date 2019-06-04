package jk

import "fmt"

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

	fmt.Println("parsing bm file", len(data))

	cursor := 0
	var header TBMHeader
	cursor += readBytes(data, cursor, &header)

	result.Header = header
	result.Images = make([]TImage, header.NumImages)

	for i := int32(0); i < header.NumImages; i++ {
		fmt.Println("reading image", i+1, "of", header.NumImages)

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
			fmt.Println("8-bit image", image, image.SizeX*image.SizeY)
			// 8-bit image
			image.Data = make([]byte, image.SizeX*image.SizeY)
		} else {
			fmt.Println("16-bit image", image, image.SizeX*image.SizeY)
			// 16-bit image
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

	return result
}
