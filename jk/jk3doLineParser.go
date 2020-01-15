package jk

import (
	"bufio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Jk3doLineParser struct {
	jk3do   Jk3doFile
	scanner *bufio.Scanner
	line    string
	done    bool
	section string
}

func NewJk3doLineParser() *Jk3doLineParser {
	return &Jk3doLineParser{
		jk3do: Jk3doFile{},
	}
}

func (p *Jk3doLineParser) ParseFromFile(filePath string) Jk3doFile {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseFromString(data)
}

func (p *Jk3doLineParser) ParseFromString(objString string) Jk3doFile {
	p.jk3do = Jk3doFile{}
	p.scanner = bufio.NewScanner(strings.NewReader(objString))
	p.line = ""
	p.done = false

	p.getNextLine() // SECTION: HEADER
	p.getNextLine() // 3DO 2.1

	p.getNextLine() // SECTION: MODELRESOURCE
	p.getNextLine() // MATERIALS %d
	p.parseMaterials()

	p.getNextLine() // SECTION: GEOMETRYDEF
	p.getNextLine() // RADIUS %f
	p.getNextLine() // INSERT OFFSET %f %f %f
	p.getNextLine() // GEOSETS %d
	p.parseGeoSets()

	p.getNextLine() // SECTION: HIERARCHYDEF
	p.getNextLine() // HIERARCHY NODES %d
	p.parseHierarchyNodes()

	cmpName := "dflt.cmp"
	p.jk3do.ColorMap = GetLoader().LoadCMP(cmpName)

	return p.jk3do
}

func (p *Jk3doLineParser) checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Jk3doLineParser) getNextLine() bool {
	for {
		ok := p.scanner.Scan()
		if !ok {
			p.done = true
			break
		}
		line := p.scanner.Text()
		line = strings.TrimSpace(line)
		line = strings.ToLower(line)
		p.line = line

		if len(line) == 0 {
			continue //blank line
		}
		if strings.HasPrefix(line, "#") {
			continue //comment
		}

		return true
	}
	return false
}

func (p *Jk3doLineParser) parseMaterials() {
	var count int
	_, err := fmt.Sscanf(p.line, "materials %d", &count)
	p.checkError(err)

	p.jk3do.Materials = make([]Material, count)

	for i := 0; i < count; i++ {
		p.getNextLine()

		var id int32
		var matName string
		_, err := fmt.Sscanf(p.line, "%d:%s", &id, &matName)
		p.checkError(err)

		material := GetLoader().LoadMAT(matName)
		material.XTile = 1.0
		material.YTile = 1.0

		p.jk3do.Materials[i] = material
	}
}

func (p *Jk3doLineParser) parseGeoSets() {
	var count int
	_, err := fmt.Sscanf(p.line, "geosets %d", &count)
	p.checkError(err)

	p.jk3do.GeoSets = make([]GeoSet, count)

	for i := 0; i < count; i++ {
		geoset := &p.jk3do.GeoSets[i]

		p.getNextLine() // GEOSET %d
		p.getNextLine() // MESHES %d
		p.parseMeshes(geoset)
	}
}

func (p *Jk3doLineParser) parseMeshes(geoset *GeoSet) {
	var count int
	_, err := fmt.Sscanf(p.line, "meshes %d", &count)
	p.checkError(err)

	geoset.Meshes = make([]Mesh, count)

	for i := 0; i < count; i++ {
		mesh := &geoset.Meshes[i]

		p.getNextLine() // MESH %d
		p.getNextLine() // NAME %s
		p.getNextLine() // RADIUS %f
		p.getNextLine() // GEOMETRYMODE %d
		p.getNextLine() // LIGHTINGMODE %d
		p.getNextLine() // TEXTUREMODE %d
		p.getNextLine() // VERTICES %d
		p.parseVertices(mesh)

		p.getNextLine() // TEXTURE VERTICES %d
		p.parseTextureVertices(mesh)

		p.getNextLine() // VERTEX NORMALS
		p.parseVertexNormals(mesh)

		p.getNextLine() // FACES %d
		p.parseFaces(mesh)

		p.getNextLine() // FACE NORMALS
		p.parseFaceNormals(mesh)
	}
}

