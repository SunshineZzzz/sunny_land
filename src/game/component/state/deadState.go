package state

import (
	eComponent "sunny_land/src/engine/component"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"

	"github.com/go-gl/mathgl/mgl32"
)

// 死亡状态
type DeadState struct {
	// 继承基础玩家状态
	playerState
}

// 确保实现了玩家状态接口
var _ IPlayerState = (*DeadState)(nil)

// 创建死亡状态
func NewDeadState(playerCom IPlayerComponent) *DeadState {
	return &DeadState{
		playerState: playerState{
			playerCom: playerCom,
		},
	}
}

// 进入状态
func (ds *DeadState) Enter() {
	// 播放死亡动画
	ds.PlayAnimation("hurt")

	// 应用击退力，只向上
	physicsCom := ds.playerCom.GetOwner().GetComponent(def.ComponentTypePhysics).(*eComponent.PhysicsComponent)
	physicsCom.Velocity = mgl32.Vec2{0.0, -200.0}

	// 禁用碰撞组件，自动掉出屏幕
	colliderCom := ds.playerCom.GetOwner().GetComponent(def.ComponentTypeCollider).(*eComponent.ColliderComponent)
	if colliderCom != nil {
		colliderCom.SetActive(false)
	}

	audioComponent := ds.playerCom.GetOwner().GetComponent(def.ComponentTypeAudio).(*eComponent.AudioComponent)
	if audioComponent != nil {
		// 播放死亡音效
		audioComponent.PlaySound("dead", false)
	}
}

// 退出状态
func (ds *DeadState) Exit() {
}

// 处理输入
func (ds *DeadState) HandleInput(physics.IContext) IPlayerState {
	// 死亡状态下不处理输入
	return nil
}

// 更新状态
func (ds *DeadState) Update(float64, physics.IContext) IPlayerState {
	// 死亡状态下不处理更新
	return nil
}
