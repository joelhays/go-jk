package jk

type JkParser interface {
	ParseFromFile(data string) interface{}
	ParseFromString(data string) interface{}
}

type Jk3doParser interface {
	ParseFromFile(data string) Jk3doFile
	ParseFromString(data string) Jk3doFile
}

type JklParser interface {
	ParseFromFile(filePath string) Jkl
	ParseFromString(jklString string) Jkl
}

type KeyParser interface {
	ParseFromFile(filePath string) Key
	ParseFromString(jklString string) Key
}

type PupParser interface {
	ParseFromFile(filePath string) Pup
	ParseFromString(jklString string) Pup
}
