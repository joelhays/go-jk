package menu

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang-ui/nuklear/nk"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
	"log"
)

type Menu struct {
	window    *glfw.Window
	context   *nk.Context
	textureId uint32
	list      *nk.ListView
}

func NewMenu(window *glfw.Window) *Menu {
	return &Menu{window: window}
}

func (m *Menu) AddButton(text string, position mgl32.Vec2) {

}

func (m *Menu) SetBackground(bmFile string) {

}

func (m *Menu) Init() {
	m.context = nk.NkPlatformInit(m.window, nk.PlatformInstallCallbacks)
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 24, nil)
	//sansFont := nk.NkFontAtlasAddDefault(atlas, 16, nil)
	nk.NkFontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(m.context, sansFont.Handle())
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

	m.list = &nk.ListView{}
}

func (m *Menu) Update() {
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
			nk.NkLabel(m.context, "Select an option:", nk.TextCentered)
		}

		nk.NkLayoutRowStatic(m.context, 30*5, 1024/3, 3)
		{
			nk.NkSpacing(m.context, 1)
			//var list nk.ListView
			nk.NkListViewBegin(m.context, m.list, "level", nk.WindowBackground, 30, 5)
			{

				nk.NkLayoutRowDynamic(m.context, 30, 1)
				{
					nk.NkLabel(m.context, "Item 1", nk.TextLeft)
					nk.NkLabel(m.context, "Item 2", nk.TextLeft)
					nk.NkLabel(m.context, "Item 3", nk.TextLeft)
					nk.NkLabel(m.context, "Item 4", nk.TextLeft)
					nk.NkLabel(m.context, "Item 5", nk.TextLeft)
				}
			}
			nk.NkListViewEnd(m.list)
		}

		nk.NkLayoutRowStatic(m.context, 30, 1024/3, 3)
		{
			nk.NkSpacing(m.context, 1)
			if nk.NkButtonLabel(m.context, "Load Level") > 0 {
				log.Println("[INFO] button pressed!")
			}
		}

	}
	nk.NkEnd(m.context)

	maxVertexBuffer := 512 * 1024
	maxElementBuffer := 128 * 1024
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
}

func (m *Menu) Unload() {
	nk.NkPlatformShutdown()
}
