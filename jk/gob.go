package jk

import (
	"bufio"
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
	FileOffset    uint32
	FileLength    uint32
	FileName      string
	UpperFileName string
}

var (
	gobManifestCache = make(map[string]GOB)
)

func loadGOBManifest(gobPath string) GOB {
	if obj, ok := gobManifestCache[gobPath]; ok {
		return obj
	}

	file, err := os.Open(gobPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := GOB{}

	fr := bufio.NewReader(file)

	var header GOBHeader
	data := make([]byte, unsafe.Sizeof(header))
	io.ReadFull(fr, data)
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &header)

	result.Header = header

	for i := int32(0); i < header.NumItems; i++ {
		tempitem := struct {
			FileOffset uint32
			FileLength uint32
			FileName   [128]byte
		}{}

		data := make([]byte, unsafe.Sizeof(tempitem))
		io.ReadFull(fr, data)
		buf := bytes.NewBuffer(data)
		binary.Read(buf, binary.LittleEndian, &tempitem)

		var item GOBItem
		item.FileOffset = tempitem.FileOffset
		item.FileLength = tempitem.FileLength

		filenameBytes := bytes.Split(tempitem.FileName[:], []byte{byte('\x00')})[0]
		item.FileName = string(filenameBytes)
		item.UpperFileName = strings.ToUpper(item.FileName)

		result.Items = append(result.Items, item)
	}

	gobManifestCache[gobPath] = result

	return result
}

func loadFileFromGOB(gobPath string, fileName string) []byte {

	gob := loadGOBManifest(gobPath)

	file, err := os.Open(gobPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileName = strings.ToUpper(fileName)

	for _, item := range gob.Items {
		if strings.Contains(item.UpperFileName, fileName) {
			file.Seek(int64(item.FileOffset), io.SeekStart)
			contentBytes := make([]byte, item.FileLength)
			file.Read(contentBytes)
			return contentBytes
		}
	}

	return nil
}
