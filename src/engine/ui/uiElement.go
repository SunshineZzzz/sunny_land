package ui

import (
	"container/list"

	"sunny_land/src/engine/context"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// UI元素抽象
type IUIElement interface {
	// 处理输入事件
	HandleInput(*context.Context) bool
	// 更新状态
	Update(float64, *context.Context)
	// 渲染
	Render(*context.Context)
	// 是否需要移除
	IsNeedRemove() bool
	// 设置父元素
	SetParent(parent IUIElement)
	// 获取父元素
	GetParent() IUIElement
	// 添加子元素
	AddChild(child IUIElement)
	// 获取(计算)元素在屏幕上位置, 相对于屏幕左上角
	GetScreenPosition() mgl32.Vec2
}

// UI元素基础实现
type UIElement struct {
	// 相对于父元素的局部位置
	position mgl32.Vec2
	// 元素大小
	size mgl32.Vec2
	// 元素当前是否可见
	visible bool
	// 是否需要移除
	needRemove bool
	// 指向父节点的非拥有指针
	parent IUIElement
	// 子元素列表
	children list.List
}

// 确保UIElement实现了IUIElement接口
var _ IUIElement = (*UIElement)(nil)

// 构建UI元素
func BuildUIElement(e *UIElement, position mgl32.Vec2, size mgl32.Vec2) {
	e.position = position
	e.size = size
	e.visible = true
	e.needRemove = false
	e.parent = nil
	e.children.Init()
}

// 处理输入事件
func (e *UIElement) HandleInput(ctx *context.Context) bool {
	if !e.visible {
		return false
	}

	// 遍历所有子节点，并删除标记了移除的元素
	for child := e.children.Front(); child != nil; {
		next := child.Next()

		uiElement := child.Value.(IUIElement)
		if !uiElement.IsNeedRemove() {
			if uiElement.HandleInput(ctx) {
				return true
			}
		} else {
			e.children.Remove(child)
		}

		child = next
	}
	return false
}

// 更新状态
func (e *UIElement) Update(dt float64, ctx *context.Context) {
	if !e.visible {
		return
	}

	// 遍历所有子节点，并删除标记了移除的元素
	for child := e.children.Front(); child != nil; {
		next := child.Next()

		uiElement := child.Value.(IUIElement)
		if !uiElement.IsNeedRemove() {
			uiElement.Update(dt, ctx)
		} else {
			e.children.Remove(child)
		}

		child = next
	}
}

// 渲染
func (e *UIElement) Render(ctx *context.Context) {
	if !e.visible {
		return
	}

	// 渲染子元素
	for child := e.children.Front(); child != nil; child = child.Next() {
		uiElement := child.Value.(IUIElement)
		uiElement.Render(ctx)
	}
}

// 是否需要移除
func (e *UIElement) IsNeedRemove() bool {
	return e.needRemove
}

// 添加子元素
func (e *UIElement) AddChild(child IUIElement) {
	if child == nil {
		return
	}
	child.SetParent(e)
	e.children.PushBack(child)
}

// 设置父元素
func (e *UIElement) SetParent(parent IUIElement) {
	e.parent = parent
}

// 将指定子元素从列表中移除，并返回其指针
func (e *UIElement) RemoveChild(child IUIElement) IUIElement {
	if child == nil {
		return nil
	}
	for element := e.children.Front(); element != nil; element = element.Next() {
		if element.Value.(IUIElement) == child {
			e.children.Remove(element)
			// 清除父指针
			child.SetParent(nil)
			// 返回被移除的子元素（可以挂载到别处）
			return child
		}
	}
	// 未找到子元素
	return nil
}

// 移除所有子元素
func (e *UIElement) RemoveAllChildren() {
	for child := e.children.Front(); child != nil; child = child.Next() {
		// 清除父指针
		child.Value.(IUIElement).SetParent(nil)
	}
	e.children.Init()
}

// 获取元素大小
func (e *UIElement) GetSize() mgl32.Vec2 {
	return e.size
}

// 获取元素位置, 相对于父元素
func (e *UIElement) GetPosition() mgl32.Vec2 {
	return e.position
}

// 是否可见
func (e *UIElement) IsVisible() bool {
	return e.visible
}

// 获取父元素
func (e *UIElement) GetParent() IUIElement {
	return e.parent
}

// 获取子元素列表
func (e *UIElement) GetChildren() *list.List {
	return &e.children
}

// 设置元素大小
func (e *UIElement) SetSize(size mgl32.Vec2) {
	e.size = size
}

// 设置元素可见性
func (e *UIElement) SetVisible(visible bool) {
	e.visible = visible
}

// 设置元素位置, 相对于父元素
func (e *UIElement) SetPosition(position mgl32.Vec2) {
	e.position = position
}

// 获取(计算)元素在屏幕上位置, 相对于屏幕左上角
func (e *UIElement) GetScreenPosition() mgl32.Vec2 {
	// 递归计算父元素的屏幕位置
	if e.parent != nil {
		return e.parent.GetScreenPosition().Add(e.position)
	}
	// 根元素的位置已经是相对屏幕的绝对位置
	return e.position
}

// 获取(计算)元素的边界(屏幕坐标)
func (e *UIElement) GetBounds() emath.Rect {
	// 元素边界是相对于屏幕的
	return emath.Rect{
		Position: e.GetScreenPosition(),
		Size:     e.size,
	}
}
