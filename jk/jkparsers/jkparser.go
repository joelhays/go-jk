package jkparsers

import "github.com/joelhays/go-jk/jk/jktypes"

type JkParser interface {
	ParseFromFile(data string) interface{}
	ParseFromString(data string) interface{}
}

type Jk3doParser interface {
	ParseFromFile(data string) jktypes.Jk3doFile
	ParseFromString(data string) jktypes.Jk3doFile
}

type JklParser interface {
	ParseFromFile(filePath string) jktypes.Jkl
	ParseFromString(jklString string) jktypes.Jkl
}

type KeyParser interface {
	ParseFromFile(filePath string) jktypes.Key
	ParseFromString(jklString string) jktypes.Key
}

type PupParser interface {
	ParseFromFile(filePath string) jktypes.Pup
	ParseFromString(jklString string) jktypes.Pup
}
