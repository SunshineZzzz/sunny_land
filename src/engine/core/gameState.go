package core

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 管理和查询游戏的全局宏观状态。
 *
 * 提供一个中心点来确定游戏当前处于哪个主要模式，
 * 以便其他系统（输入、渲染、更新等）可以相应地调整其行为。
 */
type GameState struct {
	// SDL窗口，用于获取窗口大小
	sdlWindow *sdl.Window
	// SDL渲染器，用于获取逻辑分辨率
	sdlRenderer *sdl.Renderer
	// 当前游戏状态
	currentState econtext.GameStateType
}

/**
 * @brief 构造函数，初始化游戏状态。
 * @param window SDL窗口，必须传入有效值。
 * @param renderer SDL渲染器，必须传入有效值。
 * @param initial_state 游戏的初始状态，默认为 Title
 */
func NewGameState(window *sdl.Window, renderer *sdl.Renderer, initialState econtext.GameStateType) *GameState {
	if window == nil || renderer == nil {
		slog.Error("window or renderer is nil")
		return nil
	}
	slog.Debug("new game state")
	return &GameState{
		sdlWindow:    window,
		sdlRenderer:  renderer,
		currentState: initialState,
	}
}

// 确保实现IGameState接口
var _ econtext.IGameState = (*GameState)(nil)

// 获取当前游戏状态
func (gs *GameState) GetState() econtext.GameStateType {
	return gs.currentState
}

// 设置当前游戏状态
func (gs *GameState) SetState(newState econtext.GameStateType) {
	if gs.currentState != newState {
		gs.currentState = newState
	}
}

// 获取窗口大小
func (gs *GameState) GetWindowSize() mgl32.Vec2 {
	var w, h int32
	sdl.GetWindowSize(gs.sdlWindow, &w, &h)
	return mgl32.Vec2{float32(w), float32(h)}
}

// 设置窗口大小
func (gs *GameState) SetWindowSize(size mgl32.Vec2) {
	sdl.SetWindowSize(gs.sdlWindow, int32(size.X()), int32(size.Y()))
}

// 获取逻辑分辨率
func (gs *GameState) GetLogicalSize() mgl32.Vec2 {
	var w, h int32
	sdl.GetRenderLogicalPresentation(gs.sdlRenderer, &w, &h, nil)
	return mgl32.Vec2{float32(w), float32(h)}
}

// 设置逻辑分辨率
func (gs *GameState) SetLogicalSize(size mgl32.Vec2) {
	sdl.SetRenderLogicalPresentation(gs.sdlRenderer, int32(size.X()), int32(size.Y()), sdl.LogicalPresentationLetterbox)
}

// 判断是否在标题界面
func (gs *GameState) IsInTitle() bool {
	return gs.currentState == econtext.GameStateTitle
}

// 判断是否在游戏进行中
func (gs *GameState) IsPlaying() bool {
	return gs.currentState == econtext.GameStatePlaying
}

// 判断是否在游戏暂停
func (gs *GameState) IsPaused() bool {
	return gs.currentState == econtext.GameStatePaused
}

// 判断是否在游戏结束
func (gs *GameState) IsGameOver() bool {
	return gs.currentState == econtext.GameStateGameOver
}