func (p *Jk3doLineParser) parseVertices(mesh *Mesh) {
	var count int
	_, err := fmt.Sscanf(p.line, "vertices %d", &count)
	p.checkError(err)

	mesh.Vertices = make([]mgl32.Vec3, count)

	for i := 0; i < count; i++ {
		p.getNextLine()

		_, v := parseVec3(p.line)
		mesh.Vertices[i] = v
	}
}

func (p *Jk3doLineParser) parseTextureVertices(mesh *Mesh) {
	var count int
	_, err := fmt.Sscanf(p.line, "texture vertices %d", &count)
	p.checkError(err)

	mesh.TextureVertices = make([]mgl32.Vec2, count)

	for i := 0; i < count; i++ {
		p.getNextLine()

		_, v := parseVec2(p.line)
		mesh.TextureVertices[i] = v
	}
}

func (p *Jk3doLineParser) parseVertexNormals(mesh *Mesh) {
	numVerts := len(mesh.Vertices)
	mesh.VertexNormals = make([]mgl32.Vec3, numVerts)

	for i := 0; i < numVerts; i++ {
		p.getNextLine()

		_, v := parseVec3(p.line)
		mesh.VertexNormals[i] = v
	}
}

func (p *Jk3doLineParser) parseFaces(mesh *Mesh) {
	var count int
	_, err := fmt.Sscanf(p.line, "faces %d", &count)
	p.checkError(err)

	mesh.Faces = make([]Face, count)

	for i := 0; i < count; i++ {
		p.getNextLine()

		args := strings.Fields(strings.Replace(p.line, ",", " ", -1))

		surface := Face{}

		materialID, _ := strconv.ParseInt(args[1], 10, 32)
		surface.MaterialID = materialID

		geoFlag, _ := strconv.ParseInt(args[3], 10, 32)
		surface.GeometryMode = geoFlag

		//TODO: WHAT DOES THIS VALUE MEAN?
		//if components[4] != "3" {
		//	fmt.Println("light != 3", components[5])
		//}

		numVertexIds, _ := strconv.ParseInt(args[7], 10, 32)
		vertexIds := args[8 : 8+(numVertexIds*2)]
		for v := 0; v < int(numVertexIds*2); v += 2 {
			vertexID, _ := strconv.ParseInt(strings.TrimRight(vertexIds[v], ","), 10, 32)
			texVertexID, _ := strconv.ParseInt(vertexIds[v+1], 10, 32)
			surface.VertexIds = append(surface.VertexIds, vertexID)
			surface.TextureVertexIds = append(surface.TextureVertexIds, texVertexID)

			lightIntensity := 1.0
			surface.LightIntensities = append(surface.LightIntensities, lightIntensity)
		}

		mesh.Faces[i] = surface
	}
}

func (p *Jk3doLineParser) parseFaceNormals(mesh *Mesh) {
	numFaces := len(mesh.Faces)
	mesh.FaceNormals = make([]mgl32.Vec3, numFaces)

	for i := 0; i < numFaces; i++ {
		p.getNextLine()

		_, v := parseVec3(p.line)
		mesh.FaceNormals[i] = v
	}
}

func (p *Jk3doLineParser) parseHierarchyNodes() {
	var count int
	_, err := fmt.Sscanf(p.line, "hierarchy nodes %d", &count)
	p.checkError(err)

	p.jk3do.Hierarchy = make([]HierarchyDef, count)

	for i := 0; i < count; i++ {
		p.getNextLine()

		args := strings.Fields(p.line)

		meshID, _ := strconv.ParseInt(args[3], 10, 32)
		parentID, _ := strconv.ParseInt(args[4], 10, 32)
		childID, _ := strconv.ParseInt(args[5], 10, 32)
		siblingID, _ := strconv.ParseInt(args[6], 10, 32)
		numChildren, _ := strconv.ParseInt(args[7], 10, 32)

		x, _ := strconv.ParseFloat(args[8], 32)
		y, _ := strconv.ParseFloat(args[9], 32)
		z, _ := strconv.ParseFloat(args[10], 32)

		pitch, _ := strconv.ParseFloat(args[11], 32)
		yaw, _ := strconv.ParseFloat(args[12], 32)
		roll, _ := strconv.ParseFloat(args[13], 32)

		pivotX, _ := strconv.ParseFloat(args[14], 32)
		pivotY, _ := strconv.ParseFloat(args[15], 32)
		pivotZ, _ := strconv.ParseFloat(args[16], 32)

		nodeName := args[17]

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

		p.jk3do.Hierarchy[i] = def
	}
}
