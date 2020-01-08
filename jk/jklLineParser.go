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

type JklLineParser struct {
	jkl     Jkl
	scanner *bufio.Scanner
	line    string
	done    bool
	section string
}

func NewJklLineParser() *JklLineParser {
	return &JklLineParser{
		jkl: Jkl{
			Model:          &JkMesh{},
			Jk3dos:         make(map[string]Jk3doFile),
			Jk3doTemplates: make(map[string]Template),
			Things:         nil,
		},
	}
}

func (p *JklLineParser) ParseJKLFromFile(filePath string) Jkl {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseJKLFromString(data)
}

func (p *JklLineParser) ParseJKLFromString(jklString string) Jkl {
	p.jkl = Jkl{
		Model:          &JkMesh{},
		Jk3dos:         make(map[string]Jk3doFile),
		Jk3doTemplates: make(map[string]Template),
		Things:         nil,
	}
	p.scanner = bufio.NewScanner(strings.NewReader(jklString))
	p.line = ""
	p.done = false

	p.scanner.Text()
	for {
		section, ok := p.advanceToNextSection()
		if !ok {
			break
		}

		switch section {
		case "JK":
			p.getNextLine()
			continue
		case "COPYRIGHT":
			p.getNextLine()
			continue
		case "HEADER":
			p.getNextLine()
			continue
		case "SOUNDS":
			p.getNextLine()
			continue
		case "MATERIALS":
			p.parseMaterials()
			continue
		case "GEORESOURCE":
			p.parseGeoResource()
			continue
		case "SECTORS":
			p.getNextLine()
			continue
		case "MODELS":
			p.parseModels()
			continue
		case "TEMPLATES":
			p.parseTemplates()
			continue
		case "THINGS":
			p.parseThings()
			continue
		default:
			p.getNextLine()
			continue
		}
	}

	return p.jkl
}

func (p *JklLineParser) atEndOfSection() bool {
	currentLine := strings.ToUpper(p.line)
	if strings.HasPrefix(currentLine, "SECTION: ") {
		p.section = strings.TrimPrefix(strings.ToUpper(p.line), "SECTION: ")
		return true
	}

	if currentLine == "END" {
		p.section = ""
		return true
	}

	return false
}

func (p *JklLineParser) advanceToNextSection() (string, bool) {
	if p.done {
		return "", false
	}

	if p.atEndOfSection() {
		currentLine := strings.ToUpper(p.line)
		if strings.HasPrefix(currentLine, "SECTION: ") {
			section := strings.TrimPrefix(strings.ToUpper(p.line), "SECTION: ")
			return section, true
		}
	}

	//////todo: handle current section better...
	//currentLine := strings.ToUpper(p.line)
	//if strings.HasPrefix(currentLine, "SECTION: ") {
	//	section := strings.TrimPrefix(currentLine, "SECTION: ")
	//	return section, true
	//}

	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)
		if strings.HasPrefix(line, "SECTION: ") {
			section := strings.TrimPrefix(line, "SECTION: ")
			return section, true
		}
	}

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

func (p *JklLineParser) parseGeoResource() {
	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)

		if strings.HasPrefix(line, "SECTION: ") {
			break
		}

		if strings.HasPrefix(line, "WORLD COLORMAPS") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD COLORMAPS")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseGeoResourceWorldColormaps(count)
		} else if strings.HasPrefix(line, "WORLD VERTICES") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD VERTICES")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseGeoResourceWorldVertices(count)
		} else if strings.HasPrefix(line, "WORLD TEXTURE VERTICES") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD TEXTURE VERTICES")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseGeoResourceWorldTextureVertices(count)
		} else if strings.HasPrefix(line, "WORLD SURFACES") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD SURFACES")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseGeoResourceWorldSurfaces(count)
		}
	}
}

func (p *JklLineParser) parseGeoResourceWorldColormaps(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			continue
		}
		var id int32
		var cmpName string
		n, err := fmt.Sscanf(line, "%d: %s", &id, &cmpName)
		if err != nil {
			panic(err)
		}
		if n != 2 {
			panic("Unable to get colormap information")
		}

		colorMap := GetLoader().LoadCMP(cmpName)

		p.jkl.Model.ColorMaps = append(p.jkl.Model.ColorMaps, colorMap)
	}
}

func (p *JklLineParser) parseGeoResourceWorldVertices(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			continue
		}
		var id int32
		v := mgl32.Vec3{}
		n, err := fmt.Sscanf(line, "%d: %f %f %f", &id, &v[0], &v[1], &v[2])
		if err != nil {
			panic(err)
		}
		if n != 4 {
			panic("Unable to get vertex information")
		}

		p.jkl.Model.Vertices = append(p.jkl.Model.Vertices, v)
	}
}

