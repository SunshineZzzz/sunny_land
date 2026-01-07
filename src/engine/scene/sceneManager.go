package scene

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"
)

// 待处理操作
type PendingAction int

const (
	// 没有操作
	PendingActionNone PendingAction = iota
	// 压栈操作
	PendingActionPush
	// 出栈操作
	PendingActionPop
	// 替换栈顶操作
	PendingActionReplace
)

// 场景管理器
type SceneManager struct {
	// 引擎上下文
	context *econtext.Context
	// 场景栈
	sceneStack []IScene
	// 待处理操作
	pendingAction PendingAction
	// 待处理场景
	pendingScene IScene
}

// 创建场景管理器
func NewSceneManager(context *econtext.Context) *SceneManager {
	slog.Debug("create scene manager")
	return &SceneManager{
		context:       context,
		sceneStack:    make([]IScene, 0),
		pendingAction: PendingActionNone,
		pendingScene:  nil,
	}
}

// 清理场景管理器
func (sm *SceneManager) Cleanup() {
	slog.Debug("cleanup scene manager")
	for _, scene := range sm.sceneStack {
		scene.Clean()
	}
}

// 获取当前场景
func (sm *SceneManager) GetCurrentScene() IScene {
	if len(sm.sceneStack) == 0 {
		return nil
	}
	return sm.sceneStack[len(sm.sceneStack)-1]
}

// 更新
func (sm *SceneManager) Update(dt float64) {
	// 只更新当前(栈顶)场景
	currentScene := sm.GetCurrentScene()
	if currentScene != nil {
		currentScene.Update(dt)
	}
	// 执行可能的切换场景操作
	sm.processPendingActions()
}

// 渲染
func (sm *SceneManager) Render() {
	// 渲染时需要叠加渲染所有场景，而不只是栈顶
	for _, scene := range sm.sceneStack {
		scene.Render()
	}
}

// 处理事件
func (sm *SceneManager) HandleInput() {
	// 只处理当前(栈顶)场景的事件
	currentScene := sm.GetCurrentScene()
	if currentScene != nil {
		currentScene.HandleInput()
	}
}

// 请求弹出当前场景
func (sm *SceneManager) RequestPopScene() {
	sm.pendingAction = PendingActionPop
}

// 请求替换当前场景
func (sm *SceneManager) RequestReplaceScene(scene IScene) {
	sm.pendingAction = PendingActionReplace
	sm.pendingScene = scene
}

// 请求压栈场景
func (sm *SceneManager) RequestPushScene(scene IScene) {
	sm.pendingAction = PendingActionPush
	sm.pendingScene = scene
}

// 处理待处理操作
func (sm *SceneManager) processPendingActions() {
	if sm.pendingAction == PendingActionNone {
		return
	}

	switch sm.pendingAction {
	case PendingActionPop:
		sm.popScene()
	case PendingActionReplace:
		sm.replaceScene(sm.pendingScene)
	case PendingActionPush:
		sm.pushScene(sm.pendingScene)
	}

	sm.pendingAction = PendingActionNone
	sm.pendingScene = nil
}

// 弹出场景
func (sm *SceneManager) popScene() {
	if len(sm.sceneStack) <= 0 {
		slog.Warn("pop scene with empty scene stack")
		return
	}

	currentScene := sm.GetCurrentScene()
	if currentScene != nil {
		currentScene.Clean()
	}
	sm.sceneStack[len(sm.sceneStack)-1] = nil
	sm.sceneStack = sm.sceneStack[:len(sm.sceneStack)-1]
	slog.Debug("pop scene", slog.String("scene", currentScene.GetName()))

}

// 替换场景
func (sm *SceneManager) replaceScene(scene IScene) {
	if scene == nil {
		slog.Warn("replace scene with nil scene")
		return
	}

	for i, s := range sm.sceneStack {
		s.Clean()
		sm.sceneStack[i] = nil
	}
	sm.sceneStack = sm.sceneStack[:0]

	if !scene.IsInitialized() {
		scene.Init()
	}

	sm.sceneStack = append(sm.sceneStack, scene)
	slog.Debug("replace scene", slog.String("scene", scene.GetName()))
}

// 压栈场景
func (sm *SceneManager) pushScene(scene IScene) {
	if scene == nil {
		slog.Warn("push scene with nil scene")
		return
	}

	if !scene.IsInitialized() {
		scene.Init()
	}

	sm.sceneStack = append(sm.sceneStack, scene)
	slog.Debug("push scene", slog.String("scene", scene.GetName()))
}
