package scene

type Menu interface {
	Init()
	Update()
	Unload()
}

type MenuScene struct {
	menu Menu
}

func NewMenuScene(menu Menu) *MenuScene {
	return &MenuScene{menu: menu}
}

func (s *MenuScene) Load() {
	s.menu.Init()
}

func (s *MenuScene) Unload() {
	s.menu.Unload()
}

func (s *MenuScene) Update() {
	s.menu.Update()
}
