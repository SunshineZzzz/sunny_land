package def

// 组件类型
type ComponentType uint32

// 组件类型常量
const (
	// 变换组件
	ComponentTypeTransform ComponentType = iota
	// 精灵图组件
	ComponentTypeSprite
	// 瓦片图层组件
	ComponentTypeTileLayer
	// 物理组件
	ComponentTypePhysics
	// 碰撞器组件
	ComponentTypeCollider
	// 视差组件
	ComponentTypeParallax
	// 动画组件
	ComponentTypeAnimation
	// 健康组件
	ComponentTypeHealth
	// AI组件
	ComponentTypeAI

	// 玩家组件
	ComponentTypePlayer
)
