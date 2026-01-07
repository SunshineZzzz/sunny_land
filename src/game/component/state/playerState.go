package state

import (
	"log/slog"

	eComponent "sunny_land/src/engine/component"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
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
	// 获取动画组件
	GetAnimationComponent() *eComponent.AnimationComponent
	// 获取生命值组件
	GetHealthComponent() *eComponent.HealthComponent
	// 获取摩擦系数
	GetFrictionFactor() float32
	// 获取移动力
	GetMoveForce() float32
	// 获取跳跃速度
	GetJumpSpeed() float32
	// 获取最大水平速度
	GetMaxSpeed() float32
	// 获取攀爬速度
	GetClimbSpeed() float32
	// 获取硬直时间
	GetStunnedDuration() float64
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

// 播放指定名称的动画
func (p *playerState) PlayAnimation(animationName string) {
	if p.playerCom == nil {
		slog.Error("player state play animation failed, playerCom is nil")
		return
	}

	animationCom := p.playerCom.GetOwner().GetComponent(def.ComponentTypeAnimation).(*eComponent.AnimationComponent)
	if animationCom == nil {
		slog.Error("player state play animation failed, animationCom is nil")
		return
	}

	animationCom.PlayAnimation(animationName)
}
