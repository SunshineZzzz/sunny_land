package core

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/input"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/resource"
	"sunny_land/src/engine/scene"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	// 暂时测试
	testRotation float64 = 0.0
	gameObject           = object.NewGameObject("test", "test")
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
	// 配置
	config *Config
	// 输入管理器
	inputManager *input.InputManager
	// 上下文对象
	context *econtext.Context
	// 场景管理器
	sceneManager *scene.SceneManager
	// 物理引擎
	physicsEngine *physics.PhysicsEngine
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

	// 先关闭场景管理器，确保所有场景都被清理
	g.sceneManager.Cleanup()

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
	if !g.initConfig() || !g.initSDL() || !g.initTimer() ||
		!g.initResourceManager() || !g.initRenderer() ||
		!g.initCamera() || !g.initInputManager() ||
		!g.initPhysicsEngine() || !g.initContext() ||
		!g.initSceneManager() {
		return false
	}

	// 创建第一个场景
	scene := scene.NewGameScene("GameScene", g.context, g.sceneManager)
	// 添加场景到场景管理器
	g.sceneManager.RequestPushScene(scene)

	slog.Debug("game app init")
	g.isRunning = true
	return true
}

// 初始化配置
func (g *GameApp) initConfig() bool {
	g.config = NewConfig("assets/config.json")
	slog.Debug("config init success")
	return true
}

// 初始化SDL
func (g *GameApp) initSDL() bool {
	// 初始化 SDL
	if !sdl.Init(sdl.InitVideo | sdl.InitAudio | sdl.InitEvents) {
		slog.Error("sdl init error", slog.String("error", sdl.GetError()))
		return false
	}

	var windowFlags sdl.WindowFlags
	if g.config.WindowResizable {
		windowFlags |= sdl.WindowResizable
	}
	// 创建窗口与渲染器
	g.sdlWindow = sdl.CreateWindow(g.config.WindowTitle, int32(g.config.WindowWidth), int32(g.config.WindowHeight), windowFlags)
	if g.sdlWindow == nil {
		slog.Error("sdl create window error", slog.String("error", sdl.GetError()))
		return false
	}

	g.sdlRenderer = sdl.CreateRenderer(g.sdlWindow, "")
	if g.sdlRenderer == nil {
		slog.Error("sdl create renderer error", slog.String("error", sdl.GetError()))
		return false
	}

	var vsyncMode int32 = sdl.RendererVSyncDisabled
	if g.config.VsyncEnabled {
		vsyncMode = sdl.RendererVSyncAdaptive
	}
	// 设置VSync，需要注意的是，开启后，驱动程序会尝试将帧率限制到显示器刷新率，有可能会覆盖我们手动设置的targetFps
	if !sdl.SetRenderVSync(g.sdlRenderer, vsyncMode) {
		slog.Warn("sdl set render vsync error", slog.String("error", sdl.GetError()))
	}

	// 设置渲染器的逻辑尺寸
	// sdl.LogicalPresentationLetterbox
	// 它会把游戏画面放大到窗口允许的最大尺寸，同时不改变画面的比例
	// 如果窗口比逻辑画面宽，会看到左右两侧有黑边(Letterbox)
	// 如果窗口比逻辑画面高，会看到顶部和底部有黑边(Letterbox)
	//
	// 设置逻辑分辨率为窗口大小的一半(针对像素游戏)
	if !sdl.SetRenderLogicalPresentation(g.sdlRenderer, int32(g.config.WindowWidth/2), int32(g.config.WindowHeight/2), sdl.LogicalPresentationLetterbox) {
		slog.Error("sdl set render logical presentation error", slog.String("error", sdl.GetError()))
		return false
	}

	slog.Debug("sdl init success")
	return true
}

// 初始化timer
func (g *GameApp) initTimer() bool {
	g.fpsManager.SetTargetFps(g.config.TargetFPS)
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
	g.camera = render.NewCamera(mgl32.Vec2{float32(g.config.WindowWidth / 2), float32(g.config.WindowHeight / 2)}, mgl32.Vec2{0.0, 0.0}, nil)
	slog.Debug("camera init success")
	return true
}

// 初始化输入管理器
func (g *GameApp) initInputManager() bool {
	g.inputManager = input.NewInputManager(g.sdlRenderer, &g.config.InputMappings)
	slog.Debug("input manager init success")
	return true
}

// 初始化上下文对象
func (g *GameApp) initContext() bool {
	g.context = econtext.NewContext(g.inputManager, g.renderer, g.resourceManager, g.camera, g.physicsEngine)
	slog.Debug("context init success")
	return true
}

// 初始化场景管理器
func (g *GameApp) initSceneManager() bool {
	g.sceneManager = scene.NewSceneManager(g.context)
	slog.Debug("scene manager init success")
	return true
}

