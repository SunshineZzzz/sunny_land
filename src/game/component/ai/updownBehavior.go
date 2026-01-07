package ai

import (
	"log/slog"

	"sunny_land/src/engine/component"
	"sunny_land/src/engine/utils/def"
)

/**
 * @brief AI 行为：在指定范围内上下垂直移动。
 *
 * 到达边界或碰到障碍物时会反向。
 */
type UpDownBehavior struct {
	// 巡逻范围的上边界
	patrolMinY float32
	// 巡逻范围的下边界
	patrolMaxY float32
	// 移动速度(像素/秒)
	moveSpeed float32
	// 当前是否正在向下移动
	movingDown bool
}

// 确保UpDownBehavior实现了IAIBehavior接口
var _ component.IAIBehavior = (*UpDownBehavior)(nil)

// 创建上下行为
func NewUpDownBehavior(patrolMinY, patrolMaxY, moveSpeed float32) *UpDownBehavior {
	if patrolMinY >= patrolMaxY {
		slog.Error("up down behavior patrol min y must less than max y")
		patrolMinY = patrolMaxY
	}
	// 确保移动速度为正值
	if moveSpeed <= 0.0 {
		slog.Error("up down behavior move speed must greater than 0")
		moveSpeed = 50.0
	}
	return &UpDownBehavior{
		patrolMinY: patrolMinY,
		patrolMaxY: patrolMaxY,
		moveSpeed:  moveSpeed,
		movingDown: true,
	}
}

// 进入行为
func (ub *UpDownBehavior) Enter(aiComponent *component.AIComponent) {
	animComponent := aiComponent.GetOwner().GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
	if animComponent != nil {
		animComponent.PlayAnimation("fly")
	}
	physicsCom := aiComponent.GetOwner().GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent)
	if physicsCom != nil {
		physicsCom.SetUseGravity(false)
	}
}

// 更新行为
func (ub *UpDownBehavior) Update(dt float64, aiComponent *component.AIComponent) {
	_ = dt

	// 获取必要组件
	physicComponent := aiComponent.Owner.GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent)
	transformComponent := aiComponent.Owner.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent)
	spriteComponent := aiComponent.Owner.GetComponent(def.ComponentTypeSprite).(*component.SpriteComponent)
	if physicComponent == nil || transformComponent == nil || spriteComponent == nil {
		slog.Error("up down behavior owner must have physics, transform, sprite component")
		return
	}

	// 检查碰撞和边界
	currentY := transformComponent.GetPosition().Y()
	// 到达上边界或碰到上方障碍，向下移动
	if currentY <= ub.patrolMinY || physicComponent.HasCollidedAbove() {
		physicComponent.Velocity[1] = ub.moveSpeed
		ub.movingDown = true
	} else if currentY >= ub.patrolMaxY || physicComponent.HasCollidedBelow() {
		physicComponent.Velocity[1] = -ub.moveSpeed
		ub.movingDown = false
	}
	// 不需要反转精灵图
}
