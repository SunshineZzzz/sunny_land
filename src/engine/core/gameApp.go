package core

import (
	"log/slog"

	"sunny_land/src/engine/render"
	"sunny_land/src/engine/resource"
	"sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// 逻辑画面宽度
	LogicWidth = int32(640)
	// 逻辑画面高度
	LogicHeight = int32(360)
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
	// 渲染器
	renderer *render.Renderer
	// 摄像机
	camera *render.Camera
	// 暂时测试
	testRotation float64
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
	if !g.initSDL() || !g.initTimer() || !g.initResourceManager() ||
		!g.initRenderer() || !g.initCamera() {
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
	if !sdl.CreateWindowAndRenderer("SunnyLand", LogicWidth, LogicHeight, sdl.WindowResizable, &g.sdlWindow, &g.sdlRenderer) {
		slog.Error("sdl create window and renderer error", slog.String("error", sdl.GetError()))
		return false
	}

	// 设置渲染器的逻辑尺寸
	// sdl.LogicalPresentationLetterbox
	// 它会把游戏画面放大到窗口允许的最大尺寸，同时不改变画面的比例
	// 如果窗口比逻辑画面宽，会看到左右两侧有黑边(Letterbox)
	// 如果窗口比逻辑画面高，会看到顶部和底部有黑边(Letterbox)
	if !sdl.SetRenderLogicalPresentation(g.sdlRenderer, LogicWidth, LogicHeight, sdl.LogicalPresentationLetterbox) {
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

// 初始化渲染器
func (g *GameApp) initRenderer() bool {
	g.renderer = render.NewRenderer(g.sdlRenderer, g.resourceManager)
	slog.Debug("renderer init success")
	return true
}

// 初始化摄像机
func (g *GameApp) initCamera() bool {
	g.camera = render.NewCamera(mgl32.Vec2{float32(LogicWidth), float32(LogicHeight)}, mgl32.Vec2{0.0, 0.0}, nil)
	slog.Debug("camera init success")
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
func (g *GameApp) update( /*deltaTime*/ float64) {
	g.testCamera()
}

// 渲染
func (g *GameApp) render() {
	// 清除屏幕
	g.renderer.ClearScreen()

	// 渲染代码
	g.testRenderer()

	// 显示渲染结果
	g.renderer.Present()
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

// 测试渲染器
func (g *GameApp) testRenderer() {
	spriteWorld := render.NewSprite("assets/textures/Actors/frog.png", nil, false)
	spriteUI := render.NewSprite("assets/textures/UI/buttons/Start1.png", nil, false)
	spriteParallax := render.NewSprite("assets/textures/Layers/back.png", nil, false)

	g.testRotation += 0.1

	// 注意渲染顺序
	g.renderer.DrawSpriteWithParallax(g.camera, spriteParallax, mgl32.Vec2{100, 100}, mgl32.Vec2{0.5, 0.5}, mgl32.Vec2{1.0, 1.0}, math.Vec2B{true, false})
	g.renderer.DrawSprite(g.camera, spriteWorld, mgl32.Vec2{200, 200}, mgl32.Vec2{1.0, 1.0}, g.testRotation)
	g.renderer.DrawUISprite(spriteUI, mgl32.Vec2{100, 100}, nil)
}

// 测试摄像机
func (g *GameApp) testCamera() {
	key_state := sdl.GetKeyboardState()
	if key_state[sdl.ScancodeW] {
		// 摄像机向上运动
		g.camera.Move(mgl32.Vec2{0, -1})
	}
	if key_state[sdl.ScancodeS] {
		// 摄像机向下运动
		g.camera.Move(mgl32.Vec2{0, 1})
	}
	if key_state[sdl.ScancodeA] {
		// 摄像机向左运动
		g.camera.Move(mgl32.Vec2{-1, 0})
	}
	if key_state[sdl.ScancodeD] {
		// 摄像机向右运动
		g.camera.Move(mgl32.Vec2{1, 0})
	}
}
