package jk

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func readBytes(data []byte, cursor int, object interface{}) int {
	value := reflect.Indirect(reflect.ValueOf(object))
	var sizeInBytes int
	if value.Kind() == reflect.Slice {
		sizeInBytes = value.Len()
	} else {
		sizeInBytes = int(value.Type().Size())
	}
	buf := bytes.NewBuffer(data[cursor : cursor+sizeInBytes])
	binary.Read(buf, binary.LittleEndian, object)
	return sizeInBytes
}
