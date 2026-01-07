package state

import (
	"sunny_land/src/engine/physics"

	"github.com/go-gl/mathgl/mgl32"
)

// 移动状态
type WalkState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*WalkState)(nil)

// 创建移动状态
func NewWalkState(playerCom IPlayerComponent) *WalkState {
	return &WalkState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (ws *WalkState) Enter() {
	// 播放步行动画
	ws.PlayAnimation("walk")
}

// 离开状态
func (ws *WalkState) Exit() {
}

// 处理输入
func (ws *WalkState) HandleInput(ctx physics.IContext) IPlayerState {
	inputManager := ctx.GetInputManager()
	physicsCom := ws.playerCom.GetPhysicsComponent()
	spriteCom := ws.playerCom.GetSpriteComponent()

	// 如果按下了跳跃键，则切换到跳跃状态
	if inputManager.IsActionDown("jump") {
		return NewJumpState(ws.playerCom)
	}

	// 步行状态可以左右移动
	if inputManager.IsActionDown("move_left") {
		if physicsCom.Velocity.X() > 0.0 {
			// 如果当前速度是向右的，先减速到0.0
			physicsCom.Velocity[0] = 0.0
		}
		// 添加向左的水平力
		physicsCom.AddForce(mgl32.Vec2{-ws.playerCom.GetMoveForce(), 0.0})
		// 向左移动需要反转
		spriteCom.SetIsFliped(true)
	} else if inputManager.IsActionDown("move_right") {
		if physicsCom.Velocity.X() < 0.0 {
			// 如果当前速度是向左的，先减速到0.0
			physicsCom.Velocity[0] = 0.0
		}
		// 添加向右的水平力
		physicsCom.AddForce(mgl32.Vec2{ws.playerCom.GetMoveForce(), 0.0})
		// 向右移动不需要反转
		spriteCom.SetIsFliped(false)
	} else {
		// 左右都没有，切换到空闲
		return NewIdleState(ws.playerCom)
	}
	return nil
}

// 更新
func (ws *WalkState) Update(deltaTime float64, ctx physics.IContext) IPlayerState {
	// 限制最大速度
	physicsCom := ws.playerCom.GetPhysicsComponent()
	physicsCom.Velocity[0] = mgl32.Clamp(physicsCom.Velocity.X(), -ws.playerCom.GetMaxSpeed(), ws.playerCom.GetMaxSpeed())

	// 如果下方没有碰撞，切换到下落状态
	if !ws.playerCom.IsOnGround() {
		return NewFallState(ws.playerCom)
	}
	return nil
}
