package state

import (
	"fmt"
	eComponent "sunny_land/src/engine/component"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"

	"github.com/go-gl/mathgl/mgl32"
)

// 受伤(硬值)状态
type HurtState struct {
	// 继承基础玩家状态
	playerState
	// 硬直时间(单位：秒)
	stunnedDuration float64
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*HurtState)(nil)

// 创建受伤状态
func NewHurtState(playerCom IPlayerComponent) *HurtState {
	return &HurtState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (hs *HurtState) Enter() {
	fmt.Println("enter hurt state")

	// 播放受伤动画
	hs.PlayAnimation("hurt")
	// 造成击退效果
	physicsCom := hs.playerCom.GetOwner().GetComponent(def.ComponentTypePhysics).(*eComponent.PhysicsComponent)
	spriteCom := hs.playerCom.GetOwner().GetComponent(def.ComponentTypeSprite).(*eComponent.SpriteComponent)
	// 默认左上方击退效果
	knockbackVelocity := mgl32.Vec2{-100.0, -150.0}
	// 如果玩家是向左移动的，击退方向改为右方
	if spriteCom.GetIsFliped() {
		knockbackVelocity[0] = -knockbackVelocity[0]
	}
	// 设置击退速度
	physicsCom.Velocity = knockbackVelocity
}

// 离开状态
func (hs *HurtState) Exit() {
}

// 输入
func (hs *HurtState) HandleInput(physics.IContext) IPlayerState {
	// 硬直时间不能进行任何操作
	return nil
}

// 更新状态
func (hs *HurtState) Update(dt float64, ctx physics.IContext) IPlayerState {
	_ = ctx

	// 更新硬直时间
	hs.stunnedDuration += dt
	// 两种情况离开受伤(硬值)状态
	// 1. 落地
	physicsCom := hs.playerCom.GetOwner().GetComponent(def.ComponentTypePhysics).(*eComponent.PhysicsComponent)
	if physicsCom.HasCollidedBelow() {
		if mgl32.Abs(physicsCom.Velocity.X()) < 1.0 {
			return NewIdleState(hs.playerCom)
		}
		return NewWalkState(hs.playerCom)
	}
	// 2. 硬直时间结束
	if hs.stunnedDuration > hs.playerCom.GetStunnedDuration() {
		// 重置硬直计时器
		hs.stunnedDuration = 0.0
		// 切换到下落状态
		return NewFallState(hs.playerCom)
	}

	// 否则，继续保持 hurt 状态
	return nil
}
