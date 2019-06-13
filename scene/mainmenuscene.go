package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/menu"
	"github.com/joelhays/go-jk/opengl"
	"log"
	"strings"
)

type MainMenuScene struct {
	window       *glfw.Window
	context      *nk.Context
	textureId    uint32
	sceneManager *SceneManager
	fontAtlas    *nk.FontAtlas
	font         *nk.Font
	fontHandle   *nk.UserFont
	levels       []string
}

func NewMainMenuScene(window *glfw.Window, sceneManager *SceneManager) *MainMenuScene {
	return &MainMenuScene{window: window, sceneManager: sceneManager}
}

func (m *MainMenuScene) Load() {
	m.context = nk.NkPlatformInit(m.window, nk.PlatformInstallCallbacks)
	m.fontAtlas = nk.NewFontAtlas()
	nk.NkFontStashBegin(&m.fontAtlas)
	m.font = nk.NkFontAtlasAddFromBytes(m.fontAtlas, menu.MustAsset("assets/FreeSans.ttf"), 24, nil)
	//m.font = nk.NkFontAtlasAddDefault(atlas, 16, nil)
	nk.NkFontStashEnd()
	if m.font != nil {
		m.fontHandle = m.font.Handle()
		nk.NkStyleSetFont(m.context, m.fontHandle)
	}

	bmFile := jk.GetLoader().LoadBM("bkmain.bm")
	bmRenderer := opengl.NewOpenGlBmRenderer(&bmFile, nil)
	original, ok := bmRenderer.(*opengl.OpenGlBmRenderer)
	if ok {
		//*m.context.GetStyle().GetWindow().GetFixedBackground() = nk.NkStyleItemImage(nk.NkSubimageId(int32(original.GetTextureID()), 1024, 768, nk.NkRect(0, 0, 1024, 768)))
		m.textureId = original.GetTextureID()
	}

	//bmFile2 := jk.GetLoader().LoadBM("bksingle.bm")
	//bmRenderer2 := opengl.NewOpenGlBmRenderer(&bmFile2, nil)
	//original2, ok2 := bmRenderer2.(*opengl.OpenGlBmRenderer)
	//if ok2 {
	//	*m.context.GetStyle().GetButton().GetNormal() = nk.NkStyleItemImage(nk.NkSubimageId(int32(original2.GetTextureID()), 24, 24, nk.NkRect(0, 0, 24, 24)))
	//}

	m.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)

	if len(m.levels) == 0 {
		for _, gobFileName := range jk.GetLoader().LoadJKLManifest() {
			if strings.HasPrefix(gobFileName, "jkl\\") && strings.HasSuffix(gobFileName, "jkl") {
				m.levels = append(m.levels, gobFileName)
			}
		}
	}
}

func (m *MainMenuScene) Unload() {
	//nk.NkPlatformShutdown()
}

func (m *MainMenuScene) Update() {
	nk.NkPlatformNewFrame()

	// Layout
	*m.context.GetStyle().GetWindow().GetFixedBackground() = nk.NkStyleItemImage(nk.NkSubimageId(int32(m.textureId), 1024, 768, nk.NkRect(0, 0, 1024, 768)))

	bounds := nk.NkRect(0, 0, 1024, 768)
	update := nk.NkBegin(m.context, "Demo", bounds, nk.WindowBackground)

	*m.context.GetStyle().GetWindow().GetFixedBackground() = nk.NkStyleItemHide()

	if update > 0 {

		nk.NkLayoutRowDynamic(m.context, 250, 1)
		{
		}

		nk.NkLayoutRowStatic(m.context, 30, 1024/3, 3)
		{
			nk.NkSpacing(m.context, 1)
			nk.NkLabel(m.context, "Select an level to load:", nk.TextCentered)
		}

		nk.NkLayoutRowStatic(m.context, 30*10, 1024/3, 3)
		{
			nk.NkSpacing(m.context, 1)
			var list nk.ListView
			nk.NkListViewBegin(m.context, &list, "level", nk.WindowBackground, 35, int32(len(m.levels)-1))
			{
				for l := list.Begin(); l < list.End(); l++ {
					level := m.levels[l]
					nk.NkLayoutRowDynamic(m.context, 30, 1)
					{
						if nk.NkButtonLabel(m.context, level) > 0 {
							log.Println("[INFO] button pressed! " + level)
							go m.sceneManager.LoadScene(level)
						}
						//nk.NkLabel(m.context, "Item "+strconv.Itoa(l), nk.TextLeft)
					}
				}
			}
			nk.NkListViewEnd(&list)
		}

		nk.NkLayoutRowDynamic(m.context, 30, 1)
		{
		}
		nk.NkLayoutRowStatic(m.context, 30, 1024/3, 3)
		{
			nk.NkSpacing(m.context, 1)
			if nk.NkButtonLabel(m.context, "Quit") > 0 {
				m.window.SetShouldClose(true)
			}
		}

	}
	nk.NkEnd(m.context)

	maxVertexBuffer := 512 * 1024
	maxElementBuffer := 128 * 1024
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
}
