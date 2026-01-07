package state

import (
	"sunny_land/src/engine/physics"

	"github.com/go-gl/mathgl/mgl32"
)

// 空闲状态
type IdleState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*IdleState)(nil)

// 创建空闲状态
func NewIdleState(playerCom IPlayerComponent) *IdleState {
	return &IdleState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (is *IdleState) Enter() {
	// 播放空闲动画
	is.PlayAnimation("idle")
}

// 更新状态
func (is *IdleState) Update(dt float64, ctx physics.IContext) IPlayerState {
	// 应用摩擦系数，水平方向
	// TODO: 摩擦力应该做到物理引擎中？
	physicsCom := is.playerCom.GetPhysicsComponent()
	frictionFactor := is.playerCom.GetFrictionFactor()
	physicsCom.Velocity[0] *= frictionFactor

	// 如果下方没有碰撞，则切换到下落状态
	if !physicsCom.HasCollidedBelow() {
		return NewFallState(is.playerCom)
	}

	return nil
}

// 输入
func (is *IdleState) HandleInput(ctx physics.IContext) IPlayerState {
	inputManager := ctx.GetInputManager()
	physicsCom := is.playerCom.GetPhysicsComponent()

	// 如果按"move_up"键，且与梯子重合，则切换到ClimbState
	if physicsCom.HasCollidedLadder() && inputManager.IsActionDown("move_up") {
		return NewClimbState(is.playerCom)
	}

	// 如果按下"move_down"且在梯子顶层，则切换到ClimbState
	if physicsCom.HasCollidedLadderTop() && inputManager.IsActionDown("move_down") {
		// 需要向下移动一点，确保下一帧能与梯子碰撞，否则会切换回FallState
		is.playerCom.GetTransformComponent().Translate(mgl32.Vec2{0, 2.0})
		return NewClimbState(is.playerCom)
	}

	// 如果按下了左右移动键，则切换到移动状态
	if inputManager.IsActionDown("move_left") || inputManager.IsActionDown("move_right") {
		// 切换到移动状态
		return NewWalkState(is.playerCom)
	}

	// 如果按下跳跃键，则切换到跳跃状态
	if inputManager.IsActionDown("jump") {
		return NewJumpState(is.playerCom)
	}

	return nil
}

// 退出状态
func (is *IdleState) Exit() {
}
