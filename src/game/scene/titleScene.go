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
 * @brief 标题场景类，提供4个按钮：开始游戏、加载游戏、帮助、退出
 */
type TitleScene struct {
	// 继承基础场景
	escene.Scene
	// 场景共享数据
	sessionData *data.SessionData
}

// 确保TitleScene实现IScene接口
var _ escene.IScene = (*TitleScene)(nil)

/**
 * @brief 构造函数
 * @param context 引擎上下文
 * @param scene_manager 场景管理器
 * @param game_play_state 指向游戏玩法状态的共享指针
 */
func NewTitleScene(ctx *econtext.Context, sceneManager *escene.SceneManager, sd *data.SessionData) *TitleScene {
	ts := &TitleScene{}
	escene.BuildScene(&ts.Scene, "TitleScene", ctx, sceneManager)
	ts.sessionData = sd
	if ts.sessionData == nil {
		ts.sessionData = data.NewSessionData()
	}
	slog.Debug("TitleScene created")
	return ts
}

// 初始化游戏场景
func (ts *TitleScene) Init() {
	ts.Scene.Init()

	// 加载背景地图
	if !escene.NewLevelLoader().LoadLevel("assets/maps/level0.tmj", ts) {
		slog.Error("level0.tmj init failed")
		return
	}

	// 创建UI元素
	ts.createUI()

	// 设置音量
	// 设置背景音乐音量为20%
	ts.GetContext().GetAudioPlayer().SetMusicVolume(0.2)
	// 设置音效音量为50%
	ts.GetContext().GetAudioPlayer().SetSoundVolume(0.5)

	// 设置背景音乐
	ts.GetContext().GetAudioPlayer().PlayMusic("assets/audio/platformer_level03_loop.ogg", true)

	slog.Debug("TitleScene init done")
}

