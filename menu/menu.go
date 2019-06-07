package menu

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
	"github.com/joelhays/go-jk/scene"
	"log"
)

var (
	levels = []string{"01narshadda.jkl", "02narshadda.jkl", "03katarn.jkl", "04escapehouse.jkl", "06abarons.jkl",
		"06bbarons.jkl", "07yun.jkl", "08escape88.jkl", "09fuelstation.jkl", "10cargo.jkl", "11gorc.jkl", "12escape.jkl",
		"14tower.jkl", "15maw.jkl", "16aescapeship.jkl", "16bescapeship.jkl", "17asarris.jkl", "17bsarris.jkl",
		"18ascend.jkl", "19a.jkl", "19b.jkl", "20aboc.jkl", "20bboc.jkl", "21ajarec.jkl", "21bjarec.jkl"}
)

type MainMenu struct {
	window       *glfw.Window
	context      *nk.Context
	textureId    uint32
	sceneManager *scene.SceneManager
	fontAtlas    *nk.FontAtlas
	font         *nk.Font
	fontHandle   *nk.UserFont
}

func NewMainMenu(window *glfw.Window, sceneManager *scene.SceneManager) *MainMenu {
	return &MainMenu{window: window, sceneManager: sceneManager}
}

func (m *MainMenu) Init() {
	m.context = nk.NkPlatformInit(m.window, nk.PlatformInstallCallbacks)
	m.fontAtlas = nk.NewFontAtlas()
	nk.NkFontStashBegin(&m.fontAtlas)
	m.font = nk.NkFontAtlasAddFromBytes(m.fontAtlas, MustAsset("assets/FreeSans.ttf"), 24, nil)
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
}

func (m *MainMenu) Update() {
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
			nk.NkListViewBegin(m.context, &list, "level", nk.WindowBackground, 30, int32(len(levels)-1))
			{
				for l := list.Begin(); l < list.End(); l++ {
					level := levels[l]
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

func (m *MainMenu) Unload() {
	//nk.NkPlatformShutdown()
}
