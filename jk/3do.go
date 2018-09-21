//TODO: SUPPORT MULTI-MESH MODELS

package jk

import (
	"bufio"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type Jk3do struct {
	Meshes         []JkMesh
	MeshTransforms map[string]JkMeshTransform
}

type JkMesh struct {
	Name            string
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	VertexNormals   []mgl32.Vec3
	Surfaces        []surface
	Materials       []Material
	ColorMaps       []ColorMap
}

type JkMeshTransform struct {
	Offset mgl32.Vec3
	Pitch  float32
	Yaw    float32
	Roll   float32
	Pivot  mgl32.Vec3
}

func Parse3doFromString(data string) Jk3do {
	result := Jk3do{}

	// meshRegex := regexp.MustCompile(`(?s)MESH\s\d.*?(MESH|HIERARCHYDEF)`)
	// meshMatch := meshRegex.FindAllString(data, -1)

	geosetRegex := regexp.MustCompile(`(?s)GEOSET\s\d`)
	geosetMatch := geosetRegex.Split(data, -1)

	meshRegex := regexp.MustCompile(`(?s)MESH\s\d`)
	meshMatch := meshRegex.Split(geosetMatch[1], -1)

	result.Meshes = make([]JkMesh, len(meshMatch)-1)

	result.MeshTransforms = make(map[string]JkMeshTransform)

	var meshwg sync.WaitGroup
	meshwg.Add(len(meshMatch) - 1)

	for i := 1; i < len(meshMatch); i++ {
		go func(idx int) {
			defer meshwg.Done()

			meshData := meshMatch[idx]

			nameRegex := regexp.MustCompile(`(?s)NAME\s.*?(\r|\n)`)
			nameMatch := nameRegex.FindString(meshData)

			name := strings.TrimLeft(nameMatch, "NAME")
			name = strings.TrimLeft(name, " ")
			name = strings.TrimRight(name, " ")
			name = strings.TrimRight(name, "\r")
			name = strings.TrimRight(name, "\n")

			result.Meshes[idx-1].Name = name

			parse3doVertices(meshData, &result.Meshes[idx-1])
			parse3doTextureVertices(meshData, &result.Meshes[idx-1])
			parse3doMaterials(data, &result.Meshes[idx-1])
			parse3doColormaps(meshData, &result.Meshes[idx-1])
			parse3doSurfaces(meshData, &result.Meshes[idx-1])
		}(i)
	}

	meshwg.Wait()

	parse3doHierarchy(data, &result)

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
		text := strings.Replace(scanner.Text(), ",", " ", -1)

		space := regexp.MustCompile(`\s+`)
		text = space.ReplaceAllString(text, "|")
		text = strings.TrimLeft(text, "|")
		text = strings.TrimRight(text, "|")

		components := strings.Split(text, "|")

		if strings.Contains(components[0], "#") {
			continue
		}

		callback(components)
	}
}

func parse3doVertices(data string, obj *JkMesh) {
	parse3doSection(data, `(?s)VERTICES.*?TEXTURE VERTICES`, "\\d+:.*",
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

func parse3doTextureVertices(data string, obj *JkMesh) {
	parse3doSection(data, `(?s)TEXTURE VERTICES.*?VERTEX NORMALS`, "\\d+:.*",
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

func parse3doMaterials(data string, obj *JkMesh) {
	parse3doSection(data, `(?s)MATERIALS.*?SECTION: GEOMETRYDEF`, "\\d+:.*",
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

func parse3doSurfaces(data string, obj *JkMesh) {
	parse3doSection(data, `(?s)FACES.*?FACE NORMALS`, "\\d+:.*",
		func(components []string) {
			surface := surface{}

			materialID, _ := strconv.ParseInt(components[1], 10, 32)
			surface.MaterialID = materialID

			geoFlag, _ := strconv.ParseInt(components[3], 10, 32)
			surface.Geo = geoFlag

			//TODO: WHAT DOES THIS VALUE MEAN?
			// if components[4] != "3" {
			// 	fmt.Println("light != 3", components[5])
			// }

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
		})

	parse3doSection(data, `(?s)FACE NORMALS.*?(SECTION: HIERARCHYDEF|Mesh definition|Geometry Set definition)`, "\\d+:.*",
		func(components []string) {
			if len(obj.Surfaces) == 0 {
				return
			}
			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			obj.Surfaces[surfaceID].Normal = mgl32.Vec3{float32(x), float32(y), float32(z)}
		})
}

func parse3doHierarchy(data string, obj *Jk3do) {
	parse3doSection(data, `(?s)SECTION: HIERARCHYDEF.*`, "\\d+:.*",
		func(components []string) {
			// fmt.Println(components)

			id, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)
			if id == 0 {
				return
			}

			// meshID, _ := strconv.ParseInt(components[3], 10, 32)
			// parentID, _ := strconv.ParseInt(components[4], 10, 32)

			x, _ := strconv.ParseFloat(components[8], 64)
			y, _ := strconv.ParseFloat(components[9], 64)
			z, _ := strconv.ParseFloat(components[10], 64)

			pitch, _ := strconv.ParseFloat(components[11], 64)
			yaw, _ := strconv.ParseFloat(components[12], 64)
			roll, _ := strconv.ParseFloat(components[13], 64)

			pivotX, _ := strconv.ParseFloat(components[14], 64)
			pivotY, _ := strconv.ParseFloat(components[15], 64)
			pivotZ, _ := strconv.ParseFloat(components[16], 64)

			meshName := components[17]

			obj.MeshTransforms[meshName] = JkMeshTransform{
				Offset: mgl32.Vec3{float32(x), float32(y), float32(z)},
				Pitch:  float32(pitch),
				Yaw:    float32(yaw),
				Roll:   float32(roll),
				Pivot:  mgl32.Vec3{float32(pivotX), float32(pivotY), float32(pivotZ)},
			}
		})
}

func parse3doColormaps(data string, obj *JkMesh) {
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
