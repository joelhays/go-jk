package main

import (
	"fmt"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/jk/jkparsers"
)

func testPupParser() {
	p := jkparsers.NewPupLineParser()
	pup := p.ParseFromFile("./_testfiles/rystr.pup")
	fmt.Println(fmt.Sprintf("%+v", pup))

	manifest := jk.GetLoader().LoadManifest("pup")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}

func testKeyParser() {
	p := jkparsers.NewKeyLineParser()
	r := p.ParseFromFile("./_testfiles/8twalk.key")
	fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadManifest("key")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}

func testJklParser() {
	p := jkparsers.NewJklLineParser()
	r := p.ParseFromFile("./_testfiles/jkl/01narshadda.jkl")
	_ = r
	//fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadManifest("jkl")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadEpisode(file)
		r = p.ParseFromString(string(fileBytes))
	}
}

func test3doParser() {
	p := jkparsers.NewJk3doLineParser()
	r := p.ParseFromFile("./_testfiles/3do/rystr.3do")
	_ = r
	//fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadManifest("3do")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}
