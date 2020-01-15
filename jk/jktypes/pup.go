package jktypes

type Pup struct {
	Modes  []PupMode
	Joints []PupJoint
}

type PupMode struct {
	SubModes    []PupSubMode
	BasedOn     int32
	IsInherited bool
}

type PupSubMode struct {
	Name     string
	Keyframe string
	Flags    byte
	LoPri    int32
	HiPri    int32
}

type PupJoint struct {
	Joint int32
	Node  int32
}
