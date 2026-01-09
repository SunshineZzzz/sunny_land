package ui

import (
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/render"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief UILabel 类用于创建和管理用户界面中的文本标签
 *
 * UILabel 继承自 UIElement，提供了文本渲染功能。
 * 它可以设置文本内容、字体ID、字体大小和文本颜色。
 *
 * @note 需要一个文本渲染器来获取和更新文本尺寸。
 */
type UILabel struct {
	// 继承UI元素基础实现
	UIElement
	// 需要文本渲染器，用于获取/更新文本尺寸
	textRenderer *render.TextRenderer
	// 文本内容
	text string
	// 字体ID
	fontID string
	// 字体大小
	fontSize int
	// 文本颜色
	textColor emath.FColor
}

/**
 * @brief 构造一个UILabel
 *
 * @param text_renderer 文本渲染器
 * @param text 文本内容
 * @param font_id 字体ID
 * @param font_size 字体大小
 * @param text_color 文本颜色
 * @param position 标签的局部位置
 */
func NewUILabel(textRenderer *render.TextRenderer, text string, fontID string, fontSize int, textColor emath.FColor, position mgl32.Vec2) *UILabel {
	ui := &UILabel{
		textRenderer: textRenderer,
		text:         text,
		fontID:       fontID,
		fontSize:     fontSize,
		textColor:    textColor,
	}
	BuildUIElement(&ui.UIElement, position, textRenderer.GetTextSize(text, fontID, fontSize))
	return ui
}

// 获取标签的文本内容
func (ui *UILabel) GetText() string {
	return ui.text
}

// 设置标签的文本内容
func (ui *UILabel) SetText(text string) {
	ui.text = text
	ui.SetSize(ui.textRenderer.GetTextSize(text, ui.fontID, ui.fontSize))
}

// 获取标签的字体ID
func (ui *UILabel) GetFontID() string {
	return ui.fontID
}

// 设置标签的字体ID
func (ui *UILabel) SetFontID(fontID string) {
	ui.fontID = fontID
	ui.SetSize(ui.textRenderer.GetTextSize(ui.text, fontID, ui.fontSize))
}

// 获取标签的字体大小
func (ui *UILabel) GetFontSize() int {
	return ui.fontSize
}

// 设置标签的字体大小
func (ui *UILabel) SetFontSize(fontSize int) {
	ui.fontSize = fontSize
	ui.SetSize(ui.textRenderer.GetTextSize(ui.text, ui.fontID, fontSize))
}

// 获取标签的文本颜色
func (ui *UILabel) GetTextColor() emath.FColor {
	return ui.textColor
}

// 设置标签的文本颜色
func (ui *UILabel) SetTextColor(textColor emath.FColor) {
	ui.textColor = textColor
}

// 渲染标签
func (ui *UILabel) Render(ctx *econtext.Context) {
	if !ui.visible || ui.text == "" {
		return
	}

	ui.textRenderer.DrawUIText(ui.text, ui.fontID, ui.fontSize, ui.GetScreenPosition(), ui.textColor)

	// 渲染子元素（调用基类方法）
	ui.UIElement.Render(ctx)
}
