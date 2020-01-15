package jk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type KeyLineParser struct {
	key     Key
	scanner *bufio.Scanner
	line    string
	done    bool
}

func NewKeyLineParser() *KeyLineParser {
	return &KeyLineParser{}
}

func (p *KeyLineParser) ParseFromFile(filePath string) Key {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseFromString(data)
}

func (p *KeyLineParser) ParseFromString(objString string) Key {
	p.key = Key{}
	p.scanner = bufio.NewScanner(strings.NewReader(objString))
	p.line = ""
	p.done = false

	var err error

	p.getNextLine() // SECTION: HEADER

	p.getNextLine() // FLAGS  %#x
	_, err = fmt.Sscanf(p.line, "flags %v", &p.key.Header.Flags)
	p.checkError(err)

	p.getNextLine() // TYPE   %#x
	_, err = fmt.Sscanf(p.line, "type %v", &p.key.Header.Type)
	p.checkError(err)

	p.getNextLine() // FRAMES %d
	_, err = fmt.Sscanf(p.line, "frames %v", &p.key.Header.Frames)
	p.checkError(err)

	p.getNextLine() // FPS %f
	_, err = fmt.Sscanf(p.line, "fps %v", &p.key.Header.FPS)
	p.checkError(err)

	p.getNextLine() // JOINTS %d
	_, err = fmt.Sscanf(p.line, "joints %v", &p.key.Header.Joints)
	p.checkError(err)

	p.getNextLine() // SECTION: MARKERS (optional) or SECTION: KEYFRAME NODES
	if p.line == "section: markers" {
		p.getNextLine() // MARKERS %d
		p.parseMarkers()
		p.getNextLine() // SECTION: KEYFRAME NODES
	}

	// SECTION: KEYFRAME NODES
	p.getNextLine() // NODES %d
	p.parseNodes()

	return p.key
}

func (p *KeyLineParser) checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (p *KeyLineParser) getNextLine() bool {
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

func (p *KeyLineParser) parseMarkers() {
	var count int
	_, err := fmt.Sscanf(p.line, "markers %d", &count)
	p.checkError(err)

	for i := 0; i < count; i++ {
		// todo: what are markers used for? skip for now
		p.getNextLine() // %f %d
	}
}

func (p *KeyLineParser) parseNodes() {
	var count int
	_, err := fmt.Sscanf(p.line, "nodes %d", &count)
	p.checkError(err)

	for i := 0; i < count; i++ {
		node := KeyframeNode{}

		p.getNextLine() // NODE %d
		p.getNextLine() // MESH NAME %s
		_, err = fmt.Sscanf(p.line, "mesh name %s", &node.MeshName)
		p.checkError(err)

		p.getNextLine() // ENTRIES %d
		p.parseNodeEntries(&node)

		p.key.KeyframeNodes = append(p.key.KeyframeNodes, node)
	}
}

func (p *KeyLineParser) parseNodeEntries(node *KeyframeNode) {
	var count int
	_, err := fmt.Sscanf(p.line, "entries %d", &count)
	p.checkError(err)

	for i := 0; i < count; i++ {
		entry := KeyframeNodeEntry{}

		var id int32
		p.getNextLine()
		_, err = fmt.Sscanf(p.line, "%d: %d %v %f %f %f %f %f %f", &id, &entry.Frame, &entry.Flags,
			&entry.Offset[0], &entry.Offset[1], &entry.Offset[2],
			&entry.Orientation[0], &entry.Orientation[1], &entry.Orientation[2])
		p.checkError(err)

		p.getNextLine()
		_, err = fmt.Sscanf(p.line, "%f %f %f %f %f %f",
			&entry.DeltaOffset[0], &entry.DeltaOffset[1], &entry.DeltaOffset[2],
			&entry.DeltaOrientation[0], &entry.DeltaOrientation[1], &entry.DeltaOrientation[2])
		p.checkError(err)

		node.Entries = append(node.Entries, entry)
	}
}
