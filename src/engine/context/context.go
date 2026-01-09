package context

import (
	"log/slog"

	"sunny_land/src/engine/audio"
	"sunny_land/src/engine/input"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/resource"

	"github.com/go-gl/mathgl/mgl32"
)

type GameStateType int

const (
	// 标题界面
	GameStateTitle GameStateType = iota
	// 游戏进行中
	GameStatePlaying
	// 游戏暂停
	GameStatePaused
	// 游戏结束
	GameStateGameOver
)

// 游戏状态抽象
type IGameState interface {
	// 判断是否在游戏进行中
	IsPlaying() bool
	// 设置当前游戏状态
	SetState(GameStateType)
	// 获取逻辑分辨率
	GetLogicalSize() mgl32.Vec2
}

// 持有对核心引擎模块引用的上下文对象，用于简化依赖注入，
// 传递Context对象即可获取引擎的各个模块。
type Context struct {
	// 输入管理器
	InputManager *input.InputManager
	// 渲染器
	Renderer *render.Renderer
	// 资源管理器
	ResourceManager *resource.ResourceManager
	// 相机
	Camera *render.Camera
	// 物理引擎
	PhysicsEngine *physics.PhysicsEngine
	// 音频播放器
	AudioPlayer *audio.AudioPlayer
	// 文本渲染器
	TextRenderer *render.TextRenderer
	// 游戏状态
	GameState IGameState
}

// 确保实现IContext接口
var _ physics.IContext = (*Context)(nil)

// 创建上下文对象
func NewContext(inputManager *input.InputManager, renderer *render.Renderer,
	resourceManager *resource.ResourceManager, camera *render.Camera,
	physicsEngine *physics.PhysicsEngine, audioPlayer *audio.AudioPlayer,
	textRenderer *render.TextRenderer, gameState IGameState) *Context {
	slog.Debug("create context")
	return &Context{
		InputManager:    inputManager,
		Renderer:        renderer,
		ResourceManager: resourceManager,
		Camera:          camera,
		PhysicsEngine:   physicsEngine,
		AudioPlayer:     audioPlayer,
		TextRenderer:    textRenderer,
		GameState:       gameState,
	}
}

// 获取渲染器
func (c *Context) GetRenderer() physics.IRenderer {
	return c.Renderer
}

// 获取摄像机
func (c *Context) GetCamera() physics.ICamera {
	return c.Camera
}

// 获取输入管理器
func (c *Context) GetInputManager() *input.InputManager {
	return c.InputManager
}

// 获取物理引擎
func (c *Context) GetPhysicsEngine() physics.PhysicsEngine {
	return *c.PhysicsEngine
}

// 获取音频播放器
func (c *Context) GetAudioPlayer() *audio.AudioPlayer {
	return c.AudioPlayer
}

// 获取文本渲染器
func (c *Context) GetTextRenderer() *render.TextRenderer {
	return c.TextRenderer
}

// 获取资源管理器
func (c *Context) GetResourceManager() *resource.ResourceManager {
	return c.ResourceManager
}

// 获取游戏状态
func (c *Context) GetGameState() IGameState {
	return c.GameState
}