func (p *JklLineParser) parseGeoResourceWorldTextureVertices(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			continue
		}
		var id int32
		v := mgl32.Vec2{}
		n, err := fmt.Sscanf(line, "%d: %f %f", &id, &v[0], &v[1])
		if err != nil {
			panic(err)
		}
		if n != 3 {
			panic("Unable to get texture vertex information")
		}

		p.jkl.Model.TextureVertices = append(p.jkl.Model.TextureVertices, v)
	}
}

func (p *JklLineParser) parseGeoResourceWorldSurfaces(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			continue
		}

		args := p.getLineArgs(line)

		surface := surface{}

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

	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			continue
		}

		var surfaceID int32
		v := mgl32.Vec3{}
		n, err := fmt.Sscanf(line, "%d: %f %f %f", &surfaceID, &v[0], &v[1], &v[2])
		if err != nil {
			//todo: counts are wrong in jkl file, look for 'end' instead?
			//break
			panic(err)
		}
		if n != 4 {
			panic("Unable to get world surface normal information")
		}

		p.jkl.Model.Surfaces[surfaceID].Normal = v
	}
}

func (p *JklLineParser) parseMaterials() {
	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)

		if strings.HasPrefix(line, "SECTION: ") {
			break
		}

		if strings.HasPrefix(line, "WORLD MATERIALS") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD MATERIALS")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseMaterialsWorldMaterials(count)
		}
	}
}

func (p *JklLineParser) parseMaterialsWorldMaterials(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			break
		}

		var id int32
		var matName string
		var xTile float32
		var yTile float32
		n, err := fmt.Sscanf(line, "%d: %s %f %f", &id, &matName, &xTile, &yTile)
		if err != nil {
			break
			//todo: counts are not correct in some jkl files...
			//panic(err)
		}
		if n != 4 {
			panic("Unable to get world material information")
		}

		material := GetLoader().LoadMAT(matName)

		material.XTile = xTile
		material.YTile = yTile

		p.jkl.Model.Materials = append(p.jkl.Model.Materials, material)
	}
}

func (p *JklLineParser) parseModels() {
	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)

		if strings.HasPrefix(line, "SECTION: ") {
			break
		}

		if strings.HasPrefix(line, "WORLD MODELS") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD MODELS")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseModelsWorldModels(count)
		}
	}
}

func (p *JklLineParser) parseModelsWorldModels(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			break
		}

		var id int32
		var jk3doName string
		n, err := fmt.Sscanf(line, "%d: %s", &id, &jk3doName)
		if err != nil {
			//todo: counts are wrong in jkl file, look for 'end' instead?
			break
			panic(err)
		}
		if n != 2 {
			panic("Unable to get world model information")
		}

		jk3do := GetLoader().Load3DO(jk3doName)
		if len(p.jkl.Model.ColorMaps) > 0 {
			jk3do.ColorMap = p.jkl.Model.ColorMaps[0]
		}
		p.jkl.Jk3dos[jk3doName] = jk3do
	}
}

func (p *JklLineParser) parseTemplates() {
	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)

		if strings.HasPrefix(line, "SECTION: ") {
			break
		}

		if strings.HasPrefix(line, "WORLD TEMPLATES") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD TEMPLATES")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseTemplatesWorldTemplates(count)
		}
	}
}

func (p *JklLineParser) parseTemplatesWorldTemplates(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			break
		}

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
			tmp := Template{}
			tmp.Name = name
			tmp.Jk3doName = modelName
			tmp.Size = size

			p.jkl.Jk3doTemplates[tmp.Name] = tmp
		}
	}
}

func (p *JklLineParser) parseThings() {
	for {
		line, ok := p.getNextLine()
		if !ok {
			break
		}
		line = strings.ToUpper(line)

		if strings.HasPrefix(line, "SECTION: ") {
			break
		}

		if strings.HasPrefix(line, "WORLD THINGS") {
			args := p.getLineArgsWithoutPrefix(line, "WORLD THINGS")
			count, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			p.parseThingsWorldThings(count)
		}
	}
}

func (p *JklLineParser) parseThingsWorldThings(count int) {
	for i := 0; i < count; i++ {
		line, ok := p.getNextLine()
		if !ok {
			break
		}

		if line == "end" {
			break
		}

		args := p.getLineArgs(line)

		templateName := args[1]

		x, _ := strconv.ParseFloat(args[3], 64)
		y, _ := strconv.ParseFloat(args[4], 64)
		z, _ := strconv.ParseFloat(args[5], 64)

		pitch, _ := strconv.ParseFloat(args[6], 64)
		yaw, _ := strconv.ParseFloat(args[7], 64)
		Roll, _ := strconv.ParseFloat(args[8], 64)

		t := Thing{}
		t.TemplateName = templateName
		t.Position = mgl32.Vec3{float32(x), float32(y), float32(z)}
		t.Pitch = pitch
		t.Yaw = yaw
		t.Roll = Roll

		p.jkl.Things = append(p.jkl.Things, t)
	}
}