// 初始化物理引擎
func (g *GameApp) initPhysicsEngine() bool {
	g.physicsEngine = physics.NewPhysicsEngine()
	slog.Debug("physics engine init success")
	return true
}

// 运行
func (g *GameApp) Run() {
	slog.Debug("game app run")
	if !g.init() {
		return
	}

	// 测试资源管理器
	// g.testResourceManager()
	// 测试游戏对象
	// g.testGameObject()

	for g.isRunning {
		g.fpsManager.Update()
		deltaTime := g.fpsManager.GetDeltaTime()
		// fmt.Printf("dt: %f\n", 1.0/deltaTime)
		// 每帧首先更新输入管理器
		g.inputManager.Update()
		g.HandleEvents()
		g.update(deltaTime)
		g.render()
	}

	g.Destroy()
}

// 处理事件
func (g *GameApp) HandleEvents() {
	if g.inputManager.ShouldQuit() {
		slog.Debug("received quit event, exiting")
		g.isRunning = false
	}

	// 测试输入管理器
	// g.testInputManager()
	// 处理场景事件
	g.sceneManager.HandleInput()
}

// 更新
func (g *GameApp) update(dt float64) {
	// 更新场景
	g.sceneManager.Update(dt)
}

// 渲染
func (g *GameApp) render() {
	// 清除屏幕
	g.renderer.ClearScreen()

	// 渲染代码
	// 测试代码
	// g.testRenderer()
	// gameObject.Render(g.context)
	// 渲染场景
	g.sceneManager.Render()

	// 显示渲染结果
	g.renderer.Present()
}

// 测试资源管理器
// func (g *GameApp) testResourceManager() {
// 	g.resourceManager.GetTexture("assets/textures/Actors/eagle-attack.png")
// 	g.resourceManager.GetFont("assets/fonts/VonwaonBitmap-16px.ttf", 16)
// 	g.resourceManager.GetSound("assets/audio/button_click.wav")

// 	g.resourceManager.UnloadTexture("assets/textures/Actors/eagle-attack.png")
// 	g.resourceManager.UnloadFont("assets/fonts/VonwaonBitmap-16px.ttf", 16)
// 	g.resourceManager.UnloadSound("assets/audio/button_click.wav")
// }

// 测试渲染器
// func (g *GameApp) testRenderer() {
// 	spriteWorld := render.NewSprite("assets/textures/Actors/frog.png", nil, false)
// 	spriteUI := render.NewSprite("assets/textures/UI/buttons/Start1.png", nil, false)
// 	spriteParallax := render.NewSprite("assets/textures/Layers/back.png", nil, false)

// 	testRotation += 0.1

// 	// 注意渲染顺序
// 	g.renderer.DrawSpriteWithParallax(g.camera, spriteParallax, mgl32.Vec2{100, 100}, mgl32.Vec2{0.5, 0.5}, mgl32.Vec2{1.0, 1.0}, math.Vec2B{true, false})
// 	g.renderer.DrawSprite(g.camera, spriteWorld, mgl32.Vec2{200, 200}, mgl32.Vec2{1.0, 1.0}, testRotation)
// 	g.renderer.DrawUISprite(spriteUI, mgl32.Vec2{100, 100}, nil)
// }

// 测试输入管理器
// func (g *GameApp) testInputManager() {
// 	actions := []string{
// 		"move_up",
// 		"move_down",
// 		"move_left",
// 		"move_right",
// 		"jump",
// 		"attack",
// 		"pause",
// 		"MouseLeftClick",
// 		"MouseRightClick",
// 	}

// 	for _, action := range actions {
// 		if g.inputManager.IsActionPressed(action) {
// 			slog.Info("按下 ", slog.String("action", action))
// 		}
// 		if g.inputManager.IsActionReleased(action) {
// 			slog.Info("抬起 ", slog.String("action", action))
// 		}
// 		if g.inputManager.IsActionDown(action) {
// 			slog.Info("按下中 ", slog.String("action", action))
// 		}
// 	}
// }

// func (g *GameApp) testGameObject() {
// 	transformComp := component.NewTransformComponent(mgl32.Vec2{100.0, 100.0}, mgl32.Vec2{1.0, 1.0}, 0.0)
// 	spriteComp := component.NewSpriteComponent("assets/textures/Props/big-crate.png", g.resourceManager, utils.AlignCenter, nil, false)
// 	gameObject.AddComponent(transformComp)
// 	gameObject.AddComponent(spriteComp)

// 	transformComp2 := gameObject.GetComponent(&component.TransformComponent{}).(*component.TransformComponent)
// 	if transformComp2 != transformComp {
// 		slog.Error("testGameObject: transform component mismatch")
// 	}
// 	transformComp2.SetScale(mgl32.Vec2{2.0, 2.0})
// 	transformComp2.SetRotation(30.0)
// }
