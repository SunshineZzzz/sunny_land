package component

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils"
	"sunny_land/src/engine/utils/def"
	"sunny_land/src/engine/utils/math"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// 碰撞器组件
type ColliderComponent struct {
	// 继承组件基类
	Component
	// 碰撞器
	collider physics.ICollider
	// 缓存变换组件
	transformComponent *TransformComponent
	// 偏移量，碰撞器(包围盒的)左上角相对于变换原点的偏移量
	offset mgl32.Vec2
	// 对齐方式
	align utils.Alignment
	// 是否为触发器，仅检测碰撞，不产生物理响应
	isTrigger bool
	// 是否激活
	isActive bool
}

// 确保ColliderComponent实现了IComponent接口
var _ physics.IComponent = (*ColliderComponent)(nil)

// 确保ColliderComponent实现了IColliderComponent接口
var _ physics.IColliderComponent = (*ColliderComponent)(nil)

// 创建碰撞器组件
func NewColliderComponent(collider physics.ICollider, align utils.Alignment, offset mgl32.Vec2, isTrigger, isActive bool) *ColliderComponent {
	slog.Debug("NewColliderComponent", slog.Any("collider", collider), slog.Any("align", align), slog.Any("offset", offset),
		slog.Bool("isTrigger", isTrigger), slog.Bool("isActive", isActive))
	return &ColliderComponent{
		Component: Component{
			componentType: def.ComponentTypeCollider,
		},
		collider:  collider,
		align:     align,
		offset:    offset,
		isTrigger: isTrigger,
		isActive:  isActive,
	}
}

// 初始化
func (c *ColliderComponent) Init() {
	if c.owner == nil {
		slog.Error("ColliderComponent Init: owner is nil")
		return
	}
	c.transformComponent = c.owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if c.transformComponent == nil {
		slog.Error("ColliderComponent Init: transform component is nil")
		return
	}
	// 在获取变换组件后更新偏移量
	c.updateOffset()
}

// 设置对齐方式
func (c *ColliderComponent) SetAlign(anchor utils.Alignment) {
	c.align = anchor
	if c.transformComponent != nil && c.collider != nil {
		c.updateOffset()
	}
}

// 更新偏移量
func (c *ColliderComponent) updateOffset() {
	if c.transformComponent == nil || c.collider == nil {
		return
	}

	// 获取碰撞器的包围盒尺寸
	colliderSize := c.collider.GetAABBSize()

	// 如果尺寸无效，偏移量为0
	if colliderSize.X() <= 0.0 || colliderSize.Y() <= 0.0 {
		c.offset = mgl32.Vec2{0.0, 0.0}
		return
	}
	// 获取变换组件缩放
	scale := c.transformComponent.GetScale()

	// 根据对齐方式更新偏移量
	switch c.align {
	case utils.AlignTopLeft:
		c.offset = mgl32.Vec2{0.0, 0.0}
	case utils.AlignTopCenter:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-colliderSize.X() / 2.0, 0.0}, scale)
	case utils.AlignTopRight:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-colliderSize.X(), 0.0}, scale)
	case utils.AlignCenterLeft:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{0.0, -colliderSize.Y() / 2.0}, scale)
	case utils.AlignCenter:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-colliderSize.X() / 2.0, -colliderSize.Y() / 2.0}, scale)
	case utils.AlignCenterRight:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{colliderSize.X(), -colliderSize.Y() / 2.0}, scale)
	case utils.AlignBottomLeft:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{0.0, -colliderSize.Y()}, scale)
	case utils.AlignBottomCenter:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-colliderSize.X() / 2.0, -colliderSize.Y()}, scale)
	case utils.AlignBottomRight:
		c.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-colliderSize.X(), -colliderSize.Y()}, scale)
	}
}

// 获取世界坐标系下的轴对齐包围盒
func (c *ColliderComponent) GetWorldAABB() math.Rect {
	if c.transformComponent == nil || c.collider == nil {
		return math.Rect{Position: mgl32.Vec2{0.0, 0.0}, Size: mgl32.Vec2{0.0, 0.0}}
	}
	// 计算包围盒的左上角坐标(position)
	topLeftPos := c.transformComponent.GetPosition().Add(c.offset)
	// 获取碰撞器的AABB尺寸
	aabbSize := c.collider.GetAABBSize()
	// 返回世界空间的AABB
	return math.Rect{Position: topLeftPos, Size: aabbSize}
}

// 获取缓存的变换组件
func (c *ColliderComponent) GetTransformComponent() physics.ITransformComponent {
	return c.transformComponent
}

// 获取碰撞器
func (c *ColliderComponent) GetCollider() physics.ICollider {
	return c.collider
}

// 获取当前偏移量
func (c *ColliderComponent) GetOffset() mgl32.Vec2 {
	return c.offset
}

// 获取对齐方式
func (c *ColliderComponent) GetAlign() utils.Alignment {
	return c.align
}

// 获取是否为触发器
func (c *ColliderComponent) IsTrigger() bool {
	return c.isTrigger
}

// 获取是否激活
func (c *ColliderComponent) IsActive() bool {
	return c.isActive
}

// 设置偏移量
func (c *ColliderComponent) SetOffset(offset mgl32.Vec2) {
	c.offset = offset
}

// 设置是否为触发器
func (c *ColliderComponent) SetTrigger(isTrigger bool) {
	c.isTrigger = isTrigger
}

// 设置是否激活
func (c *ColliderComponent) SetActive(isActive bool) {
	c.isActive = isActive
}
