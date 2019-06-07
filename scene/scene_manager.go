package scene

import (
	"log"
)

type SceneManager struct {
	scenes      map[string]Scene
	activeScene string
	loading     bool
}

func NewSceneManager() *SceneManager {
	return &SceneManager{scenes: make(map[string]Scene)}
}

func (m *SceneManager) Add(key string, scene Scene) {
	m.scenes[key] = scene
}

func (m *SceneManager) LoadScene(key string) {
	if m.loading {
		return
	}

	m.loading = true

	if m.activeScene == key {
		m.loading = false
		return
	}

	log.Println("[INFO] started loading " + key)

	if scene, ok := m.scenes[m.activeScene]; ok {
		scene.Unload()
	}

	if scene, ok := m.scenes[key]; ok {
		scene.Load()
	}
	m.activeScene = key

	log.Println("[INFO] finished loading " + key)

	m.loading = false
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
