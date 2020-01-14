package jk

import (
	"strings"
	"sync"
)

var (
	instance         *Loader
	once             sync.Once
	resourceGobFiles = []string{"J:\\Resource\\Res2.gob", "J:\\Resource\\Res1hi.gob"}
	episodeGobFiles  = []string{"J:\\Episode\\JK1.GOB", "J:\\Episode\\JK1CTF.GOB", "J:\\Episode\\JK1MP.GOB"}
)

type Loader struct {
	cache       map[string]interface{}
	jklParser   JklParser
	jk3doParser Jk3doParser
}

func GetLoader() *Loader {
	once.Do(func() {
		instance = &Loader{}
		instance.cache = make(map[string]interface{})
		instance.jklParser = NewJklLineParser()
		instance.jk3doParser = NewJk3doLineParser2()
	})
	return instance
}

func (l *Loader) getGobFiles(gobFiles []string, prefix string, suffix string) []string {
	var files []string
	for _, gob := range gobFiles {
		for _, gobData := range loadGOBManifest(gob).Items {
			if strings.HasPrefix(gobData.FileName, prefix) && strings.HasSuffix(gobData.FileName, suffix) {
				files = append(files, gobData.FileName)
			}
		}
	}

	return files
}

func (l *Loader) LoadJKLManifest() []string {
	return l.getGobFiles(episodeGobFiles, "jkl\\", "jkl")
}

func (l *Loader) LoadBMManifest() []string {
	return l.getGobFiles(resourceGobFiles, "ui\\bm\\", "bm")
}

func (l *Loader) Load3DOManifest() []string {
	return l.getGobFiles(resourceGobFiles, "3do\\", "3do")
}

func (l *Loader) LoadJKL(filename string) Jkl {
	for _, gob := range episodeGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		jklLevel := l.jklParser.ParseJKLFromString(string(fileBytes))
		return jklLevel
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		jklLevel := l.jklParser.ParseJKLFromString(string(fileBytes))
		return jklLevel
	}

	return Jkl{}
}

func (l *Loader) Load3DO(filename string) Jk3doFile {
	var obj Jk3doFile

	if obj, ok := l.cache[filename]; ok {
		return obj.(Jk3doFile)
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		obj = l.jk3doParser.Parse3doFromString(string(fileBytes))
		l.cache[filename] = obj
		return obj
	}

	return Jk3doFile{}
}

func (l *Loader) LoadMAT(filename string) Material {
	var mat Material

	if mat, ok := l.cache[filename]; ok {
		return mat.(Material)
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		mat = parseMatFile(fileBytes)
		l.cache[filename] = mat
		return mat
	}

	return Material{}
}

func (l *Loader) LoadCMP(filename string) ColorMap {
	var cmp ColorMap

	if cmp, ok := l.cache[filename]; ok {
		return cmp.(ColorMap)
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		cmp = parseCmpFile(fileBytes)
		l.cache[filename] = cmp
		return cmp
	}

	return ColorMap{}
}

func (l *Loader) LoadBM(filename string) BMFile {
	var bm BMFile

	if bm, ok := l.cache[filename]; ok {
		return bm.(BMFile)
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		bm = parseBmFile(fileBytes)
		l.cache[filename] = bm
		return bm
	}

	return BMFile{}
}

func (l *Loader) LoadSFT(filename string) SFTFile {
	var sft SFTFile

	if sft, ok := l.cache[filename]; ok {
		return sft.(SFTFile)
	}

	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		sft = parseSFTFile(fileBytes)
		l.cache[filename] = sft
		return sft
	}

	return SFTFile{}
}
