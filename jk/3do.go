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

type Jk3doFile struct {
	Materials []Material
	GeoSets   []GeoSet
	Hierarchy []HierarchyDef
	ColorMap  ColorMap
}

type GeoSet struct {
	Meshes []Mesh
}

type Mesh struct {
	GeometryMode    int64
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	VertexNormals   []mgl32.Vec3
	Faces           []Face
	FaceNormals     []mgl32.Vec3
}

type Face struct {
	VertexIds        []int64
	TextureVertexIds []int64
	LightIntensities []float64
	GeometryMode     int64
	MaterialID       int64
}

type HierarchyDef struct {
	MeshID      int64
	ParentID    int64
	ChildID     int64
	SiblingID   int64
	NumChildren int64
	Position    mgl32.Vec3
	Pitch       float64
	Yaw         float64
	Roll        float64
	Pivot       mgl32.Vec3
	NodeName    string
}

func Parse3doFile(data string) Jk3doFile {
	result := Jk3doFile{}

	parse3doFileMaterials(data, &result)
	parse3doFileHierarchy(data, &result)

	geosetRegex := regexp.MustCompile(`(?s)GEOSET\s\d`)
	geosetMatches := geosetRegex.Split(data, -1)[1:]

	meshRegex := regexp.MustCompile(`(?s)MESH\s\d`)

	for _, geosetMatch := range geosetMatches {
		meshMatches := meshRegex.Split(geosetMatch, -1)[1:]

		geoset := GeoSet{}
		geoset.Meshes = make([]Mesh, len(meshMatches))

		result.GeoSets = append(result.GeoSets, geoset)

		var meshwg sync.WaitGroup
		meshwg.Add(len(meshMatches))

		for i := 0; i < len(meshMatches); i++ {
			go func(idx int) {
				defer meshwg.Done()

				meshData := meshMatches[idx]

				mesh := &geoset.Meshes[idx]

				parse3doFileVertices(meshData, mesh)
				parse3doFileTextureVertices(meshData, mesh)
				//TODO: PARSE VERTEX NORMALS
				parse3doFileSurfaces(meshData, mesh)
			}(i)
		}

		meshwg.Wait()
	}

	cmpName := "dflt.cmp"
	var cmpBytes []byte
	for _, file := range GobFiles {
		cmpBytes = LoadFileFromGOB(file, cmpName)
		if cmpBytes != nil {
			break
		}
	}
	result.ColorMap = ParseCmpFile(cmpBytes)

	return result
}

