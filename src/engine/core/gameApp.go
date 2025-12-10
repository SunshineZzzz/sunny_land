package core

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 主游戏应用，初始化SDL，管理游戏循环
type GameApp struct {
	// SDL窗口
	sdlWindow *sdl.Window
	// SDL渲染器
	sdlRenderer *sdl.Renderer
	// 是否运行中
	isRunning bool
	// 帧率管理器
	fpsManager *FPS
}

// 创建游戏应用
func NewGameApp() *GameApp {
	return &GameApp{
		fpsManager: NewFPS(),
	}
}

// 销毁游戏应用
func (g *GameApp) Destroy() {
	slog.Debug("game app destroy")
	if g.isRunning {
		slog.Warn("game app is running, destroy")
	}

	// 清理SDL资源
	if g.sdlRenderer != nil {
		sdl.DestroyRenderer(g.sdlRenderer)
		g.sdlRenderer = nil
	}
	if g.sdlWindow != nil {
		sdl.DestroyWindow(g.sdlWindow)
		g.sdlWindow = nil
	}
	sdl.Quit()
}

// 初始化
func (g *GameApp) init() bool {
	slog.Debug("game app init")

	// 初始化 SDL
	if !sdl.Init(sdl.InitVideo | sdl.InitAudio | sdl.InitEvents) {
		slog.Error("sdl init error", slog.String("error", sdl.GetError()))
		return false
	}

	// 创建窗口与渲染器
	if !sdl.CreateWindowAndRenderer("SunnyLand", 1280, 720, sdl.WindowResizable, &g.sdlWindow, &g.sdlRenderer) {
		slog.Error("sdl create window and renderer error", slog.String("error", sdl.GetError()))
		return false
	}

	// 设置渲染器的逻辑尺寸
	// sdl.LogicalPresentationLetterbox
	// 它会把游戏画面放大到窗口允许的最大尺寸，同时不改变画面的比例
	// 如果窗口比逻辑画面宽，会看到左右两侧有黑边(Letterbox)
	// 如果窗口比逻辑画面高，会看到顶部和底部有黑边(Letterbox)
	if !sdl.SetRenderLogicalPresentation(g.sdlRenderer, 1280, 720, sdl.LogicalPresentationLetterbox) {
		slog.Error("sdl set render logical presentation error", slog.String("error", sdl.GetError()))
		return false
	}

	g.isRunning = true
	return true
}

// 运行
func (g *GameApp) Run() {
	slog.Debug("game app run")
	if !g.init() {
		return
	}

	g.fpsManager.SetTargetFps(60)
	for g.isRunning {
		g.fpsManager.Update()
		deltaTime := g.fpsManager.GetDeltaTime()
		// fmt.Printf("dt: %f\n", 1.0/deltaTime)
		g.handleEvents()
		g.update(deltaTime)
		g.render()
	}

	g.Destroy()
}

// 处理事件
func (g *GameApp) handleEvents() {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		if event.Type() == sdl.EventQuit {
			g.isRunning = false
			return
		}
	}
}

// 更新
func (g *GameApp) update(deltaTime float64) {
}

// 渲染
func (g *GameApp) render() {
}
