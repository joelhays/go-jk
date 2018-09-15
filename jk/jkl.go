package jk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	GobFiles = []string{"J:\\Resource\\Res2.gob"} //"J:\\Episode\\JK1.GOB", "J:\\Episode\\JK1CTF.GOB", "J:\\Episode\\JK1MP.GOB", "J:\\Resource\\Res2.gob", "J:\\Resource\\Res1hi.gob"}
)

// Jkl contains the information extracted from the Jedi Knight Level (.jkl) file
type Jkl struct {
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	Surfaces        []surface
	Materials       []Material
	ColorMaps       []ColorMap
}

type surface struct {
	VertexIds        []int64
	TextureVertexIds []int64
	LightIntensities []float64
	Normal           mgl32.Vec3
	Geo              int64
	MaterialID       int64
}

// ReadJKLFromFile will read a .jkl file and return a struct containing all necessary information
func ReadJKLFromFile(filePath string) Jkl {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return ReadJKLFromString(data)
}

// ReadJKLFromString will parse a string as a .jkl file
func ReadJKLFromString(jklString string) Jkl {
	data := jklString

	jklResult := Jkl{}

	parseVertices(data, &jklResult)
	parseTextureVertices(data, &jklResult)
	parseMaterials(data, &jklResult)
	parseColormaps(data, &jklResult)
	parseSurfaces(data, &jklResult)

	return jklResult
}

func parseSection(data string, regex string, componentRegex string, callback func(components []string)) {
	sectionRegex := regexp.MustCompile(regex)
	sectionMatch := sectionRegex.FindAllString(data, -1)

	if len(sectionMatch) == 0 {
		callback([]string{})
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(sectionMatch[0]))
	for scanner.Scan() {
		match, _ := regexp.MatchString(componentRegex, scanner.Text())
		if match != true {
			continue
		}

		components := strings.Fields(scanner.Text())

		callback(components)
	}
}

func parseVertices(data string, jklResult *Jkl) {
	parseSection(data, `(?s)World vertices.*World texture vertices`, "\\d+:.*",
		func(components []string) {
			var err error

			x, err := strconv.ParseFloat(components[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.ParseFloat(components[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			z, err := strconv.ParseFloat(components[3], 64)
			if err != nil {
				log.Fatal(err)
			}

			jklResult.Vertices = append(jklResult.Vertices, mgl32.Vec3{float32(x), float32(y), float32(z)})
		})
}

func parseTextureVertices(data string, jklResult *Jkl) {
	parseSection(data, `(?s)World texture vertices.*World adjoins`, "\\d+:.*",
		func(components []string) {
			var err error

			u, err := strconv.ParseFloat(components[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			v, err := strconv.ParseFloat(components[2], 64)
			if err != nil {
				log.Fatal(err)
			}

			jklResult.TextureVertices = append(jklResult.TextureVertices, mgl32.Vec2{float32(u), float32(v)})
		})
}

func parseMaterials(data string, jklResult *Jkl) {
	parseSection(data, `(?s)World materials.*SECTION: GEORESOURCE`, "\\d+:.*",
		func(components []string) {
			var err error

			matName := components[1]

			xTile, err := strconv.ParseFloat(components[2], 64)
			if err != nil {
				log.Fatal(err)
			}

			yTile, err := strconv.ParseFloat(components[3], 64)
			if err != nil {
				log.Fatal(err)
			}

			var matBytes []byte
			for _, file := range GobFiles {
				matBytes = LoadFileFromGOB(file, matName)
				if matBytes != nil {
					break
				}
			}

			material := ParseMatFile(matBytes)

			material.XTile = float32(xTile)
			material.YTile = float32(yTile)

			jklResult.Materials = append(jklResult.Materials, material)
		})
}

func parseColormaps(data string, jklResult *Jkl) {
	parseSection(data, `(?s)World Colormaps.*World vertices`, "\\d+:.*",
		func(components []string) {
			var cmpName string
			if len(components) == 0 {
				cmpName = "dflt.cmp"
			} else {
				cmpName = components[1]
			}

			fmt.Println(cmpName)

			var cmpBytes []byte
			for _, file := range GobFiles {
				cmpBytes = LoadFileFromGOB(file, cmpName)
				if cmpBytes != nil {
					break
				}
			}

			colorMap := ParseCmpFile(cmpBytes)

			jklResult.ColorMaps = append(jklResult.ColorMaps, colorMap)
		})
}

func parseSurfaces(data string, jklResult *Jkl) {
	parseSection(data, `(?s)World surfaces.*\#--- Surface normals ---`, "\\d+:.*",
		func(components []string) {
			surface := surface{}

			materialID, _ := strconv.ParseInt(components[1], 10, 32)
			surface.MaterialID = materialID

			geoFlag, _ := strconv.ParseInt(components[4], 10, 32)
			surface.Geo = geoFlag

			if components[5] != "3" {
				fmt.Println("light != 3", components[5])
			}

			numVertexIds, _ := strconv.ParseInt(components[9], 10, 32)
			vertexIds := components[10 : 10+numVertexIds]
			for idx, vertexIDPair := range vertexIds {
				splitVertexIDPair := strings.Split(vertexIDPair, ",")
				vertexID, _ := strconv.ParseInt(splitVertexIDPair[0], 10, 64)
				texVertexID, _ := strconv.ParseInt(splitVertexIDPair[1], 10, 64)
				surface.VertexIds = append(surface.VertexIds, vertexID)
				surface.TextureVertexIds = append(surface.TextureVertexIds, texVertexID)

				lightIntensity, _ := strconv.ParseFloat(components[10+numVertexIds:][idx], 64)
				surface.LightIntensities = append(surface.LightIntensities, lightIntensity)
			}

			jklResult.Surfaces = append(jklResult.Surfaces, surface)
		})

	parseSection(data, `(?s)\#--- Surface normals ---.*Section: SECTORS`, "\\d+:.*",
		func(components []string) {
			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			jklResult.Surfaces[surfaceID].Normal = mgl32.Vec3{float32(x), float32(y), float32(z)}
		})
}
