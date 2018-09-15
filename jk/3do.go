package jk

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type Jk3do struct {
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	VertexNormals   []mgl32.Vec3
	Surfaces        []surface
	Materials       []Material
	ColorMaps       []ColorMap
}

func Parse3doFromString(data string) Jk3do {
	result := Jk3do{}

	parse3doVertices(data, &result)
	parse3doTextureVertices(data, &result)
	parse3doMaterials(data, &result)
	parse3doColormaps(data, &result)
	parse3doSurfaces(data, &result)

	return result
}

func parse3doSection(data string, regex string, componentRegex string, callback func(components []string)) {
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

func parse3doVertices(data string, obj *Jk3do) {
	parse3doSection(data, `(?s)VERTICES.*TEXTURE VERTICES`, "\\d+:.*",
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

			obj.Vertices = append(obj.Vertices, mgl32.Vec3{float32(x), float32(y), float32(z)})
		})
}

func parse3doTextureVertices(data string, obj *Jk3do) {
	parse3doSection(data, `(?s)TEXTURE VERTICES.*VERTEX NORMALS`, "\\d+:.*",
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

			obj.TextureVertices = append(obj.TextureVertices, mgl32.Vec2{float32(u), float32(v)})
		})
}

func parse3doMaterials(data string, obj *Jk3do) {
	parse3doSection(data, `(?s)MATERIALS.*SECTION: GEOMETRYDEF`, "\\d+:.*",
		func(components []string) {
			matName := components[1]

			var matBytes []byte
			for _, file := range GobFiles {
				matBytes = LoadFileFromGOB(file, matName)
				if matBytes != nil {
					break
				}
			}

			material := ParseMatFile(matBytes)

			material.XTile = 1.0
			material.YTile = 1.0

			obj.Materials = append(obj.Materials, material)
		})
}

func parse3doSurfaces(data string, obj *Jk3do) {
	parse3doSection(data, `(?s)FACES.*FACE NORMALS`, "\\d+:.*",
		func(components []string) {
			surface := surface{}

			materialID, _ := strconv.ParseInt(components[1], 10, 32)
			surface.MaterialID = materialID

			geoFlag, _ := strconv.ParseInt(components[3], 10, 32)
			surface.Geo = geoFlag

			if components[4] != "3" {
				fmt.Println("light != 3", components[5])
			}

			numVertexIds, _ := strconv.ParseInt(components[7], 10, 32)
			vertexIds := components[8 : 8+(numVertexIds*2)]
			for i := 0; i < int(numVertexIds*2); i += 2 {
				vertexID, _ := strconv.ParseInt(strings.TrimRight(vertexIds[i], ","), 10, 64)
				texVertexID, _ := strconv.ParseInt(vertexIds[i+1], 10, 64)
				surface.VertexIds = append(surface.VertexIds, vertexID)
				surface.TextureVertexIds = append(surface.TextureVertexIds, texVertexID)

				lightIntensity := 1.0
				surface.LightIntensities = append(surface.LightIntensities, lightIntensity)
			}
			obj.Surfaces = append(obj.Surfaces, surface)

			fmt.Println(surface.VertexIds)
			fmt.Println(surface.TextureVertexIds)
		})

	parseSection(data, `(?s)FACE NORMALS.*SECTION: HIERARCHYDEF`, "\\d+:.*",
		func(components []string) {
			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			obj.Surfaces[surfaceID].Normal = mgl32.Vec3{float32(x), float32(y), float32(z)}
		})
}

func parse3doColormaps(data string, obj *Jk3do) {
	var cmpName string
	cmpName = "dflt.cmp"

	var cmpBytes []byte
	for _, file := range GobFiles {
		cmpBytes = LoadFileFromGOB(file, cmpName)
		if cmpBytes != nil {
			break
		}
	}

	colorMap := ParseCmpFile(cmpBytes)

	obj.ColorMaps = append(obj.ColorMaps, colorMap)
}
