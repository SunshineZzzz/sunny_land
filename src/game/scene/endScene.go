package scene

import (
	"log/slog"
	"strconv"
	econtext "sunny_land/src/engine/context"
	escene "sunny_land/src/engine/scene"
	"sunny_land/src/engine/ui"
	emath "sunny_land/src/engine/utils/math"
	"sunny_land/src/game/data"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @class EndScene
 * @brief 显示游戏结束（胜利或失败）信息的叠加场景。
 *
 * 提供重新开始或返回主菜单的选项。
 */
type EndScene struct {
	// 继承基础场景
	escene.Scene
	// 会话数据
	sessionData *data.SessionData
}

// 确保EndScene实现IScene接口
var _ escene.IScene = (*EndScene)(nil)

/**
 * @brief 构造函数
 * @param context 引擎上下文
 * @param scene_manager 场景管理器
 * @param session_data 指向游戏数据状态的共享指针
 */
func NewEndScene(context *econtext.Context, sceneManager *escene.SceneManager, sessionData *data.SessionData) *EndScene {
	es := &EndScene{
		sessionData: sessionData,
	}
	escene.BuildScene(&es.Scene, "EndScene", context, sceneManager)
	slog.Debug("EndScene created")
	return es
}

// 初始化场景
func (es *EndScene) Init() {
	es.Scene.Init()

	// 设置游戏状态为 GameOver
	es.GetContext().GetGameState().SetState(econtext.GameStateGameOver)

	es.createUI()
}

// 创建UI元素
func (es *EndScene) createUI() {
	window_size := es.GetContext().GetGameState().GetLogicalSize()
	if !es.UIManager.Init(window_size) {
		slog.Error("EndScene UI manager init failed")
		return
	}

	isWin := es.sessionData.GetIsWin()

	// --- 主文字标签 ---
	main_message := "YOU WIN! CONGRATS!"
	if !isWin {
		main_message = "YOU DIED! TRY AGAIN!"
	}
	// 赢了是绿色，输了是红色
	messageColor := emath.FColor{R: 0.0, G: 1.0, B: 0.0, A: 1.0}
	if !isWin {
		messageColor = emath.FColor{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
	}

	mainLabel := ui.NewUILabel(es.GetContext().GetTextRenderer(),
		main_message,
		"assets/fonts/VonwaonBitmap-16px.ttf",
		48,
		messageColor,
		mgl32.Vec2{0.0, 0.0},
	)

	// 标签居中
	label_size := mainLabel.GetSize()
	main_label_pos := mgl32.Vec2{(window_size.X() - label_size.X()) / 2.0, window_size.Y() * 0.3}
	mainLabel.SetPosition(main_label_pos)
	es.UIManager.AddElement(mainLabel)

	// 得分标签
	currentScore := es.sessionData.GetCurrentScore()
	highScore := es.sessionData.GetHighScore()
	scoreColor := emath.FColor{R: 1.0, G: 1.0, B: 1.0, A: 1.0}
	scoreFontSize := 24

	// 当前得分
	scoreText := "Score: " + strconv.Itoa(currentScore)
	score_label := ui.NewUILabel(es.GetContext().GetTextRenderer(),
		scoreText,
		"assets/fonts/VonwaonBitmap-16px.ttf",
		scoreFontSize,
		scoreColor,
		mgl32.Vec2{0.0, 0.0},
	)
	scoreLabelSize := score_label.GetSize()
	// x方向居中，y方向在主标签下方20像素
	score_label.SetPosition(mgl32.Vec2{(window_size.X() - scoreLabelSize.X()) / 2.0, main_label_pos.Y() + label_size.Y() + 20.0})
	es.UIManager.AddElement(score_label)

	// 最高分
	highScoreText := "High Score: " + strconv.Itoa(highScore)
	high_score_label := ui.NewUILabel(es.GetContext().GetTextRenderer(),
		highScoreText,
		"assets/fonts/VonwaonBitmap-16px.ttf",
		scoreFontSize,
		scoreColor,
		mgl32.Vec2{0.0, 0.0},
	)
	highScoreLabelSize := high_score_label.GetSize()
	// x方向居中，y方向在当前得分下方10像素
	high_score_label.SetPosition(mgl32.Vec2{(window_size.X() - highScoreLabelSize.X()) / 2.0, score_label.GetPosition().Y() + scoreLabelSize.Y() + 10.0})
	es.UIManager.AddElement(high_score_label)

	// UI按钮
	// 让按钮更大一点
	buttonSize := mgl32.Vec2{120.0, 40.0}
	buttonSpacing := float32(20.0)
	totalButtonWidth := buttonSize.X()*2.0 + buttonSpacing

	// 按钮放在右下角，与边缘间隔30像素
	buttonsX := window_size.X() - totalButtonWidth - 30.0
	buttonsY := window_size.Y() - buttonSize.Y() - 30.0
	// Back Button
	backButton := ui.NewUIButton(es.GetContext(),
		"assets/textures/UI/buttons/Back1.png",
		"assets/textures/UI/buttons/Back2.png",
		"assets/textures/UI/buttons/Back3.png",
		mgl32.Vec2{buttonsX, buttonsY},
		buttonSize,
		func() { es.onBackClick() },
	)
	es.UIManager.AddElement(backButton)

	// Restart Button
	restartButton := ui.NewUIButton(es.GetContext(),
		"assets/textures/UI/buttons/Restart1.png",
		"assets/textures/UI/buttons/Restart2.png",
		"assets/textures/UI/buttons/Restart3.png",
		mgl32.Vec2{buttonsX + buttonSize.X() + buttonSpacing, buttonsY},
		buttonSize,
		func() { es.onRestartClick() },
	)
	es.UIManager.AddElement(restartButton)
}

// 返回按钮点击事件
func (es *EndScene) onBackClick() {
	slog.Info("Back button clicked")
	es.SceneManager.RequestReplaceScene(NewTitleScene(es.GetContext(), es.SceneManager, es.sessionData))
}

// 重新开始按钮点击事件
func (es *EndScene) onRestartClick() {
	slog.Info("Restart button clicked")
	es.sessionData.Reset()
	es.SceneManager.RequestReplaceScene(NewGameScene(es.GetContext(), es.SceneManager, es.sessionData))
}
