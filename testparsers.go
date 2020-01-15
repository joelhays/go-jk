package main

import (
	"fmt"
	"github.com/joelhays/go-jk/jk"
)

func testPupParser() {
	p := jk.NewPupLineParser()
	pup := p.ParseFromFile("./_testfiles/rystr.pup")
	fmt.Println(fmt.Sprintf("%+v", pup))

	manifest := jk.GetLoader().LoadResourceManifest("misc\\pup\\", "pup")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}

func testKeyParser() {
	p := jk.NewKeyLineParser()
	r := p.ParseFromFile("./_testfiles/8twalk.key")
	fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadResourceManifest("3do\\key\\", "key")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}

func testJklParser() {
	p := jk.NewJklLineParser()
	r := p.ParseFromFile("./_testfiles/jkl/01narshadda.jkl")
	_ = r
	//fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadEpisodeManifest("jkl\\", "jkl")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r = p.ParseFromString(string(fileBytes))
	}
}

func test3doParser() {
	p := jk.NewJk3doLineParser()
	r := p.ParseFromFile("./_testfiles/3do/rystr.3do")
	_ = r
	//fmt.Println(fmt.Sprintf("%+v", r))

	manifest := jk.GetLoader().LoadResourceManifest("3do\\", "3do")
	for _, file := range manifest {
		fmt.Println(file)
		fileBytes := jk.GetLoader().LoadResource(file)
		r := p.ParseFromString(string(fileBytes))
		_ = r
	}
}
