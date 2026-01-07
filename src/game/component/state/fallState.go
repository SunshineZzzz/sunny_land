package state

import (
	"sunny_land/src/engine/physics"

	"github.com/go-gl/mathgl/mgl32"
)

// 下落状态
type FallState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*FallState)(nil)

// 创建下落状态
func NewFallState(playerCom IPlayerComponent) *FallState {
	return &FallState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (fs *FallState) Enter() {
	// 播放下落动画
	fs.PlayAnimation("fall")
}

// 离开状态
func (fs *FallState) Exit() {
}

// 更新状态
func (fs *FallState) Update(deltaTime float64, ctx physics.IContext) IPlayerState {
	// 限制最大速度
	physicsCom := fs.playerCom.GetPhysicsComponent()
	physicsCom.Velocity[0] = mgl32.Clamp(physicsCom.Velocity.X(), -fs.playerCom.GetMaxSpeed(), fs.playerCom.GetMaxSpeed())

	// 如果下方有碰撞，根据水平速度决定切换到移动状态还是空闲状态
	if physicsCom.HasCollidedBelow() {
		if physicsCom.Velocity.X() != 0.0 {
			return NewWalkState(fs.playerCom)
		}
		// 切换为空闲状态
		return NewIdleState(fs.playerCom)
	}
	return nil
}

// 处理输入
func (fs *FallState) HandleInput(ctx physics.IContext) IPlayerState {
	inputManager := ctx.GetInputManager()
	physicsCom := fs.playerCom.GetPhysicsComponent()
	spriteCom := fs.playerCom.GetSpriteComponent()

	// 如果按下上下键，且与梯子重合，则切换到ClimbState
	if (inputManager.IsActionDown("move_up") || inputManager.IsActionDown("move_down")) && physicsCom.HasCollidedLadder() {
		return NewClimbState(fs.playerCom)
	}

	// 跳跃状态下可以左右移动
	if inputManager.IsActionDown("move_left") {
		// 如果当前是向右移动，速度清空
		if physicsCom.Velocity.X() > 0.0 {
			physicsCom.Velocity = mgl32.Vec2{0.0, physicsCom.Velocity.Y()}
		}
		physicsCom.AddForce(mgl32.Vec2{-fs.playerCom.GetMoveForce(), 0.0})
		// 切换为向左移动
		spriteCom.SetIsFliped(true)
	} else if inputManager.IsActionDown("move_right") {
		// 如果当前是向左移动，速度清空
		if physicsCom.Velocity.X() < 0.0 {
			physicsCom.Velocity = mgl32.Vec2{0.0, physicsCom.Velocity.Y()}
		}
		physicsCom.AddForce(mgl32.Vec2{fs.playerCom.GetMoveForce(), 0.0})
		// 切换为向右移动
		spriteCom.SetIsFliped(false)
	}
	return nil
}
