package scene

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"
	escene "sunny_land/src/engine/scene"
	"sunny_land/src/engine/ui"
	emath "sunny_land/src/engine/utils/math"
	"sunny_land/src/game/data"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 游戏暂停时显示的菜单场景，提供继续、保存、返回、退出等选项。
 * 该场景通常被推送到 GameScene 之上。
 */
type MenuScene struct {
	// 继承基础场景
	escene.Scene
	// 场景共享数据
	sessionData *data.SessionData
}

// 确保MenuScene实现IScene接口
var _ escene.IScene = (*MenuScene)(nil)

/**
 * @brief MenuScene 的构造函数
 * @param context 引擎上下文的引用
 * @param scene_manager 场景管理器的引用
 * @param session_data_ 场景间传递的游戏数据
 */
func NewMenuScene(ctx *econtext.Context, sceneManager *escene.SceneManager, sd *data.SessionData) *MenuScene {
	ms := &MenuScene{}
	escene.BuildScene(&ms.Scene, "MenuScene", ctx, sceneManager)
	ms.sessionData = sd
	if sd == nil {
		slog.Error("MenuScene: sessionData is nil")
	}
	slog.Debug("MenuScene created")
	return ms
}

// 初始化游戏场景
func (ms *MenuScene) Init() {
	ms.Scene.Init()

	ms.GetContext().GetGameState().SetState(econtext.GameStatePaused)
	ms.CreateUI()
}

// 处理输入事件
func (ms *MenuScene) HandleInput() {
	// 先让 UIManager 处理交互
	ms.Scene.HandleInput()

	// 检查暂停键，允许按暂停键恢复游戏
	if ms.GetContext().GetInputManager().IsActionPressed("pause") {
		// 弹出自身以恢复底层的GameScene
		ms.SceneManager.RequestPopScene()
		ms.GetContext().GetGameState().SetState(econtext.GameStatePlaying)
	}
}

// 创建UI元素
func (ms *MenuScene) CreateUI() {
	screenSize := ms.GetContext().GetGameState().GetLogicalSize()
	if !ms.UIManager.Init(screenSize) {
		slog.Error("UIManager init failed")
		return
	}

	// "PAUSE"标签
	pauseLabel := ui.NewUILabel(ms.GetContext().GetTextRenderer(),
		"PAUSE",
		"assets/fonts/VonwaonBitmap-16px.ttf",
		32,
		emath.FColor{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		mgl32.Vec2{0.0, 0.0},
	)

	// 放在中间靠上的位置
	size := pauseLabel.GetSize()
	labelY := screenSize.Y() * 0.2
	pauseLabel.SetPosition(mgl32.Vec2{(screenSize.X() - size.X()) / 2.0, labelY})
	ms.UIManager.AddElement(pauseLabel)

	// 创建按钮(4个按钮，设定好大小、间距)
	// 按钮稍微小一点
	buttonWidth := float32(96.0)
	buttonHeight := float32(32.0)
	buttonSpacing := float32(10.0)
	// 从标签下方开始，增加间距
	startY := labelY + 80.0
	// 水平居中
	buttonX := (screenSize.X() - buttonWidth) / 2.0

	// Resume Button
	resumeButton := ui.NewUIButton(ms.GetContext(),
		"assets/textures/UI/buttons/Resume1.png",
		"assets/textures/UI/buttons/Resume2.png",
		"assets/textures/UI/buttons/Resume3.png",
		mgl32.Vec2{buttonX, startY},
		mgl32.Vec2{buttonWidth, buttonHeight},
		func() { ms.onResumeClicked() },
	)
	ms.UIManager.AddElement(resumeButton)

	// Save Button
	startY += buttonHeight + buttonSpacing
	saveButton := ui.NewUIButton(ms.GetContext(),
		"assets/textures/UI/buttons/Save1.png",
		"assets/textures/UI/buttons/Save2.png",
		"assets/textures/UI/buttons/Save3.png",
		mgl32.Vec2{buttonX, startY},
		mgl32.Vec2{buttonWidth, buttonHeight},
		func() { ms.onSaveClicked() },
	)
	ms.UIManager.AddElement(saveButton)

	// Back Button
	startY += buttonHeight + buttonSpacing
	backButton := ui.NewUIButton(ms.GetContext(),
		"assets/textures/UI/buttons/Back1.png",
		"assets/textures/UI/buttons/Back2.png",
		"assets/textures/UI/buttons/Back3.png",
		mgl32.Vec2{buttonX, startY},
		mgl32.Vec2{buttonWidth, buttonHeight},
		func() { ms.onBackClicked() },
	)
	ms.UIManager.AddElement(backButton)

	// Quit Button
	startY += buttonHeight + buttonSpacing
	quitButton := ui.NewUIButton(ms.GetContext(),
		"assets/textures/UI/buttons/Quit1.png",
		"assets/textures/UI/buttons/Quit2.png",
		"assets/textures/UI/buttons/Quit3.png",
		mgl32.Vec2{buttonX, startY},
		mgl32.Vec2{buttonWidth, buttonHeight},
		func() { ms.onQuitClicked() },
	)
	ms.UIManager.AddElement(quitButton)
}

// 按钮回调函数实现
// 继续游戏按钮回调
func (ms *MenuScene) onResumeClicked() {
	// 弹出当前场景，恢复游戏
	ms.SceneManager.RequestPopScene()
	ms.GetContext().GetGameState().SetState(econtext.GameStatePlaying)
}

// 保存游戏按钮回调
func (ms *MenuScene) onSaveClicked() {
	// 保存游戏数据
	if ms.sessionData != nil && ms.sessionData.SaveToFile("assets/save.json") {
		slog.Debug("save game success")
	} else {
		slog.Error("save game failed")
	}
}

// 返回按钮回调
func (ms *MenuScene) onBackClicked() {
	// 直接替换为TitleScene
	ms.SceneManager.RequestReplaceScene(
		NewTitleScene(ms.GetContext(), ms.SceneManager, ms.sessionData))
}

// 退出按钮回调
func (ms *MenuScene) onQuitClicked() {
	// 输入管理器设置退出标志
	ms.GetContext().GetInputManager().SetShouldQuit(true)
}
