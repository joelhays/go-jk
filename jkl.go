package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type jkl struct {
	Vertices []vertex
	Surfaces []surface
}

type vertex struct {
	X float32
	Y float32
	Z float32
}

type surface struct {
	VertexIds []int64
	Normal    vertex
}

// ReadJKL will read a .jkl file and return a struct containing all necessary information
func ReadJKL(filePath string) jkl {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	fmt.Println("got data", len(data))

	jklResult := jkl{}

	parseVertices(data, &jklResult)
	// for _, vertex := range jklResult.Vertices {
	// 	fmt.Println(vertex.X, vertex.Y, vertex.Z)
	// }

	parseSurfaces(data, &jklResult)
	// for _, surface := range jklResult.Surfaces {
	// 	fmt.Println(surface.VertexIds)
	// }

	return jklResult
}

func parseSection(data string, regex string, callback func(components []string)) {
	sectionRegex := regexp.MustCompile(regex)
	sectionMatch := sectionRegex.FindAllString(data, -1)
	scanner := bufio.NewScanner(strings.NewReader(sectionMatch[0]))
	for scanner.Scan() {
		match, _ := regexp.MatchString("\\d+:.*", scanner.Text())
		if match != true {
			continue
		}

		components := strings.Fields(scanner.Text())

		callback(components)
	}
}

func parseVertices(data string, jklResult *jkl) {
	parseSection(data, `(?s)\#----- Vertices Subsection -----.*\#-- Texture Verts Subsection ---`,
		func(components []string) {
			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			jklResult.Vertices = append(jklResult.Vertices, vertex{X: float32(x), Y: float32(y), Z: float32(z)})
		})
}

func parseSurfaces(data string, jklResult *jkl) {
	parseSection(data, `(?s)\#----- Surfaces Subsection -----.*\#--- Surface normals ---`,
		func(components []string) {
			numVertexIds, _ := strconv.ParseInt(components[9], 10, 32)
			vertexIds := components[10 : 10+numVertexIds]

			surface := surface{}

			for _, vertexIDPair := range vertexIds {
				splitVertexIDPair := strings.Split(vertexIDPair, ",")
				vertexID, _ := strconv.ParseInt(splitVertexIDPair[0], 10, 64)
				surface.VertexIds = append(surface.VertexIds, vertexID)
			}

			jklResult.Surfaces = append(jklResult.Surfaces, surface)
		})

	parseSection(data, `(?s)\#--- Surface normals ---.*\#\#\#\#\#\# Sector information \#\#\#\#\#\#`,
		func(components []string) {
			// fmt.Println(components)

			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			jklResult.Surfaces[surfaceID].Normal = vertex{X: float32(x), Y: float32(y), Z: float32(z)}
		})
}
