package state

import (
	econtext "sunny_land/src/engine/context"

	"github.com/go-gl/mathgl/mgl32"
)

// UI元素抽象
type IUIElement interface {
	// 处理输入事件
	HandleInput(*econtext.Context) bool
	// 更新状态
	Update(float64, *econtext.Context)
	// 渲染
	Render(*econtext.Context)
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
	// 检查给定点是否在元素的边界内
	IsPointInside(mgl32.Vec2) bool
}

// 可交互UI元素的抽象
type IUIInteractive interface {
	// 继承UI元素抽象
	IUIElement
	// 如果有点击事件，则重写该方法
	Clicked()
	// 设置当前显示的精灵
	SetSprite(string)
	// 播放音效
	PlaySound(string)
}

// UI状态的抽象
type IUIState interface {
	// 进入状态
	Enter()
	// 处理输入
	HandleInput(*econtext.Context) IUIState
}

// 基础UI状态
type UIState struct {
	// 指向父节点
	owner IUIInteractive
}