func parse3doFileMaterials(data string, obj *Jk3doFile) {
	parse3doFileSection(data, `(?s)MATERIALS.*?SECTION: GEOMETRYDEF`, "\\d+:.*",
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

func parse3doFileHierarchy(data string, obj *Jk3doFile) {
	parse3doFileSection(data, `(?s)SECTION: HIERARCHYDEF.*`, "\\d+:.*",
		func(components []string) {

			// id, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)
			// if id == 0 {
			// 	return
			// }

			meshID, _ := strconv.ParseInt(components[3], 10, 32)
			parentID, _ := strconv.ParseInt(components[4], 10, 32)
			childID, _ := strconv.ParseInt(components[5], 10, 32)
			siblingID, _ := strconv.ParseInt(components[6], 10, 32)
			numChildren, _ := strconv.ParseInt(components[7], 10, 32)

			x, _ := strconv.ParseFloat(components[8], 32)
			y, _ := strconv.ParseFloat(components[9], 32)
			z, _ := strconv.ParseFloat(components[10], 32)

			pitch, _ := strconv.ParseFloat(components[11], 32)
			yaw, _ := strconv.ParseFloat(components[12], 32)
			roll, _ := strconv.ParseFloat(components[13], 32)

			pivotX, _ := strconv.ParseFloat(components[14], 32)
			pivotY, _ := strconv.ParseFloat(components[15], 32)
			pivotZ, _ := strconv.ParseFloat(components[16], 32)

			nodeName := components[17]

			def := HierarchyDef{
				MeshID:      meshID,
				ParentID:    parentID,
				ChildID:     childID,
				SiblingID:   siblingID,
				NumChildren: numChildren,
				Position:    mgl32.Vec3{float32(x), float32(y), float32(z)},
				Pitch:       pitch,
				Yaw:         yaw,
				Roll:        roll,
				Pivot:       mgl32.Vec3{float32(pivotX), float32(pivotY), float32(pivotZ)},
				NodeName:    nodeName,
			}

			obj.Hierarchy = append(obj.Hierarchy, def)
		})
}

func parse3doFileSection(data string, regex string, componentRegex string, callback func(components []string)) {
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

func parse3doFileVertices(data string, obj *Mesh) {
	parse3doFileSection(data, `(?s)VERTICES.*?TEXTURE VERTICES`, "\\d+:.*",
		func(components []string) {
			var err error

			x, err := strconv.ParseFloat(components[1], 32)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.ParseFloat(components[2], 32)
			if err != nil {
				log.Fatal(err)
			}
			z, err := strconv.ParseFloat(components[3], 32)
			if err != nil {
				log.Fatal(err)
			}

			obj.Vertices = append(obj.Vertices, mgl32.Vec3{float32(x), float32(y), float32(z)})
		})
}

func parse3doFileTextureVertices(data string, obj *Mesh) {
	parse3doFileSection(data, `(?s)TEXTURE VERTICES.*?VERTEX NORMALS`, "\\d+:.*",
		func(components []string) {
			var err error

			u, err := strconv.ParseFloat(components[1], 32)
			if err != nil {
				log.Fatal(err)
			}
			v, err := strconv.ParseFloat(components[2], 32)
			if err != nil {
				log.Fatal(err)
			}

			obj.TextureVertices = append(obj.TextureVertices, mgl32.Vec2{float32(u), float32(v)})
		})
}

func parse3doFileSurfaces(data string, obj *Mesh) {
	parse3doFileSection(data, `(?s)FACES.*?FACE NORMALS`, "\\d+:.*",
		func(components []string) {
			surface := Face{}

			materialID, _ := strconv.ParseInt(components[1], 10, 32)
			surface.MaterialID = materialID

			geoFlag, _ := strconv.ParseInt(components[3], 10, 32)
			surface.GeometryMode = geoFlag

			//TODO: WHAT DOES THIS VALUE MEAN?
			// if components[4] != "3" {
			// 	fmt.Println("light != 3", components[5])
			// }

			numVertexIds, _ := strconv.ParseInt(components[7], 10, 32)
			vertexIds := components[8 : 8+(numVertexIds*2)]
			for i := 0; i < int(numVertexIds*2); i += 2 {
				vertexID, _ := strconv.ParseInt(strings.TrimRight(vertexIds[i], ","), 10, 32)
				texVertexID, _ := strconv.ParseInt(vertexIds[i+1], 10, 32)
				surface.VertexIds = append(surface.VertexIds, vertexID)
				surface.TextureVertexIds = append(surface.TextureVertexIds, texVertexID)

				lightIntensity := 1.0
				surface.LightIntensities = append(surface.LightIntensities, lightIntensity)
			}
			obj.Faces = append(obj.Faces, surface)
		})

	obj.FaceNormals = make([]mgl32.Vec3, len(obj.Faces))

	parse3doFileSection(data, `(?s)FACE NORMALS.*?(SECTION: HIERARCHYDEF|Mesh definition|Geometry Set definition)`, "\\d+:.*",
		func(components []string) {
			if len(obj.Faces) == 0 {
				return
			}
			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 32)
			y, _ := strconv.ParseFloat(components[2], 32)
			z, _ := strconv.ParseFloat(components[3], 32)

			obj.FaceNormals[surfaceID] = mgl32.Vec3{float32(x), float32(y), float32(z)}
		})
}
