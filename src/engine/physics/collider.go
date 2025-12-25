package physics

import "github.com/go-gl/mathgl/mgl32"

// 碰撞器类型
type ColliderType int

const (
	// 无碰撞器
	ColliderTypeNone ColliderType = iota
	// 轴对齐包围盒碰撞器
	ColliderTypeAABB
	// 圆形碰撞器
	ColliderTypeCircle
)

// 碰撞器抽象
type ICollider interface {
	// 设置轴对齐包围盒的尺寸
	SetAABBSize(mgl32.Vec2)
	// 获取轴对齐包围盒的尺寸
	GetAABBSize() mgl32.Vec2
	// 获取类型
	GetType() ColliderType
}

// 碰撞器基类
type Collider struct {
	// 轴对齐最小包围盒的尺寸
	aabbSize mgl32.Vec2
}

// 确保Collider实现了ICollider接口
var _ ICollider = (*Collider)(nil)

// 设置轴对齐包围盒的尺寸
func (c *Collider) SetAABBSize(size mgl32.Vec2) {
	c.aabbSize = size
}

// 获取轴对齐包围盒的尺寸
func (c *Collider) GetAABBSize() mgl32.Vec2 {
	return c.aabbSize
}

// 获取类型
func (c *Collider) GetType() ColliderType {
	return ColliderTypeNone
}

// 轴对齐(意味着盒子的每一条边都必须和屏幕的水平轴或垂直轴平行，不可旋转)包围盒碰撞器
type AABBCollider struct {
	// 继承碰撞器基类
	Collider
	// 包围盒尺寸
	size mgl32.Vec2
}

// 确保Collider实现了ICollider接口
var _ ICollider = (*AABBCollider)(nil)

// 创建轴对齐包围盒碰撞器
func NewAABBCollider(size mgl32.Vec2) *AABBCollider {
	return &AABBCollider{
		Collider: Collider{
			aabbSize: size,
		},
		size: size,
	}
}

// 获取类型
func (c *AABBCollider) GetType() ColliderType {
	return ColliderTypeAABB
}

// 获取轴对齐包围盒尺寸
func (c *AABBCollider) GetAABBSize() mgl32.Vec2 {
	return c.size
}

// 设置轴对齐包围盒尺寸
func (c *AABBCollider) SetAABBSize(size mgl32.Vec2) {
	c.size = size
	c.aabbSize = size
}

// 圆形碰撞器
type CircleCollider struct {
	// 继承碰撞器基类
	Collider
	// 半径
	radius float32
}

// 确保Collider实现了ICollider接口
var _ ICollider = (*CircleCollider)(nil)

// 创建圆形碰撞器
func NewCircleCollider(radius float32) *CircleCollider {
	return &CircleCollider{
		Collider: Collider{
			aabbSize: mgl32.Vec2{radius * 2.0, radius * 2.0},
		},
		radius: radius,
	}
}

// 获取类型
func (c *CircleCollider) GetType() ColliderType {
	return ColliderTypeCircle
}

// 获取半径
func (c *CircleCollider) GetRadius() float32 {
	return c.radius
}

// 设置半径
func (c *CircleCollider) SetRadius(radius float32) {
	c.radius = radius
	c.aabbSize = mgl32.Vec2{radius * 2.0, radius * 2.0}
}
