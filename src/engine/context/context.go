package context

import (
	"log/slog"
	"sunny_land/src/engine/input"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/resource"
)

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
}

// 确保实现IContext接口
var _ physics.IContext = (*Context)(nil)

// 创建上下文对象
func NewContext(inputManager *input.InputManager, renderer *render.Renderer,
	resourceManager *resource.ResourceManager, camera *render.Camera,
	physicsEngine *physics.PhysicsEngine) *Context {
	slog.Debug("create context")
	return &Context{
		InputManager:    inputManager,
		Renderer:        renderer,
		ResourceManager: resourceManager,
		Camera:          camera,
		PhysicsEngine:   physicsEngine,
	}
}

// 获取渲染器
func (c *Context) GetRenderer() *render.Renderer {
	return c.Renderer
}

// 获取摄像机
func (c *Context) GetCamera() *render.Camera {
	return c.Camera
}
