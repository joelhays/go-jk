package jk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"log"
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

func parseVec3(line string) (int32, mgl32.Vec3) {
	var id int32
	v := mgl32.Vec3{}
	n, err := fmt.Sscanf(line, "%d: %f %f %f", &id, &v[0], &v[1], &v[2])
	if err != nil {
		log.Fatal(err)
	}
	if n != 4 {
		panic("Unable to get vec3 from line: " + line)
	}

	return id, v
}

func parseVec2(line string) (int32, mgl32.Vec2) {
	var id int32
	v := mgl32.Vec2{}
	n, err := fmt.Sscanf(line, "%d: %f %f", &id, &v[0], &v[1])
	if err != nil {
		log.Fatal(err)
	}
	if n != 3 {
		panic("Unable to get vec2 from line: " + line)
	}

	return id, v
}
