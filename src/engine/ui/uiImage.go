package ui

import (
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/ui/state"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 一个用于显示纹理或部分纹理的UI元素。
 *
 * 继承自UIElement并添加了渲染图像的功能。
 */
type UIImage struct {
	// 继承UI元素基础实现
	UIElement
	// 精灵图对象
	sprite *render.Sprite
}

// 确保UIImage实现IUIElement接口
var _ state.IUIElement = (*UIImage)(nil)

/**
 * @brief 构造一个UIImage对象。
 *
 * @param textureId 要显示的纹理ID。
 * @param position 图像的局部位置。
 * @param size 图像元素的大小。（如果为{0,0}，则使用纹理的原始尺寸）
 * @param sourceRect 可选：要绘制的纹理部分。（如果为空，则使用纹理的整个区域）
 * @param isFlipped 可选：精灵是否应该水平翻转。
 */
func NewUIImage(textureId string, position mgl32.Vec2, size mgl32.Vec2, sourceRect *sdl.FRect, isFlipped bool) *UIImage {
	ui := &UIImage{
		sprite: render.NewSprite(textureId, sourceRect, isFlipped),
	}
	BuildUIElement(&ui.UIElement, position, size)
	return ui
}

// 渲染
func (i *UIImage) Render(ctx *econtext.Context) {
	if !i.visible || i.sprite == nil || i.sprite.GetTextureId() == "" {
		// 如果不可见或没有分配纹理则不渲染
		return
	}

	// 渲染自身
	position := i.GetScreenPosition()
	// 如果尺寸为0，则使用纹理的原始尺寸
	if i.GetSize().X() <= 0.0 && i.GetSize().Y() <= 0.0 {
		ctx.GetRenderer().DrawUISprite(i.sprite, position, nil)
	} else {
		ctx.GetRenderer().DrawUISprite(i.sprite, position, &i.size)
	}

	// 渲染子元素，调用基类方法
	i.UIElement.Render(ctx)
}

// 获取纹理ID
func (i *UIImage) GetTextureId() string {
	return i.sprite.GetTextureId()
}
