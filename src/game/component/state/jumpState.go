package state

import (
	"sunny_land/src/engine/physics"

	"github.com/go-gl/mathgl/mgl32"
)

// 跳跃状态
type JumpState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*JumpState)(nil)

// 创建跳跃状态
func NewJumpState(playerCom IPlayerComponent) *JumpState {
	return &JumpState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (js *JumpState) Enter() {
	// 播放跳跃动画
	js.PlayAnimation("jump")
	physicsCom := js.playerCom.GetPhysicsComponent()
	// 向上跳跃
	physicsCom.Velocity[1] = -js.playerCom.GetJumpSpeed()
}

// 更新状态
func (js *JumpState) Update(deltaTime float64, ctx physics.IContext) IPlayerState {
	// 限制最大速度
	physicsCom := js.playerCom.GetPhysicsComponent()
	physicsCom.Velocity[0] = mgl32.Clamp(physicsCom.Velocity.X(), -js.playerCom.GetMaxSpeed(), js.playerCom.GetMaxSpeed())

	// 如果向下移动，切换到下落状态
	if physicsCom.Velocity.Y() > 0.0 {
		return NewFallState(js.playerCom)
	}
	return nil
}

// 输入
func (js *JumpState) HandleInput(ctx physics.IContext) IPlayerState {
	inputManager := ctx.GetInputManager()
	physicsCom := js.playerCom.GetPhysicsComponent()
	spriteCom := js.playerCom.GetSpriteComponent()

	// 如果按下上下键，且与梯子重合，则切换到ClimbState
	if physicsCom.HasCollidedLadder() &&
		(inputManager.IsActionDown("move_up") || inputManager.IsActionDown("move_down")) {
		return NewClimbState(js.playerCom)
	}

	// 跳跃状态下可以左右移动
	if inputManager.IsActionDown("move_left") {
		// 如果当前是向右移动，速度清空
		if physicsCom.Velocity.X() > 0.0 {
			physicsCom.Velocity = mgl32.Vec2{0.0, physicsCom.Velocity.Y()}
		}
		physicsCom.AddForce(mgl32.Vec2{-js.playerCom.GetMoveForce(), 0.0})
		// 切换为向左移动
		spriteCom.SetIsFliped(true)
	} else if inputManager.IsActionDown("move_right") {
		// 如果当前是向左移动，速度清空
		if physicsCom.Velocity.X() < 0.0 {
			physicsCom.Velocity = mgl32.Vec2{0.0, physicsCom.Velocity.Y()}
		}
		physicsCom.AddForce(mgl32.Vec2{js.playerCom.GetMoveForce(), 0.0})
		// 切换为向右移动
		spriteCom.SetIsFliped(false)
	}
	return nil
}

// 退出状态
func (js *JumpState) Exit() {
}
