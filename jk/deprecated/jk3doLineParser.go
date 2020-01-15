package deprecated

import (
	"bufio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/jk"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Jk3doLineParser struct {
	jk3do   jk.Jk3doFile
	scanner *bufio.Scanner
	line    string
	done    bool
	section string
}

func NewJk3doLineParser() *Jk3doLineParser {
	return &Jk3doLineParser{
		jk3do: jk.Jk3doFile{},
	}
}

func (p *Jk3doLineParser) Parse3doFromFile(filePath string) jk.Jk3doFile {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.Parse3doFromString(data)
}

func (p *Jk3doLineParser) Parse3doFromString(objString string) jk.Jk3doFile {
	p.jk3do = jk.Jk3doFile{}
	p.scanner = bufio.NewScanner(strings.NewReader(objString))
	p.line = ""
	p.done = false

	p.scanner.Text()
	for {
		section, ok := p.advanceToNextSection()
		if !ok {
			break
		}

		switch section {
		case "header":
		case "modelresource":
			p.parseModelResource()
		case "geometrydef":
			p.parseGeometryDef()
		case "hierarchydef":
			p.parseHierarchyDef()
		}
	}

	cmpName := "dflt.cmp"
	p.jk3do.ColorMap = jk.GetLoader().LoadCMP(cmpName)

	return p.jk3do
}

func (p *Jk3doLineParser) checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Jk3doLineParser) atEndOfSection() bool {
	if strings.HasPrefix(p.line, "section: ") {
		section := strings.TrimPrefix(p.line, "section: ")
		if section != p.section {
			return true
		}
	}

	if p.line == "end" {
		return true
	}

	return false
}

func (p *Jk3doLineParser) advanceToNextSection() (string, bool) {
	if p.done {
		p.section = ""
		return "", false
	}

	if p.atEndOfSection() {
		if strings.HasPrefix(p.line, "section: ") {
			section := strings.TrimPrefix(p.line, "section: ")
			p.section = section
			return section, true
		}
	}

	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		if strings.HasPrefix(line, "section: ") {
			section := strings.TrimPrefix(line, "section: ")
			p.section = section
			return section, true
		}
	}

	p.section = ""
	return "", false
}

func (p *Jk3doLineParser) getNextLine() (string, bool) {
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

		return line, true
	}
	return "", false
}

func (p *Jk3doLineParser) getLineArgs(line string) []string {
	return p.getLineArgsWithoutPrefix(line, "")
}

func (p *Jk3doLineParser) getLineArgsWithoutPrefix(line string, ignore string) []string {
	if len(ignore) != 0 {
		line = strings.TrimPrefix(line, ignore)
	}
	return strings.Fields(line)
}

func (p *Jk3doLineParser) processSection(callback func(string)) {
	p.processNLines(-1, callback)
}

func (p *Jk3doLineParser) processNLines(numToProcess int, callback func(string)) {
	if numToProcess == 0 {
		return
	}

	numProcessed := 0
	for {
		if numToProcess > 0 && numProcessed == numToProcess {
			break
		}

		line, ok := p.getNextLine()
		if !ok {
			break
		}

		if p.atEndOfSection() {
			break
		}

		callback(line)

		numProcessed++
	}
}

func (p *Jk3doLineParser) parseModelResource() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "materials %d", &count); args == 1 {
			p.processNLines(count, func(l string) {
				var id int32
				var matName string
				n, err := fmt.Sscanf(l, "%d: %s", &id, &matName)
				p.checkError(err)
				if n != 2 {
					panic("Unable to get material information")
				}

				material := jk.GetLoader().LoadMAT(matName)
				material.XTile = 1.0
				material.YTile = 1.0

				p.jk3do.Materials = append(p.jk3do.Materials, material)
			})
		}
	})
}

func (p *Jk3doLineParser) parseGeometryDef() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "geosets %d", &count); args == 1 {
			for g := 0; g < count; g++ {

				geoset := jk.GeoSet{}

				p.getNextLine() // GEOSET %d
				p.getNextLine() // MESHES %d

				meshCount := 0
				args, _ = fmt.Sscanf(p.line, "meshes %d", &meshCount)

				geoset.Meshes = make([]jk.Mesh, meshCount)
				p.jk3do.GeoSets = append(p.jk3do.GeoSets, geoset)

				for m := 0; m < meshCount; m++ {
					mesh := &geoset.Meshes[m]

					p.getNextLine() // MESH %d

					p.getNextLine() // NAME %s
					p.getNextLine() // RADIUS %f
					p.getNextLine() // GEOMETRYMODE %d
					p.getNextLine() // LIGHTINGMODE %d
					p.getNextLine() // TEXTUREMODE %d

					p.getNextLine() // VERTICES %d
					vtxCount := 0
					args, _ = fmt.Sscanf(p.line, "vertices %d", &vtxCount)
					p.processNLines(vtxCount, func(l string) {
						_, v := jk.parseVec3(l)
						mesh.Vertices = append(mesh.Vertices, v)
					})

					p.getNextLine() // TEXTURE VERTICES %d
					texVtxCount := 0
					args, _ = fmt.Sscanf(p.line, "texture vertices %d", &texVtxCount)
					p.processNLines(texVtxCount, func(l string) {
						_, v := jk.parseVec2(l)
						mesh.TextureVertices = append(mesh.TextureVertices, v)
					})

					p.getNextLine() // VERTEX NORMALS
					p.processNLines(vtxCount, func(l string) {
						_, v := jk.parseVec3(l)
						mesh.VertexNormals = append(mesh.VertexNormals, v)
					})

					p.getNextLine() // FACES %d
					faceCount := 0
					args, _ = fmt.Sscanf(p.line, "faces %d", &faceCount)
					p.processNLines(faceCount, func(l string) {
						args := strings.Fields(strings.Replace(l, ",", " ", -1))

						surface := jk.Face{}

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
						mesh.Faces = append(mesh.Faces, surface)
					})

					p.getNextLine() // FACE NORMALS
					p.processNLines(faceCount, func(l string) {
						_, v := jk.parseVec3(l)
						mesh.FaceNormals = append(mesh.FaceNormals, v)
					})
				}
			}
		}
	})
}

func (p *Jk3doLineParser) parseHierarchyDef() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "hierarchy nodes %d", &count); args == 1 {
			p.processNLines(count, func(l string) {
				args := p.getLineArgs(l)

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

				def := jk.HierarchyDef{
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

				p.jk3do.Hierarchy = append(p.jk3do.Hierarchy, def)
			})
		}
	})
}
