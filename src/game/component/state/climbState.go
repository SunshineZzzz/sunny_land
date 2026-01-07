package state

import (
	"sunny_land/src/engine/physics"
)

// 攀爬状态
type ClimbState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*ClimbState)(nil)

// 创建攀爬状态
func NewClimbState(playerCom IPlayerComponent) *ClimbState {
	return &ClimbState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (cs *ClimbState) Enter() {
	// 播放攀爬动画
	cs.PlayAnimation("climb")
	physicsComponent := cs.playerCom.GetPhysicsComponent()
	if physicsComponent != nil {
		// 禁用重力
		physicsComponent.SetUseGravity(false)
	}
}

// 退出状态
func (cs *ClimbState) Exit() {
	physicsComponent := cs.playerCom.GetPhysicsComponent()
	if physicsComponent != nil {
		// 启用重力
		physicsComponent.SetUseGravity(true)
	}
}

// 处理输入
func (cs *ClimbState) HandleInput(ctx physics.IContext) IPlayerState {
	inputManager := ctx.GetInputManager()
	physicsCom := cs.playerCom.GetPhysicsComponent()
	animCom := cs.playerCom.GetAnimationComponent()

	// 攀爬状态下，按键则移动，不按键则静止
	isUp := inputManager.IsActionDown("move_up")
	isDown := inputManager.IsActionDown("move_down")
	isLeft := inputManager.IsActionDown("move_left")
	isRight := inputManager.IsActionDown("move_right")
	speed := cs.playerCom.GetClimbSpeed()

	if isUp {
		physicsCom.Velocity[1] = -speed
	} else if !isUp && isDown {
		physicsCom.Velocity[1] = speed
	} else if !isUp && !isDown {
		physicsCom.Velocity[1] = 0.0
	}

	if isLeft {
		physicsCom.Velocity[0] = -speed
	} else if !isLeft && isRight {
		physicsCom.Velocity[0] = speed
	} else if !isLeft && !isRight {
		physicsCom.Velocity[0] = 0.0
	}

	// 根据是否有按键决定动画播放情况
	if isUp || isDown || isLeft || isRight {
		// 有按键则恢复动画播放
		animCom.ResumeAnimation()
	} else {
		// 无按键则停止动画播放
		animCom.StopAnimation()
	}

	// 按跳跃键主动离开攀爬状态
	if inputManager.IsActionDown("jump") {
		return NewJumpState(cs.playerCom)
	}
	return nil
}

// 更新状态
func (cs *ClimbState) Update(float64, physics.IContext) IPlayerState {
	physicsCom := cs.playerCom.GetPhysicsComponent()
	// 如果着地，则切换到IdleState
	if physicsCom.HasCollidedBelow() {
		return NewIdleState(cs.playerCom)
	}
	// 如果离开梯子区域，则切换到FallState，能走到这里，说明非着地状态
	if !physicsCom.HasCollidedLadder() {
		return NewFallState(cs.playerCom)
	}
	return nil
}
