package ui

import (
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/ui/state"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 用于分组其他UI元素的容器UI元素
 *
 * Panel通常用于布局和组织。
 * 可以选择是否绘制背景色(纯色)。
 */
type UIPanel struct {
	// 继承UI元素基础实现
	UIElement
	// 背景颜色
	backgroundColor *emath.FColor
}

// 确保UIPanel实现IUIElement接口
var _ state.IUIElement = (*UIPanel)(nil)

// 创建UIPanel实例
func NewUIPanel(position mgl32.Vec2, size mgl32.Vec2, backgroundColor *emath.FColor) *UIPanel {
	up := &UIPanel{
		backgroundColor: backgroundColor,
	}
	BuildUIElement(&up.UIElement, position, size)
	return up
}

// 设置背景颜色
func (p *UIPanel) SetBackgroundColor(color *emath.FColor) {
	p.backgroundColor = color
}

// 获取背景颜色
func (p *UIPanel) GetBackgroundColor() *emath.FColor {
	return p.backgroundColor
}

// 渲染
func (p *UIPanel) Render(ctx *econtext.Context) {
	if !p.visible {
		return
	}

	if p.backgroundColor != nil {
		ctx.GetRenderer().DrawUIFilledRect(p.GetBounds(), *p.backgroundColor)
	}

	// 调用基类渲染方法(绘制子节点)
	p.UIElement.Render(ctx)
}
