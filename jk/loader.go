package jk

import (
	"fmt"
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
}

func GetLoader() *Loader {
	once.Do(func() {
		instance = &Loader{}
	})
	return instance
}

func (l *Loader) getGobFiles(gobFiles []string, suffix string) []string {
	var files []string
	for _, gob := range gobFiles {
		for _, gobData := range loadGOBManifest(gob).Items {
			if strings.HasSuffix(gobData.FileName, suffix) {
				files = append(files, gobData.FileName)
			}
		}
	}

	return files
}

func (l *Loader) LoadManifest(resourceType string) []string {
	return l.getGobFiles(resourceGobFiles, "."+resourceType)
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
