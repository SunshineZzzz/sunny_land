package component

import (
	"log/slog"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"

	"github.com/go-gl/mathgl/mgl32"
)

// 变换组件，负责游戏对象的位置、旋转和缩放
type TransformComponent struct {
	// 继承基础组件
	Component
	// 位置
	position mgl32.Vec2
	// 缩放
	scale mgl32.Vec2
	// 旋转（角度）
	rotation float64
}

// 确保TransformComponent实现了IComponent接口
var _ object.IComponent = (*TransformComponent)(nil)

// 创建变换组件
func NewTransformComponent(position mgl32.Vec2, scale mgl32.Vec2, rotation float64) *TransformComponent {
	slog.Debug("create transform component", slog.Any("position", position), slog.Any("scale", scale), slog.Float64("rotation", rotation))
	return &TransformComponent{
		position: position,
		scale:    scale,
		rotation: rotation,
	}
}

// 设置缩放
func (tc *TransformComponent) SetScale(scale mgl32.Vec2) {
	tc.scale = scale
	if tc.owner != nil {
		spriteComp := tc.owner.GetComponent(&SpriteComponent{}).(*SpriteComponent)
		if spriteComp != nil {
			spriteComp.updateOffset()
		}
	}
}

// 获取位置
func (tc *TransformComponent) GetPosition() mgl32.Vec2 {
	return tc.position
}

// 获取旋转角度
func (tc *TransformComponent) GetRotation() float64 {
	return tc.rotation
}

// 获取缩放
func (tc *TransformComponent) GetScale() mgl32.Vec2 {
	return tc.scale
}

// 设置位置
func (tc *TransformComponent) SetPosition(position mgl32.Vec2) {
	tc.position = position
}

// 设置旋转角度
func (tc *TransformComponent) SetRotation(rotation float64) {
	tc.rotation = rotation
}

// 平移位置
func (tc *TransformComponent) Translate(delta mgl32.Vec2) {
	tc.position = tc.position.Add(delta)
}

// 更新
func (tc *TransformComponent) Update(float64, *econtext.Context) {
}
