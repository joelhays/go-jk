package jkparsers

import (
	"bufio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/jk/jktypes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type JklLineParser struct {
	jkl     jktypes.Jkl
	scanner *bufio.Scanner
	line    string
	done    bool
	section string
}

func NewJklLineParser() *JklLineParser {
	p := &JklLineParser{}
	p.init("")
	return p
}

func (p *JklLineParser) ParseFromFile(filePath string) jktypes.Jkl {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseFromString(data)
}

func (p *JklLineParser) ParseFromString(jklString string) jktypes.Jkl {
	p.init(jklString)

	p.scanner.Text()
	for {
		section, ok := p.advanceToNextSection()
		if !ok {
			break
		}

		switch section {
		case "jk":
		case "copyright":
		case "header":
		case "sounds":
		case "materials":
			p.parseMaterials()
		case "georesource":
			p.parseGeoResource()
		case "sectors":
		case "models":
			p.parseModels()
		case "templates":
			p.parseTemplates()
		case "things":
			p.parseThings()
		}
	}

	return p.jkl
}

func (p *JklLineParser) init(jklString string) {
	p.jkl = jktypes.Jkl{
		Model:          &jktypes.JkMesh{},
		Jk3dos:         make(map[string]jktypes.Jk3doFile),
		Jk3doTemplates: make(map[string]jktypes.Template),
		Things:         nil,
	}
	p.scanner = bufio.NewScanner(strings.NewReader(jklString))
	p.line = ""
	p.done = false
	p.section = ""
}

func (p *JklLineParser) checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (p *JklLineParser) atEndOfSection() bool {
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

func (p *JklLineParser) advanceToNextSection() (string, bool) {
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

func (p *JklLineParser) getNextLine() (string, bool) {
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

func (p *JklLineParser) getLineArgs(line string) []string {
	return p.getLineArgsWithoutPrefix(line, "")
}

func (p *JklLineParser) getLineArgsWithoutPrefix(line string, ignore string) []string {
	if len(ignore) != 0 {
		line = strings.TrimPrefix(line, ignore)
	}
	return strings.Fields(line)
}

func (p *JklLineParser) processSection(callback func(string)) {
	p.processNLines(-1, callback)
}

func (p *JklLineParser) processNLines(numToProcess int, callback func(string)) {
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

func (p *JklLineParser) parseGeoResource() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "world colormaps %d", &count); args == 1 {
			p.processNLines(count, p.parseGeoResourceWorldColormap)
		} else if args, _ = fmt.Sscanf(line, "world vertices %d", &count); args == 1 {
			p.processNLines(count, func(l string) {
				_, v := parseVec3(l)
				p.jkl.Model.Vertices = append(p.jkl.Model.Vertices, v)
			})
		} else if args, _ = fmt.Sscanf(line, "world texture vertices %d", &count); args == 1 {
			p.processNLines(count, func(l string) {
				_, v := parseVec2(l)
				p.jkl.Model.TextureVertices = append(p.jkl.Model.TextureVertices, v)
			})
		} else if args, _ = fmt.Sscanf(line, "world surfaces %d", &count); args == 1 {
			p.processNLines(count, p.parseGeoResourceWorldSurface)
			p.processNLines(count, func(l string) {
				id, v := parseVec3(l)
				p.jkl.Model.Surfaces[id].Normal = v
			})
		}
	})
}

func (p *JklLineParser) parseGeoResourceWorldColormap(line string) {
	var id int32
	var cmpName string
	n, err := fmt.Sscanf(line, "%d: %s", &id, &cmpName)
	p.checkError(err)
	if n != 2 {
		panic("Unable to get colormap information")
	}

	var colorMap jktypes.ColorMap
	fileBytes := jk.GetLoader().LoadResource(cmpName)
	if fileBytes != nil {
		colorMap = NewCmpParser().ParseFromBytes(fileBytes)
	}

	p.jkl.Model.ColorMaps = append(p.jkl.Model.ColorMaps, colorMap)
}

