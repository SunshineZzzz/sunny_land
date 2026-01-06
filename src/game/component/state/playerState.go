package state

import (
	eComponent "sunny_land/src/engine/component"
	"sunny_land/src/engine/physics"
)

// 玩家组件接口，定义玩家组件需要实现的方法
type IPlayerComponent interface {
	// 继承基础组件接口
	physics.IComponent
	// 获取变换组件
	GetTransformComponent() *eComponent.TransformComponent
	// 获取精灵图组件
	GetSpriteComponent() *eComponent.SpriteComponent
	// 获取物理组件
	GetPhysicsComponent() *eComponent.PhysicsComponent
	// 获取摩擦系数
	GetFrictionFactor() float32
	// 获取移动力
	GetMoveForce() float32
	// 获取跳跃速度
	GetJumpSpeed() float32
	// 获取最大水平速度
	GetMaxSpeed() float32
}

// 玩家状态机抽象
type IPlayerState interface {
	// 进入状态
	Enter()
	// 更新状态
	Update(float64, physics.IContext) IPlayerState
	// 输入
	HandleInput(physics.IContext) IPlayerState
	// 退出状态
	Exit()
}

// 基础玩家状态
type playerState struct {
	// 玩家组件
	playerCom IPlayerComponent
}
