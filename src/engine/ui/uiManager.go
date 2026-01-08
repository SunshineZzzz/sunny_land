package ui

import (
	"log/slog"
	"sunny_land/src/engine/context"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 管理特定场景中的UI元素集合。
 *
 * 负责UI元素的生命周期管理（通过根元素）、渲染调用和输入事件分发。
 * 每个需要UI的场景（如菜单、游戏HUD）应该拥有一个UIManager实例。
 */
type UIManager struct {
	// 一个UIPanel作为根节点(UI元素)
	rootElement *UIPanel
}

// 创建UIManager实例
func NewUIManager() *UIManager {
	slog.Debug("create ui manager")
	return &UIManager{
		rootElement: NewUIPanel(mgl32.Vec2{}, mgl32.Vec2{}, nil),
	}
}

// 初始化
func (um *UIManager) Init(windowSize mgl32.Vec2) bool {
	slog.Debug("init ui manager")
	um.rootElement.SetSize(windowSize)
	return true
}

// 添加一个UI元素到根节点的child_容器中
func (um *UIManager) AddElement(element IUIElement) {
	um.rootElement.AddChild(element)
}

// 获取根UIPanel元素的指针
func (um *UIManager) GetRootElement() *UIPanel {
	return um.rootElement
}

// 清除所有UI元素，通常用于重置UI状态
func (um *UIManager) ClearElements() {
	um.rootElement.RemoveAllChildren()
}

// 处理输入事件，如果事件被处理则返回true。
func (um *UIManager) HandleInput(ctx *context.Context) bool {
	if !um.rootElement.IsVisible() {
		return false
	}
	// 从根元素开始向下分发事件
	return um.rootElement.HandleInput(ctx)
}

// 更新
func (um *UIManager) Update(dt float64, context *context.Context) {
	if !um.rootElement.IsVisible() {
		return
	}
	// 从根元素开始向下更新
	um.rootElement.Update(dt, context)
}

// 渲染
func (um *UIManager) Render(context *context.Context) {
	if !um.rootElement.IsVisible() {
		return
	}
	// 从根元素开始向下渲染
	um.rootElement.Render(context)
}