func (p *JklLineParser) parseGeoResourceWorldSurface(line string) {
	args := p.getLineArgs(line)

	surface := jktypes.Surface{}

	materialID, _ := strconv.ParseInt(args[1], 10, 32)
	surface.MaterialID = materialID

	geoFlag, _ := strconv.ParseInt(args[4], 10, 32)
	surface.Geo = geoFlag

	// TODO: WHAT DOES THIS VALUE MEAN?
	//if args[5] != "3" {
	//	fmt.Println("light != 3", args[5])
	//}

	numVertexIds, _ := strconv.ParseInt(args[9], 10, 32)
	vertexIds := args[10 : 10+numVertexIds]
	for idx, vertexIDPair := range vertexIds {
		splitVertexIDPair := strings.Split(vertexIDPair, ",")
		vertexID, _ := strconv.ParseInt(splitVertexIDPair[0], 10, 64)
		texVertexID, _ := strconv.ParseInt(splitVertexIDPair[1], 10, 64)
		surface.VertexIds = append(surface.VertexIds, vertexID)
		surface.TextureVertexIds = append(surface.TextureVertexIds, texVertexID)

		lightIntensity, _ := strconv.ParseFloat(args[10+numVertexIds:][idx], 64)
		surface.LightIntensities = append(surface.LightIntensities, lightIntensity)
	}

	p.jkl.Model.Surfaces = append(p.jkl.Model.Surfaces, surface)
}

func (p *JklLineParser) parseMaterials() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "world materials %d", &count); args == 1 {
			p.processNLines(count, p.parseMaterialsWorldMaterial)
		}
	})
}

func (p *JklLineParser) parseMaterialsWorldMaterial(line string) {
	var id int32
	var matName string
	var xTile float32
	var yTile float32
	n, err := fmt.Sscanf(line, "%d: %s %f %f", &id, &matName, &xTile, &yTile)
	p.checkError(err)
	if n != 4 {
		panic("Unable to get world material information")
	}

	var material jktypes.Material
	fileBytes := jk.GetLoader().LoadResource(matName)
	if fileBytes != nil {
		material = NewMatParser().ParseFromBytes(fileBytes)
	}

	material.XTile = xTile
	material.YTile = yTile

	p.jkl.Model.Materials = append(p.jkl.Model.Materials, material)
}

func (p *JklLineParser) parseModels() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "world models %d", &count); args == 1 {
			p.processNLines(count, p.parseModelsWorldModel)
		}
	})
}

func (p *JklLineParser) parseModelsWorldModel(line string) {
	var id int32
	var jk3doName string
	n, err := fmt.Sscanf(line, "%d: %s", &id, &jk3doName)
	p.checkError(err)
	if n != 2 {
		panic("Unable to get world model information")
	}

	var jk3do jktypes.Jk3doFile
	fileBytes := jk.GetLoader().LoadResource(jk3doName)
	if fileBytes != nil {
		jk3do = NewJk3doLineParser().ParseFromString(string(fileBytes))
	}

	if len(p.jkl.Model.ColorMaps) > 0 {
		jk3do.ColorMap = p.jkl.Model.ColorMaps[0]
	}
	p.jkl.Jk3dos[jk3doName] = jk3do
}

func (p *JklLineParser) parseTemplates() {
	p.processSection(func(line string) {
		var count int
		var args int
		var err error
		if args, err = fmt.Sscanf(line, "world templates %d", &count); args == 1 {
			p.processNLines(count, p.parseTemplatesWorldTemplate)
		}
		p.checkError(err)
	})
}

func (p *JklLineParser) parseTemplatesWorldTemplate(line string) {
	args := p.getLineArgs(line)

	name := args[0]
	var modelName string
	size := 1.0
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "size=") {
			size, _ = strconv.ParseFloat(strings.TrimPrefix(args[i], "size="), 32)
		}
		if strings.HasPrefix(args[i], "model3d=") {
			modelName = strings.TrimPrefix(args[i], "model3d=")
		}
	}

	if modelName != "" {
		tmp := jktypes.Template{}
		tmp.Name = name
		tmp.Jk3doName = modelName
		tmp.Size = size

		p.jkl.Jk3doTemplates[tmp.Name] = tmp
	}
}

func (p *JklLineParser) parseThings() {
	p.processSection(func(line string) {
		var count int
		var args int
		if args, _ = fmt.Sscanf(line, "world things %d", &count); args == 1 {
			p.processNLines(count, p.parseThingsWorldThing)
		}
	})
}

func (p *JklLineParser) parseThingsWorldThing(line string) {
	args := p.getLineArgs(line)

	templateName := args[1]

	x, _ := strconv.ParseFloat(args[3], 64)
	y, _ := strconv.ParseFloat(args[4], 64)
	z, _ := strconv.ParseFloat(args[5], 64)

	pitch, _ := strconv.ParseFloat(args[6], 64)
	yaw, _ := strconv.ParseFloat(args[7], 64)
	Roll, _ := strconv.ParseFloat(args[8], 64)

	t := jktypes.Thing{}
	t.TemplateName = templateName
	t.Position = mgl32.Vec3{float32(x), float32(y), float32(z)}
	t.Pitch = pitch
	t.Yaw = yaw
	t.Roll = Roll

	p.jkl.Things = append(p.jkl.Things, t)
}
