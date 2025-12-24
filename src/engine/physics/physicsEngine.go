package physics

import (
	"log/slog"

	"sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// 变换组件抽象
type ITransformComponent interface {
	// 平移
	Translate(mgl32.Vec2)
}

// 物理组件抽象
type IPhysicsComponent interface {
	// 组件是否启用
	IsEnabled() bool
	// 组件是否受重力影响
	IsUseGravity() bool
	// 获取质量
	GetMass() float32
	// 添加力
	AddForce(mgl32.Vec2)
	// 获取力
	GetForce() mgl32.Vec2
	// 清除力
	ClearForce()
	// 获取变换组件
	GetTransformComponent() ITransformComponent
	// 获取速度
	GetVelocity() mgl32.Vec2
	// 设置速度
	SetVelocity(mgl32.Vec2)
}

// 物理引擎，负责管理和模拟物理行为，碰撞检测
type PhysicsEngine struct {
	// 注册的物理组件容器
	physicsComponents []IPhysicsComponent
	// 默认重力值{0.0, 980.0}，单位：像素每二次方秒，现实中是，9.8米/s^2，游戏中是，100像素 * 9.8米/s^2 = 980.0像素/s^2
	gravity mgl32.Vec2
	// 最大速度值{-500.0, -500.0}/{500.0, 500.0}，单位：像素/秒
	maxSpeed float32
}

// 创建物理引擎
func NewPhysicsEngine() *PhysicsEngine {
	slog.Debug("new physics engine")
	return &PhysicsEngine{
		physicsComponents: make([]IPhysicsComponent, 0),
		gravity:           mgl32.Vec2{0.0, 980.0},
		maxSpeed:          500.0,
	}
}

// 注册物理组件
func (pe *PhysicsEngine) RegisterPhysicsComponent(component IPhysicsComponent) {
	slog.Debug("register physics component")
	pe.physicsComponents = append(pe.physicsComponents, component)
}

// 移除注册物理组件
func (pe *PhysicsEngine) UnregisterComponent(component IPhysicsComponent) {
	slog.Debug("remove physics component")
	for i, comp := range pe.physicsComponents {
		if comp == component {
			pe.physicsComponents = append(pe.physicsComponents[:i], pe.physicsComponents[i+1:]...)
			return
		}
	}
}

// 更新
func (pe *PhysicsEngine) Update(deltaTime float64) {
	for _, pc := range pe.physicsComponents {
		if pc == nil || !pc.IsEnabled() {
			continue
		}

		// 是否使用重力，如果组件接受重力影响，F = m * a
		if pc.IsUseGravity() {
			pc.AddForce(pe.gravity.Mul(pc.GetMass()))
		}
		// 还可以添加其他影响，比如风力，摩擦力，目前不考虑

		// 更新速度，v += a * dt，其中 a = F / m
		pc.SetVelocity(
			pc.GetVelocity().Add(
				pc.GetForce().Mul(1.0 / pc.GetMass()).Mul(float32(deltaTime)),
			),
		)
		// 清除当前帧的力
		pc.ClearForce()

		// 更新位置
		tc := pc.GetTransformComponent()
		tc.Translate(pc.GetVelocity().Mul(float32(deltaTime)))

		// 限制最大速度
		pc.SetVelocity(
			math.Mgl32Vec2Clamp(
				pc.GetVelocity(),
				mgl32.Vec2{-pe.maxSpeed, -pe.maxSpeed},
				mgl32.Vec2{pe.maxSpeed, pe.maxSpeed},
			),
		)
	}
}
