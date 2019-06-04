package scene

type SceneManager struct {
	scenes      map[string]Scene
	activeScene string
}

func NewSceneManager() *SceneManager {
	return &SceneManager{scenes: make(map[string]Scene)}
}

func (m *SceneManager) Add(key string, scene Scene) {
	m.scenes[key] = scene
}

func (m *SceneManager) LoadScene(key string) {
	if scene, ok := m.scenes[m.activeScene]; ok {
		scene.Unload()
	}

	m.activeScene = key
	if scene, ok := m.scenes[m.activeScene]; ok {
		scene.Load()
	}
}

func (m *SceneManager) Update() {
	if scene, ok := m.scenes[m.activeScene]; ok {
		scene.Update()
	}
}

func (m *SceneManager) Unload() {
	if scene, ok := m.scenes[m.activeScene]; ok {
		scene.Unload()
	}
}
