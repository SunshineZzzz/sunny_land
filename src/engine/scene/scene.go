package scene

import (
	"container/list"
	"log/slog"

	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/resource"
)

// 场景接口，负责管理场景中的游戏对象和场景生命周期
type IScene interface {
	// 初始化场景
	Init()
	// 更新场景状态
	Update(float64)
	// 渲染场景
	Render()
	// 处理输入事件
	HandleInput()
	// 清理场景
	Clean()
	// 直接向场景中添加一个游戏对象，初始化时可用，游戏进行中不安全
	AddGameObject(*object.GameObject)
	// 安全地添加游戏对象，添加到pending_additions中
	SafeAddGameObject(*object.GameObject)
	// 直接从场景中移除一个游戏对象，一般不使用，但保留实现的逻辑
	RemoveGameObject(*object.GameObject)
	// 安全地移除游戏对象，设置need_remove_标记
	SafeRemoveGameObject(*object.GameObject)
	// 获取场景名称
	GetName() string
	// 判断场景是否已初始化
	IsInitialized() bool
	// 获取资源管理器
	GetResourceManager() *resource.ResourceManager
}

// 基础场景
type scene struct {
	// 场景名称
	sceneName string
	// 上下文
	ctx *econtext.Context
	// 场景管理器
	sceneManager *SceneManager
	// 是否初始化
	initialized bool
	// 场景中的游戏对象容器
	gameObjects *list.List
	// 待添加的游戏对象容器，延迟添加
	pendingAdditions []*object.GameObject
}

// 确保实现了IScene接口
var _ IScene = (*scene)(nil)

// 构建场景
func buildScene(s *scene, sceneName string, ctx *econtext.Context, sceneManager *SceneManager) {
	s.sceneName = sceneName
	s.ctx = ctx
	s.sceneManager = sceneManager
	s.initialized = false
	s.gameObjects = list.New()
	s.pendingAdditions = make([]*object.GameObject, 0)
}

// 初始化场景
func (s *scene) Init() {
	if s.initialized {
		slog.Warn("Scene already initialized", slog.String("sceneName", s.sceneName))
		return
	}
	s.initialized = true
}

// 更新场景状态
func (s *scene) Update(dt float64) {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}

	// 先更新物理引擎
	s.ctx.PhysicsEngine.Update(dt)

	// 更新所有游戏对象，并删除需要移除的对象
	for e := s.gameObjects.Front(); e != nil; {
		next := e.Next()

		gt := e.Value.(*object.GameObject)
		if gt.NeedRemove() {
			s.gameObjects.Remove(e)
			gt.Clean()
		} else {
			gt.Update(dt, s.ctx)
		}

		e = next
	}

	// 处理待添加(延时添加)的游戏对象
	s.processPendingAdditions()
}

// 处理待添加(延时添加)的游戏对象
func (s *scene) processPendingAdditions() {
	for _, gt := range s.pendingAdditions {
		s.gameObjects.PushBack(gt)
	}
	s.pendingAdditions = make([]*object.GameObject, 0)
}

// 渲染场景
func (s *scene) Render() {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}
	// 渲染所有游戏对象
	for e := s.gameObjects.Front(); e != nil; e = e.Next() {
		gt := e.Value.(*object.GameObject)
		gt.Render(s.ctx)
	}
}

// 处理输入事件
func (s *scene) HandleInput() {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}

	// 处理所有游戏对象的输入事件, 并删除需要移除的对象
	for e := s.gameObjects.Front(); e != nil; {
		next := e.Next()

		gt := e.Value.(*object.GameObject)
		if gt.NeedRemove() {
			s.gameObjects.Remove(e)
			gt.Clean()
		} else {
			gt.HandleInput(s.ctx)
		}

		e = next
	}
}

// 清理场景
func (s *scene) Clean() {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}
	s.initialized = false

	// 清理所有游戏对象
	for e := s.gameObjects.Front(); e != nil; e = e.Next() {
		gt := e.Value.(*object.GameObject)
		gt.Clean()
	}
	s.gameObjects.Init()
	slog.Debug("Scene cleaned", slog.String("sceneName", s.sceneName))
}

// 直接向场景中添加一个游戏对象，初始化时可用，游戏进行中不安全
func (s *scene) AddGameObject(gt *object.GameObject) {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}
	if gt == nil {
		slog.Warn("GameObject is nil", slog.String("sceneName", s.sceneName))
		return
	}
	s.gameObjects.PushBack(gt)
}

// 安全地添加游戏对象，添加到pending_additions中
func (s *scene) SafeAddGameObject(gt *object.GameObject) {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}
	if gt == nil {
		slog.Warn("GameObject is nil", slog.String("sceneName", s.sceneName))
		return
	}
	s.pendingAdditions = append(s.pendingAdditions, gt)
}

// 直接从场景中移除一个游戏对象，一般不使用，但保留实现的逻辑
func (s *scene) RemoveGameObject(gt *object.GameObject) {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return
	}
	if gt == nil {
		slog.Warn("GameObject is nil", slog.String("sceneName", s.sceneName))
		return
	}

	for e := s.gameObjects.Front(); e != nil; e = e.Next() {
		if e.Value.(*object.GameObject) == gt {
			s.gameObjects.Remove(e)
			gt.Clean()
			slog.Debug("GameObject removed", slog.String("sceneName", s.sceneName), slog.String("gameObjectName", gt.GetName()))
			return
		}
	}
	slog.Warn("GameObject not found in scene", slog.String("sceneName", s.sceneName), slog.String("gameObjectName", gt.GetName()))
}

// 安全地移除游戏对象，设置need_remove_标记
func (s *scene) SafeRemoveGameObject(gt *object.GameObject) {
	gt.SetNeedRemove(true)
}

// 获取场景名称
func (s *scene) GetName() string {
	return s.sceneName
}

// 判断场景是否已初始化
func (s *scene) IsInitialized() bool {
	return s.initialized
}

// 根据名称查找游戏对象
func (s *scene) FindGameObjectByName(name string) *object.GameObject {
	if !s.initialized {
		slog.Warn("Scene not initialized", slog.String("sceneName", s.sceneName))
		return nil
	}
	if name == "" {
		slog.Warn("GameObject name is empty", slog.String("sceneName", s.sceneName))
		return nil
	}

	for e := s.gameObjects.Front(); e != nil; e = e.Next() {
		gt := e.Value.(*object.GameObject)
		if gt.GetName() == name {
			return gt
		}
	}
	slog.Warn("GameObject not found in scene", slog.String("sceneName", s.sceneName), slog.String("gameObjectName", name))
	return nil
}

// 获取资源管理器
func (s *scene) GetResourceManager() *resource.ResourceManager {
	return s.ctx.ResourceManager
}
