package scene

type Scene interface {
	Load()
	Unload()
	Update()
}
