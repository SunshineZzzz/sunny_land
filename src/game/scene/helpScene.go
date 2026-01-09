package scene

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"
	escene "sunny_land/src/engine/scene"
	"sunny_land/src/engine/ui"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 显示帮助信息的场景。
 *
 * 显示一张帮助图片，按鼠标左键退出。
 */
type HelpsScene struct {
	// 继承基础场景
	escene.Scene
}

// 确保HelpScene实现IScene接口
var _ escene.IScene = (*HelpsScene)(nil)

// 构造函数
func NewHelpsScene(context *econtext.Context, sceneManager *escene.SceneManager) *HelpsScene {
	hs := &HelpsScene{}
	escene.BuildScene(&hs.Scene, "HelpScene", context, sceneManager)
	slog.Debug("HelpScene created")
	return hs
}

// 初始化游戏场景
func (hs *HelpsScene) Init() {
	hs.Scene.Init()

	screenSize := hs.GetContext().GetGameState().GetLogicalSize()
	if !hs.UIManager.Init(screenSize) {
		slog.Error("ui manager init failed")
		return
	}

	// 创建帮助图片 UIImage （让它覆盖整个屏幕）
	helpImage := ui.NewUIImage("assets/textures/UI/instructions.png", mgl32.Vec2{0.0, 0.0}, screenSize, nil, false)
	hs.UIManager.AddElement(helpImage)
}

// 处理输入事件
func (hs *HelpsScene) HandleInput() {
	if !hs.IsInitialized() {
		return
	}

	// 检测是否按下鼠标左键
	if hs.GetContext().GetInputManager().IsActionPressed("MouseLeftClick") {
		// 鼠标左键被按下, 退出
		hs.SceneManager.RequestPopScene()
	}
}
