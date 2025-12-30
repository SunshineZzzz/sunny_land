package component

import (
	"log/slog"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/physics"

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
}

// 确保SpriteComponent实现了IComponent接口
var _ object.IComponent = (*PhysicsComponent)(nil)

// 确保SpriteComponent实现了IPhysicsComponent接口
var _ physics.IPhysicsComponent = (*PhysicsComponent)(nil)

// 创建物理组件
func NewPhysicsComponent(physicsEngine *physics.PhysicsEngine, mass float32, useGravity bool) *PhysicsComponent {
	slog.Debug("create physics component", slog.Float64("mass", float64(mass)), slog.Bool("useGravity", useGravity))
	return &PhysicsComponent{
		physicsEngine:      physicsEngine,
		transformComponent: nil,
		mass:               mass,
		useGravity:         useGravity,
		isEnable:           true,
	}
}

// 初始化
func (pc *PhysicsComponent) Init() {
	if pc.owner == nil {
		slog.Error("physics component owner is nil")
		return
	}
	if pc.physicsEngine == nil {
		slog.Error("physics component physics engine is nil")
		return
	}
	// 从物体中获取变换组件
	pc.transformComponent = pc.owner.GetComponent(&TransformComponent{}).(*TransformComponent)
	if pc.transformComponent == nil {
		slog.Warn("physics component transform component is nil")
	}
	// 注册到物理引擎
	pc.physicsEngine.RegisterPhysicsComponent(pc)
	slog.Debug("physics component init", slog.String("gameObject.Name", pc.owner.GetName()))
}

// 清理
func (pc *PhysicsComponent) Clean() {
	pc.physicsEngine.UnregisterComponent(pc)
	slog.Debug("physics component clean", slog.String("gameObject.Name", pc.owner.GetName()))
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

// 获取碰撞组件
func (pc *PhysicsComponent) GetColliderComponent() physics.IColliderComponent {
	if pc.owner == nil {
		slog.Error("physics component owner is nil")
		return nil
	}
	return pc.owner.GetComponent(&ColliderComponent{}).(*ColliderComponent)
}

// 获取游戏对象
func (pc *PhysicsComponent) GetGameObject() any {
	return pc.GetOwner()
}

// 获取游戏对象标签
func (pc *PhysicsComponent) GetGameObjectTag() string {
	return pc.owner.GetTag()
}
