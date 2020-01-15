package jk

import (
	"fmt"
	"github.com/joelhays/go-jk/jk/jktypes"
	"log"
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
	cmpParser   *CmpParser
	bmParser    *BmParser
	matParser   *MatParser
	sftParser   *SftParser
}

func GetLoader() *Loader {
	once.Do(func() {
		instance = &Loader{}
		instance.cache = make(map[string]interface{})
		instance.jklParser = NewJklLineParser()
		instance.jk3doParser = NewJk3doLineParser()
		instance.cmpParser = NewCmpParser()
		instance.bmParser = NewBmParser()
		instance.matParser = NewMatParser()
		instance.sftParser = NewSftParser()
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

func (l *Loader) LoadResourceManifest(prefix string, suffix string) []string {
	return l.getGobFiles(resourceGobFiles, prefix, suffix)
}

func (l *Loader) LoadEpisodeManifest(prefix string, suffix string) []string {
	return l.getGobFiles(episodeGobFiles, prefix, suffix)
}

func (l *Loader) LoadJKL(filename string) jktypes.Jkl {
	fileBytes := l.LoadEpisode(filename)
	if fileBytes == nil {
		return jktypes.Jkl{}
	}

	jklLevel := l.jklParser.ParseFromString(string(fileBytes))
	return jklLevel
}

func (l *Loader) Load3DO(filename string) jktypes.Jk3doFile {
	var obj jktypes.Jk3doFile

	if obj, ok := l.cache[filename]; ok {
		return obj.(jktypes.Jk3doFile)
	}

	fileBytes := l.LoadResource(filename)
	obj = l.jk3doParser.ParseFromString(string(fileBytes))
	l.cache[filename] = obj
	return obj
}

func (l *Loader) LoadMAT(filename string) jktypes.Material {
	var mat jktypes.Material

	if mat, ok := l.cache[filename]; ok {
		return mat.(jktypes.Material)
	}

	fileBytes := l.LoadResource(filename)
	if fileBytes == nil {
		return jktypes.Material{}
	}

	mat = l.matParser.ParseFromBytes(fileBytes)
	l.cache[filename] = mat
	return mat
}

func (l *Loader) LoadCMP(filename string) jktypes.ColorMap {
	var cmp jktypes.ColorMap

	if cmp, ok := l.cache[filename]; ok {
		return cmp.(jktypes.ColorMap)
	}

	fileBytes := l.LoadResource(filename)
	if fileBytes == nil {
		return jktypes.ColorMap{}
	}

	cmp = l.cmpParser.ParseFromBytes(fileBytes)
	l.cache[filename] = cmp
	return cmp
}

func (l *Loader) LoadBM(filename string) jktypes.BMFile {
	var bm jktypes.BMFile

	if bm, ok := l.cache[filename]; ok {
		return bm.(jktypes.BMFile)
	}

	fileBytes := l.LoadResource(filename)
	if fileBytes == nil {
		return jktypes.BMFile{}
	}

	bm = l.bmParser.ParseFromBytes(fileBytes)
	l.cache[filename] = bm
	return bm
}

func (l *Loader) LoadSFT(filename string) jktypes.SFTFile {
	var sft jktypes.SFTFile

	if sft, ok := l.cache[filename]; ok {
		return sft.(jktypes.SFTFile)
	}

	fileBytes := l.LoadResource(filename)
	if fileBytes == nil {
		return jktypes.SFTFile{}
	}

	sft = l.sftParser.ParseFromBytes(fileBytes)
	l.cache[filename] = sft
	return sft
}

func (l *Loader) LoadResource(filename string) []byte {
	for _, gob := range resourceGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		return fileBytes
	}

	log.Println(fmt.Errorf("unable to find %s", filename))
	return nil
}

func (l *Loader) LoadEpisode(filename string) []byte {
	for _, gob := range episodeGobFiles {
		fileBytes := loadFileFromGOB(gob, filename)
		if fileBytes == nil {
			continue
		}
		return fileBytes
	}

	log.Println(fmt.Errorf("unable to find %s", filename))
	return nil
}
