package render

import (
	"log/slog"

	"sunny_land/src/engine/utils/math"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// 相机
type Camera struct {
	// 视口大小(屏幕大小)
	viewportSize mgl32.Vec2
	// 相机左上角的世界坐标
	position mgl32.Vec2
	// 限制相机在世界中的移动范围，nil表示不限制
	limitBounds *emath.Rect
}

// 创建相机
func NewCamera(viewportSize, position mgl32.Vec2, limitBounds *emath.Rect) *Camera {
	slog.Debug("create camera", slog.Any("viewportSize", viewportSize), slog.Any("position", position), slog.Any("limitBounds", limitBounds))
	return &Camera{
		viewportSize: viewportSize,
		position:     position,
		limitBounds:  limitBounds,
	}
}

// 设置相机位置
func (c *Camera) SetPosition(position mgl32.Vec2) {
	c.position = position
	c.ClampPosition()
}

// 更新
func (c *Camera) Update(deltaTime float64) {
}

// 移动相机
func (c *Camera) Move(direction mgl32.Vec2) {
	c.position = c.position.Add(direction)
	c.ClampPosition()
}

// 设置相机限制范围
func (c *Camera) SetLimitBounds(limitBounds *math.Rect) {
	c.limitBounds = limitBounds
	c.ClampPosition()
}

// 获取相机位置
func (c *Camera) GetPosition() mgl32.Vec2 {
	return c.position
}

// 限制相机位置在限制范围内
func (c *Camera) ClampPosition() {
	if c.limitBounds == nil || c.limitBounds.Size.X() <= 0.0 || c.limitBounds.Size.Y() <= 0.0 {
		return
	}

	// 计算允许相机移动位置范围
	minCamPos := c.limitBounds.Position
	maxCamPos := c.limitBounds.Position.Add(c.limitBounds.Size).Sub(c.viewportSize)

	// 确保maxCamPos不小于minCamPos，视口可能比世界还大
	maxCamPos[0] = min(maxCamPos.X(), minCamPos.X())
	maxCamPos[1] = min(maxCamPos.Y(), minCamPos.Y())

	// 限制相机位置在范围内
	c.position = math.Mgl32Vec2Clamp(c.position, minCamPos, maxCamPos)
}

// 世界坐标转换为屏幕坐标(视口坐标)
func (c *Camera) WorldToScreen(worldPos mgl32.Vec2) mgl32.Vec2 {
	return worldPos.Sub(c.position)
}

// 世界坐标转换为屏幕坐标(视口坐标)，考虑视差
// scrollFactor，视差系数，用于计算视差效果，0.0表示没有视差(固定背景)，1.0表示背景跟着相机移动，0.0~1.0之间视差
// 移动得越快，看起来离玩家越近, 视差系数越接近1.0。移动得越慢，看起来离玩家越远，视差系数越接近0.0
func (c *Camera) WorldToScreenWithParallax(worldPos mgl32.Vec2, scrollFactor mgl32.Vec2) mgl32.Vec2 {
	return worldPos.Sub(math.Mgl32Vec2MulElem(c.position, scrollFactor))
}

// 屏幕坐标(视口坐标)转换为世界坐标
func (c *Camera) ScreenToWorld(screenPos mgl32.Vec2) mgl32.Vec2 {
	return screenPos.Add(c.position)
}

// 获取相机视口大小(屏幕大小)
func (c *Camera) GetViewportSize() mgl32.Vec2 {
	return c.viewportSize
}

// 获取相机限制范围
func (c *Camera) GetLimitBounds() *math.Rect {
	return c.limitBounds
}
