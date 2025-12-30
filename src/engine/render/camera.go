package render

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/math"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// 相机移动距离阈值，低于该值时，相机位置直接设置到目标位置
	SNAP_THRESHOLD = float32(1.0)
)

// 相机
type Camera struct {
	// 视口大小(屏幕大小)
	viewportSize mgl32.Vec2
	// 相机左上角的世界坐标
	position mgl32.Vec2
	// 限制相机在世界中的移动范围，nil表示不限制
	limitBounds *emath.Rect
	// 相机跟随目标变化组件，空置表示不跟随
	targetTC physics.ITransformComponent
	// 相机平滑移动速度
	smoothSpeed float32
}

// 确保Camera实现了ICamera接口
var _ physics.ICamera = (*Camera)(nil)

// 创建相机
func NewCamera(viewportSize, position mgl32.Vec2, limitBounds *emath.Rect) *Camera {
	slog.Debug("create camera", slog.Any("viewportSize", viewportSize), slog.Any("position", position), slog.Any("limitBounds", limitBounds))
	return &Camera{
		viewportSize: viewportSize,
		position:     position,
		limitBounds:  limitBounds,
		smoothSpeed:  5.0,
	}
}

// 设置相机位置
func (c *Camera) SetPosition(position mgl32.Vec2) {
	c.position = position
	c.ClampPosition()
}

// 更新
func (c *Camera) Update(deltaTime float64) {
	if c.targetTC == nil {
		return
	}

	// 计算相机目标位置
	targetTCPos := c.targetTC.GetPosition()
	// 目标向量与视口中心向量差值
	desiredPos := targetTCPos.Sub(c.viewportSize.Mul(0.5))
	// 计算相机当前位置与想要去的位置差值
	distance := c.position.Sub(desiredPos).Len()

	if distance < SNAP_THRESHOLD {
		// 如果相机距离小于阈值，直接设置相机位置到目标位置
		c.position = desiredPos
	} else {
		// 如果相机距离大于阈值，平滑移动相机位置到目标位置
		// 这个算法应该准确的说是平滑减速
		// 第 1 帧：你走总路程的 50%。
		// 第 2 帧：你走剩余路程的 50%。
		// 第 3 帧：你走剩余之剩余路程的 50%。
		c.position = emath.Mgl32Vec2Mix(c.position, desiredPos, c.smoothSpeed*float32(deltaTime))
		// 取整一下，要不画面撕裂
		c.position = mgl32.Vec2{
			mgl32.Round(c.position.X(), 0),
			mgl32.Round(c.position.Y(), 0),
		}
	}

	c.ClampPosition()
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
	maxCamPos[0] = max(maxCamPos.X(), minCamPos.X())
	maxCamPos[1] = max(maxCamPos.Y(), minCamPos.Y())

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

// 设置相机跟随目标
func (c *Camera) SetTargetTC(targetTC physics.ITransformComponent) {
	c.targetTC = targetTC
}

// 获取相机跟随目标
func (c *Camera) GetTargetTC() physics.ITransformComponent {
	return c.targetTC
}
