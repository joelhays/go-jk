package jk

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strings"
	"unsafe"
)

type GOB struct {
	Header GOBHeader
	Items  []GOBItem
}

type GOBHeader struct {
	FileType        [3]byte
	Version         byte
	FirstFileOffset int32
	NumItemsOffset  int32
	NumItems        int32
}

type GOBItem struct {
	FileOffset uint32
	FileLength uint32
	FileName   [128]byte
}

func loadGOB(gobPath string) GOB {
	file, err := os.Open(gobPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := GOB{}

	var header GOBHeader
	data := make([]byte, unsafe.Sizeof(header))
	file.Read(data)
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &header)

	result.Header = header

	for i := int32(0); i < header.NumItems; i++ {
		var item GOBItem
		data := make([]byte, unsafe.Sizeof(item))
		file.Read(data)
		buf := bytes.NewBuffer(data)
		binary.Read(buf, binary.LittleEndian, &item)

		result.Items = append(result.Items, item)
	}

	return result
}

func loadFileFromGOB(gobPath string, fileName string) []byte {

	gob := loadGOB(gobPath)

	file, err := os.Open(gobPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, item := range gob.Items {
		pathAndName := string(item.FileName[:])
		if strings.Contains(strings.ToUpper(pathAndName), strings.ToUpper(fileName)) {
			file.Seek(int64(item.FileOffset), io.SeekStart)
			contentBytes := make([]byte, item.FileLength)
			file.Read(contentBytes)
			return contentBytes
		}
	}

	return nil
}