// 创建UI元素
func (ts *TitleScene) createUI() {
	screenSize := mgl32.Vec2{640.0, 360.0}

	if !ts.UIManager.Init(screenSize) {
		slog.Error("ui manager init failed")
		return
	}

	// 创建标题图片 (假设不知道大小)
	titleImage := ui.NewUIImage("assets/textures/UI/title-screen.png", mgl32.Vec2{0.0, 0.0}, mgl32.Vec2{0.0, 0.0}, nil, false)
	size := ts.GetContext().GetResourceManager().GetTextureSize(titleImage.GetTextureId())
	// 放大为2倍
	titleImage.SetSize(size.Mul(2.0))

	// 水平居中
	titleImage.SetPosition(screenSize.Sub(titleImage.GetSize()).Mul(0.5).Sub(mgl32.Vec2{0.0, 50.0}))
	ts.UIManager.AddElement(titleImage)

	// 创建按钮面板并居中，4个按钮，设定好大小、间距
	buttonWidth := float32(96.0)
	buttonHeight := float32(32.0)
	buttonSpacing := float32(20.0)
	numButtons := 4

	// 计算面板总宽高
	panelWidth := float32(numButtons)*buttonWidth + (float32(numButtons-1))*buttonSpacing
	panelHeight := buttonHeight

	// 计算面板位置使其居中
	panelX := (screenSize.X() - panelWidth) / 2.0
	// 垂直位置中间靠下
	panelY := screenSize.Y() * 0.65

	// 创建按钮面板
	buttonPanel := ui.NewUIPanel(mgl32.Vec2{panelX, panelY}, mgl32.Vec2{panelWidth, panelHeight}, nil)

	// 创建按钮并添加到 UIPanel (位置是相对于 UIPanel 的 0,0)
	curButtonPos := mgl32.Vec2{0.0, 0.0}
	buttonSize := mgl32.Vec2{buttonWidth, buttonHeight}
	// Start Button
	startButton := ui.NewUIButton(ts.GetContext(),
		"assets/textures/UI/buttons/Start1.png",
		"assets/textures/UI/buttons/Start2.png",
		"assets/textures/UI/buttons/Start3.png",
		curButtonPos,
		buttonSize,
		func() { ts.onStartGameClick() },
	)
	buttonPanel.AddChild(startButton)

	// Load Button
	curButtonPos[0] += buttonWidth + buttonSpacing
	loadButton := ui.NewUIButton(ts.GetContext(),
		"assets/textures/UI/buttons/Load1.png",
		"assets/textures/UI/buttons/Load2.png",
		"assets/textures/UI/buttons/Load3.png",
		curButtonPos,
		buttonSize,
		func() { ts.onLoadGameClick() },
	)
	buttonPanel.AddChild(loadButton)

	// Helps Button
	curButtonPos[0] += buttonWidth + buttonSpacing
	helpsButton := ui.NewUIButton(ts.GetContext(),
		"assets/textures/UI/buttons/Helps1.png",
		"assets/textures/UI/buttons/Helps2.png",
		"assets/textures/UI/buttons/Helps3.png",
		curButtonPos,
		buttonSize,
		func() { ts.onHelpsClick() },
	)
	buttonPanel.AddChild(helpsButton)

	// Quit Button
	curButtonPos[0] += buttonWidth + buttonSpacing
	quitButton := ui.NewUIButton(ts.GetContext(),
		"assets/textures/UI/buttons/Quit1.png",
		"assets/textures/UI/buttons/Quit2.png",
		"assets/textures/UI/buttons/Quit3.png",
		curButtonPos,
		buttonSize,
		func() { ts.onQuitClick() },
	)
	buttonPanel.AddChild(quitButton)

	// 将 UIPanel 添加到UI管理器
	ts.UIManager.AddElement(buttonPanel)

	// 创建 Credits 标签
	creditsLabel := ui.NewUILabel(ts.GetContext().GetTextRenderer(),
		"SunnyLand Credits: XXX - 2026",
		"assets/fonts/VonwaonBitmap-16px.ttf",
		16,
		emath.FColor{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
		mgl32.Vec2{0.0, 0.0},
	)
	creditsLabel.SetPosition(mgl32.Vec2{(screenSize.X() - creditsLabel.GetSize().X()) / 2.0,
		screenSize.Y() - creditsLabel.GetSize().Y() - 10.0})
	ts.UIManager.AddElement(creditsLabel)
}

// 更新
func (ts *TitleScene) Update(dt float64) {
	ts.Scene.Update(dt)
	// 相机自动向右移动
	ts.GetContext().GetCamera().Move(mgl32.Vec2{float32(dt) * 100.0, 0.0})
}

// 按钮回调实现
// 开始游戏按钮点击回调
func (ts *TitleScene) onStartGameClick() {
	if ts.sessionData != nil {
		ts.sessionData.Reset()
	}
	ts.SceneManager.RequestReplaceScene(NewGameScene(ts.GetContext(), ts.SceneManager, ts.sessionData))
}

// 加载游戏按钮点击回调
func (ts *TitleScene) onLoadGameClick() {
	if ts.sessionData == nil {
		slog.Error("sessionData is nil")
		return
	}

	if ts.sessionData.LoadFromFile("assets/save.json") {
		slog.Debug("save file load success, start game...")
		ts.SceneManager.RequestReplaceScene(NewGameScene(ts.GetContext(), ts.SceneManager, ts.sessionData))
	} else {
		slog.Warn("load save file failed")
	}
}

// 帮助按钮点击回调
func (ts *TitleScene) onHelpsClick() {
	slog.Debug("helps button click")
	ts.SceneManager.RequestPushScene(NewHelpsScene(ts.GetContext(), ts.SceneManager))
}

// 退出按钮点击回调
func (ts *TitleScene) onQuitClick() {
	slog.Debug("quit button click")
	ts.GetContext().GetInputManager().SetShouldQuit(true)
}
