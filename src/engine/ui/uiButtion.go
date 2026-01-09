package ui

import (
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/ui/state"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 按钮UI元素
 *
 * 继承自UIInteractive，用于创建可交互的按钮。
 * 支持三种状态：正常、悬停、按下。
 * 支持回调函数，当按钮被点击时调用。
 */
type UIButton struct {
	// 继承UI交互元素基础实现
	UIInteractive
	// 可自定义的函数
	callback func()
}

// 确保UIButton实现了IUIInteractive接口
var _ state.IUIInteractive = (*UIButton)(nil)

// 确保UIButton实现了IUIElement接口
var _ state.IUIElement = (*UIButton)(nil)

/**
 * @brief 构造函数
 * @param normal_sprite_id 正常状态的精灵ID
 * @param hover_sprite_id 悬停状态的精灵ID
 * @param pressed_sprite_id 按下状态的精灵ID
 * @param position 位置
 * @param size 大小
 * @param callback 回调函数
 */
func NewUIButton(context *econtext.Context, normalSpriteId string, hoverSpriteId string,
	pressedSpriteId string, position mgl32.Vec2, size mgl32.Vec2, callback func()) *UIButton {
	ub := &UIButton{
		callback: callback,
	}
	BuildUIInteractive(&ub.UIInteractive, context, position, size)
	ub.AddSprite("normal", render.NewSprite(normalSpriteId, nil, false))
	ub.AddSprite("hover", render.NewSprite(hoverSpriteId, nil, false))
	ub.AddSprite("pressed", render.NewSprite(pressedSpriteId, nil, false))

	// 设置默认状态为"normal"
	ub.SetState(state.NewUINormalState(ub))

	// 设置默认音效
	ub.AddSound("hover", "assets/audio/button_hover.wav")
	ub.AddSound("pressed", "assets/audio/button_click.wav")
	return ub
}

// 设置点击回调函数
func (ub *UIButton) SetClickCallback(callback func()) {
	ub.callback = callback
}

// 返回点击回调函数
func (ub *UIButton) GetClickCallback() func() {
	return ub.callback
}

// 重写基类方法，当按钮被点击时调用回调函数
func (ub *UIButton) Clicked() {
	if ub.callback != nil {
		ub.callback()
	}
}
