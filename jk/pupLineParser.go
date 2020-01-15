package jk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type PupLineParser struct {
	pup     Pup
	scanner *bufio.Scanner
	line    string
	done    bool
}

func NewPupLineParser() *PupLineParser {
	return &PupLineParser{
		pup: Pup{},
	}
}

func (p *PupLineParser) ParseFromFile(filePath string) Pup {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)

	return p.ParseFromString(data)
}

func (p *PupLineParser) ParseFromString(objString string) Pup {
	p.pup = Pup{}
	p.scanner = bufio.NewScanner(strings.NewReader(objString))
	p.done = false

	var mode *PupMode
	for {
		p.getNextLine()
		if p.done {
			break
		}

		if strings.HasPrefix(p.line, "mode=") {
			var args int
			var modeNum int32
			var basedon int32
			args, _ = fmt.Sscanf(p.line, "mode=%d, basedon=%d colormaps %d", &modeNum, &basedon)
			p.pup.Modes = append(p.pup.Modes, PupMode{
				SubModes:    make([]PupSubMode, 0),
				BasedOn:     basedon,
				isInherited: args == 2,
			})
			mode = &p.pup.Modes[len(p.pup.Modes)-1]
			continue
		} else if p.line == "joints" {
			for {
				p.getNextLine()
				if p.line == "end" {
					break
				}
				joint := PupJoint{}
				_, _ = fmt.Sscanf(p.line, "%d=%d", &joint.Joint, &joint.Node)
				p.pup.Joints = append(p.pup.Joints, joint)
			}
		} else {
			if mode == nil {
				panic("Processing submode without a mode")
			}
			subMode := PupSubMode{}
			_, _ = fmt.Sscanf(p.line, "%s %s %v %d %d", &subMode.Name, &subMode.Keyframe, &subMode.Flags,
				&subMode.LoPri, &subMode.HiPri)
			mode.SubModes = append(mode.SubModes, subMode)
		}
	}

	return p.pup
}

func (p *PupLineParser) getNextLine() bool {
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
