package core

import (
	"log/slog"

	"sunny_land/src/engine/resource"

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
	// 资源管理器
	resourceManager *resource.ResourceManager
}

// 创建游戏应用
func NewGameApp() *GameApp {
	return &GameApp{
		fpsManager: NewFPS(),
	}
}

// 销毁游戏应用
func (g *GameApp) Destroy() {
	if g.isRunning {
		slog.Warn("game app is running, destroy")
	}

	// 清理资源管理器
	if g.resourceManager != nil {
		g.resourceManager.Clear()
		g.resourceManager = nil
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

	slog.Debug("game app destroy")
}

// 初始化
func (g *GameApp) init() bool {
	if !g.initSDL() {
		return false
	}

	if !g.initTimer() {
		return false
	}

	if !g.initResourceManager() {
		return false
	}

	slog.Debug("game app init")
	g.isRunning = true
	return true
}

// 初始化SDL
func (g *GameApp) initSDL() bool {
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

	slog.Debug("sdl init success")
	return true
}

// 初始化timer
func (g *GameApp) initTimer() bool {
	g.fpsManager.SetTargetFps(144)
	slog.Debug("fps manager init success")
	return true
}

// 初始化资源管理器
func (g *GameApp) initResourceManager() bool {
	g.resourceManager = resource.NewResourceManager(g.sdlRenderer)
	slog.Debug("resource manager init success")
	return true
}

// 运行
func (g *GameApp) Run() {
	slog.Debug("game app run")
	if !g.init() {
		return
	}

	// 测试资源管理器
	g.testResourceManager()

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

// 测试资源管理器
func (g *GameApp) testResourceManager() {
	g.resourceManager.GetTexture("assets/textures/Actors/eagle-attack.png")
	g.resourceManager.GetFont("assets/fonts/VonwaonBitmap-16px.ttf", 16)
	g.resourceManager.GetSound("assets/audio/button_click.wav")

	g.resourceManager.UnloadTexture("assets/textures/Actors/eagle-attack.png")
	g.resourceManager.UnloadFont("assets/fonts/VonwaonBitmap-16px.ttf", 16)
	g.resourceManager.UnloadSound("assets/audio/button_click.wav")
}
