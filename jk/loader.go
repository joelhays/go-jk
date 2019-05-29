package jk

import "sync"

var (
	instance         *Loader
	once             sync.Once
	resourceGobFiles = []string{"J:\\Resource\\Res2.gob", "J:\\Resource\\Res1hi.gob"}
	episodeGobFiles  = []string{"J:\\Episode\\JK1.GOB", "J:\\Episode\\JK1CTF.GOB", "J:\\Episode\\JK1MP.GOB"}
)

type Loader struct {
	objCache map[string]Jk3doFile
	matCache map[string]Material
	cmpCache map[string]ColorMap
	bmCache  map[string]BMFile
}

func GetLoader() *Loader {
	once.Do(func() {
		instance = &Loader{}
		instance.objCache = make(map[string]Jk3doFile)
		instance.matCache = make(map[string]Material)
		instance.cmpCache = make(map[string]ColorMap)
		instance.bmCache = make(map[string]BMFile)
	})
	return instance
}

func (l *Loader) LoadJKL(filename string) Jkl {
	for _, gob := range episodeGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		jklLevel := readJKLFromString(string(fileBytes))
		return jklLevel
	}

	return Jkl{}
}

func (l *Loader) Load3DO(filename string) Jk3doFile {
	var obj Jk3doFile

	if obj, ok := l.objCache[filename]; ok {
		return obj
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		obj = Parse3doFile(string(fileBytes))
		l.objCache[filename] = obj
		return obj
	}

	return Jk3doFile{}
}

func (l *Loader) LoadMAT(filename string) Material {
	var mat Material

	if mat, ok := l.matCache[filename]; ok {
		return mat
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		mat = parseMatFile(fileBytes)
		l.matCache[filename] = mat
		return mat
	}

	return Material{}
}

func (l *Loader) LoadCMP(filename string) ColorMap {
	var cmp ColorMap

	if cmp, ok := l.cmpCache[filename]; ok {
		return cmp
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		cmp = parseCmpFile(fileBytes)
		l.cmpCache[filename] = cmp
		return cmp
	}

	return ColorMap{}
}

func (l *Loader) LoadBM(filename string) BMFile {
	var bm BMFile

	if bm, ok := l.bmCache[filename]; ok {
		return bm
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		bm = parseBmFile(fileBytes)
		l.bmCache[filename] = bm
		return bm
	}

	return BMFile{}
}
