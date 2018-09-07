package main

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
)

type GOB struct {
}

func LoadFileFromGOB(gobPath string, fileName string) []byte {
	file, err := os.Open(gobPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read header
	// GOBHeader=record
	// array[0..2] of char;    {'GOB '}
	// byte;                   {Apparently - version= x20 }
	// Longint;                {Offset to first file size from begining of file= x14 }
	// Longint;                {Offset to #of GobItems from file beginning= x0C }
	// Longint;                {# of items in gob file }
	headerBytes := make([]byte, 3)
	file.Read(headerBytes)

	versionBytes := make([]byte, 1)
	file.Read(versionBytes)

	firstFileOffsetBytes := make([]byte, 4)
	file.Read(firstFileOffsetBytes)

	numItemsOffsetBytes := make([]byte, 4)
	file.Read(numItemsOffsetBytes)

	numItemsBytes := make([]byte, 4)
	file.Read(numItemsBytes)
	numItems := binary.LittleEndian.Uint32(numItemsBytes)

	// fmt.Println("Header", string(headerBytes))
	// fmt.Println("First File Offset", binary.LittleEndian.Uint32(firstFileOffsetBytes))
	// fmt.Println("Number Of Items Offset", binary.LittleEndian.Uint32(numItemsOffsetBytes))
	// fmt.Println("Number of Items", numItems)

	// read items
	// GobItems =record
	// Longint:              {offset from begining of file}
	// Longint;              {Length of file}
	// array[0..127] of char {path and name of file}

	for i := uint32(0); i < numItems; i++ {
		offsetBytes := make([]byte, 4)
		file.Read(offsetBytes)
		offset := binary.LittleEndian.Uint32(offsetBytes)

		fileLengthBytes := make([]byte, 4)
		file.Read(fileLengthBytes)
		fileLength := binary.LittleEndian.Uint32(fileLengthBytes)

		pathAndNameBytes := make([]byte, 128)
		file.Read(pathAndNameBytes)
		pathAndName := string(pathAndNameBytes)

		// fmt.Println("Offset", offset)
		// fmt.Println("File Length", fileLength)
		// fmt.Println("Path and Name", pathAndName)

		if strings.Contains(strings.ToUpper(pathAndName), strings.ToUpper(fileName)) {
			file.Seek(int64(offset), io.SeekStart)
			contentBytes := make([]byte, fileLength)
			file.Read(contentBytes)
			return contentBytes
		}

	}

	return nil
}
