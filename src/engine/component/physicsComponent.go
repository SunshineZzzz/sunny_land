package component

import (
	"log/slog"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"

	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsComponent struct {
	// 继承基础组件
	Component
	// 物体速度
	Velocity mgl32.Vec2
	// 物理引擎
	physicsEngine *physics.PhysicsEngine
	// 缓存对象的变换组件
	transformComponent *TransformComponent
	// 当前帧受到的力，单位：Newton，1牛，如果一个力作用在1kg物体上，物体的加速度恰好是1像素每二次方秒
	force mgl32.Vec2
	// 质量，单位：千克
	mass float32
	// 是否受重力影响
	useGravity bool
	// 是否启用
	isEnable bool
	// 碰撞标志位
	// 是否与底部碰撞
	collidedBelow bool
	// 是否与顶部碰撞
	collidedAbove bool
	// 是否与左侧碰撞
	collidedLeft bool
	// 是否与右侧碰撞
	collidedRight bool
}

// 确保SpriteComponent实现了IComponent接口
var _ physics.IComponent = (*PhysicsComponent)(nil)

// 确保SpriteComponent实现了IPhysicsComponent接口
var _ physics.IPhysicsComponent = (*PhysicsComponent)(nil)

// 创建物理组件
func NewPhysicsComponent(physicsEngine *physics.PhysicsEngine, mass float32, useGravity bool) *PhysicsComponent {
	slog.Debug("create physics component", slog.Float64("mass", float64(mass)), slog.Bool("useGravity", useGravity))
	return &PhysicsComponent{
		Component: Component{
			ComponentType: def.ComponentTypePhysics,
		},
		physicsEngine:      physicsEngine,
		transformComponent: nil,
		mass:               mass,
		useGravity:         useGravity,
		isEnable:           true,
	}
}

// 初始化
func (pc *PhysicsComponent) Init() {
	if pc.Owner == nil {
		slog.Error("physics component owner is nil")
		return
	}
	if pc.physicsEngine == nil {
		slog.Error("physics component physics engine is nil")
		return
	}
	// 从物体中获取变换组件
	pc.transformComponent = pc.Owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if pc.transformComponent == nil {
		slog.Warn("physics component transform component is nil")
	}
	// 注册到物理引擎
	pc.physicsEngine.RegisterPhysicsComponent(pc)
	slog.Debug("physics component init", slog.String("gameObject.Name", pc.Owner.GetName()))
}

// 清理
func (pc *PhysicsComponent) Clean() {
	pc.physicsEngine.UnregisterComponent(pc)
	slog.Debug("physics component clean", slog.String("gameObject.Name", pc.Owner.GetName()))
}

// 组件是否启用
func (pc *PhysicsComponent) IsEnabled() bool {
	return pc.isEnable
}

// 组件是否受重力影响
func (pc *PhysicsComponent) IsUseGravity() bool {
	return pc.useGravity
}

// 设置是否受重力影响
func (pc *PhysicsComponent) SetUseGravity(useGravity bool) {
	pc.useGravity = useGravity
}

// 获取质量
func (pc *PhysicsComponent) GetMass() float32 {
	return pc.mass
}

// 添加力
func (pc *PhysicsComponent) AddForce(force mgl32.Vec2) {
	if pc.isEnable {
		pc.force = pc.force.Add(force)
	}
}

// 清除力
func (pc *PhysicsComponent) ClearForce() {
	pc.force = mgl32.Vec2{0.0, 0.0}
}

// 获取力
func (pc *PhysicsComponent) GetForce() mgl32.Vec2 {
	return pc.force
}

// 获取变换组件
func (pc *PhysicsComponent) GetTransformComponent() physics.ITransformComponent {
	return pc.transformComponent
}

// 获取速度
func (pc *PhysicsComponent) GetVelocity() mgl32.Vec2 {
	return pc.Velocity
}

// 设置速度
func (pc *PhysicsComponent) SetVelocity(velocity mgl32.Vec2) {
	pc.Velocity = velocity
}

// 重置所有碰撞标志
func (pc *PhysicsComponent) ResetCollisionFlags() {
	pc.collidedBelow = false
	pc.collidedAbove = false
	pc.collidedLeft = false
	pc.collidedRight = false
}

// 设置下方碰撞标志位
func (pc *PhysicsComponent) SetCollidedBelow(collided bool) {
	pc.collidedBelow = collided
}

// 设置上方碰撞标志位
func (pc *PhysicsComponent) SetCollidedAbove(collided bool) {
	pc.collidedAbove = collided
}

// 设置左侧碰撞标志位
func (pc *PhysicsComponent) SetCollidedLeft(collided bool) {
	pc.collidedLeft = collided
}

// 设置右侧碰撞标志位
func (pc *PhysicsComponent) SetCollidedRight(collided bool) {
	pc.collidedRight = collided
}

// 检查是否与底部碰撞
func (pc *PhysicsComponent) HasCollidedBelow() bool {
	return pc.collidedBelow
}

// 检查是否与顶部碰撞
func (pc *PhysicsComponent) HasCollidedAbove() bool {
	return pc.collidedAbove
}

// 检查是否与左侧碰撞
func (pc *PhysicsComponent) HasCollidedLeft() bool {
	return pc.collidedLeft
}

// 检查是否与右侧碰撞
func (pc *PhysicsComponent) HasCollidedRight() bool {
	return pc.collidedRight
}
