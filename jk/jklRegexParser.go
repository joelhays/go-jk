package jk

import (
	"bufio"
	"github.com/joelhays/go-jk/jk/jktypes"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type JklRegexParser struct {
	jkl     jktypes.Jkl
	scanner *bufio.Scanner
	line    string
	done    bool
}

func NewJklRegexParser() *JklRegexParser {
	return &JklRegexParser{
		jkl: jktypes.Jkl{
			Model:          &jktypes.JkMesh{},
			Jk3dos:         make(map[string]jktypes.Jk3doFile),
			Jk3doTemplates: make(map[string]jktypes.Template),
			Things:         nil,
		},
	}
}

// ReadJKLFromFile will read a .jkl file and return a struct containing all necessary information
func (p *JklRegexParser) ParseJKLFromFile(filePath string) jktypes.Jkl {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseJKLFromString(data)
}

// ReadJKLFromString will parse a string as a .jkl file
func (p *JklRegexParser) ParseJKLFromString(jklString string) jktypes.Jkl {
	data := jklString

	jklResult := jktypes.Jkl{}
	jklResult.Model = &jktypes.JkMesh{}

	jklResult.Jk3dos = make(map[string]jktypes.Jk3doFile)
	jklResult.Jk3doTemplates = make(map[string]jktypes.Template)

	p.parseVertices(data, &jklResult)
	p.parseTextureVertices(data, &jklResult)
	p.parseMaterials(data, &jklResult)
	p.parseColormaps(data, &jklResult)
	p.parseSurfaces(data, &jklResult)

	p.parse3dos(data, &jklResult)
	p.parse3doTemplates(data, &jklResult)
	p.parseThings(data, &jklResult)

	return jklResult
}

func (p *JklRegexParser) parseSection(data string, regex string, componentRegex string, callback func(components []string)) {
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

func (p *JklRegexParser) parseVertices(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World vertices.*World texture vertices`, "\\d+:.*",
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

			jklResult.Model.Vertices = append(jklResult.Model.Vertices, mgl32.Vec3{float32(x), float32(y), float32(z)})
		})
}

func (p *JklRegexParser) parseTextureVertices(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World texture vertices.*World adjoins`, "\\d+:.*",
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

			jklResult.Model.TextureVertices = append(jklResult.Model.TextureVertices, mgl32.Vec2{float32(u), float32(v)})
		})
}

func (p *JklRegexParser) parseMaterials(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World materials.*SECTION: GEORESOURCE`, "\\d+:.*",
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

			material := GetLoader().LoadMAT(matName)

			material.XTile = float32(xTile)
			material.YTile = float32(yTile)

			jklResult.Model.Materials = append(jklResult.Model.Materials, material)
		})
}

func (p *JklRegexParser) parseColormaps(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World Colormaps.*World vertices`, "\\d+:.*",
		func(components []string) {
			var cmpName string
			if len(components) == 0 {
				cmpName = "dflt.cmp"
			} else {
				cmpName = components[1]
			}

			colorMap := GetLoader().LoadCMP(cmpName)

			jklResult.Model.ColorMaps = append(jklResult.Model.ColorMaps, colorMap)
		})
}

func (p *JklRegexParser) parseSurfaces(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World surfaces.*\#--- Surface normals ---`, "\\d+:.*",
		func(components []string) {
			surface := jktypes.Surface{}

			materialID, _ := strconv.ParseInt(components[1], 10, 32)
			surface.MaterialID = materialID

			geoFlag, _ := strconv.ParseInt(components[4], 10, 32)
			surface.Geo = geoFlag

			// TODO: WHAT DOES THIS VALUE MEAN?
			//if components[5] != "3" {
			//	fmt.Println("light != 3", components[5])
			//}

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

			jklResult.Model.Surfaces = append(jklResult.Model.Surfaces, surface)
		})

	p.parseSection(data, `(?s)\#--- Surface normals ---.*Section: SECTORS`, "\\d+:.*",
		func(components []string) {
			surfaceID, _ := strconv.ParseInt(strings.TrimRight(components[0], ":"), 10, 32)

			x, _ := strconv.ParseFloat(components[1], 64)
			y, _ := strconv.ParseFloat(components[2], 64)
			z, _ := strconv.ParseFloat(components[3], 64)

			jklResult.Model.Surfaces[surfaceID].Normal = mgl32.Vec3{float32(x), float32(y), float32(z)}
		})
}

func (p *JklRegexParser) parse3dos(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World models.*Section: SPRITES`, "\\d+:.*",
		func(components []string) {
			jk3doName := components[1]

			jk3do := GetLoader().Load3DO(jk3doName)
			if len(jklResult.Model.ColorMaps) > 0 {
				jk3do.ColorMap = jklResult.Model.ColorMaps[0]
			}
			jklResult.Jk3dos[jk3doName] = jk3do
		})
}

func (p *JklRegexParser) parse3doTemplates(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World templates.*Section: Things`, ".*",
		func(components []string) {
			if len(components) < 3 {
				return
			}

			name := components[0]
			var modelName string
			size := 1.0
			for i := 0; i < len(components); i++ {
				if strings.HasPrefix(components[i], "size=") {
					size, _ = strconv.ParseFloat(strings.TrimPrefix(components[i], "size="), 32)
				}
				if strings.HasPrefix(components[i], "model3d=") {
					modelName = strings.TrimPrefix(components[i], "model3d=")
				}
			}

			if modelName != "" {
				tmp := jktypes.Template{}
				tmp.Name = name
				tmp.Jk3doName = modelName
				tmp.Size = size

				jklResult.Jk3doTemplates[tmp.Name] = tmp
			}
		})
}

func (p *JklRegexParser) parseThings(data string, jklResult *jktypes.Jkl) {
	p.parseSection(data, `(?s)World things.*end`, "\\d+:.*",
		func(components []string) {
			templateName := components[1]

			x, _ := strconv.ParseFloat(components[3], 64)
			y, _ := strconv.ParseFloat(components[4], 64)
			z, _ := strconv.ParseFloat(components[5], 64)

			pitch, _ := strconv.ParseFloat(components[6], 64)
			yaw, _ := strconv.ParseFloat(components[7], 64)
			Roll, _ := strconv.ParseFloat(components[8], 64)

			t := jktypes.Thing{}
			t.TemplateName = templateName
			t.Position = mgl32.Vec3{float32(x), float32(y), float32(z)}
			t.Pitch = pitch
			t.Yaw = yaw
			t.Roll = Roll

			jklResult.Things = append(jklResult.Things, t)
		})
}
