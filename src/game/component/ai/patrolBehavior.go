package ai

import (
	"log/slog"

	"sunny_land/src/engine/component"
	"sunny_land/src/engine/utils/def"
)

/**
 * @brief AI 行为：在指定范围内左右巡逻。
 *
 * 遇到墙壁或到达巡逻边界时会转身。
 */
type PatrolBehavior struct {
	// 巡逻范围的左边界
	patrolMinX float32
	// 巡逻范围的右边界
	patrolMaxX float32
	// 移动速度(像素/秒)
	moveSpeed float32
	// 当前是否向右移动
	movingRight bool
}

// 确保PatrolBehavior实现了IAIBehavior接口
var _ component.IAIBehavior = (*PatrolBehavior)(nil)

// 创建巡逻行为
func NewPatrolBehavior(minX, maxX, speed float32) *PatrolBehavior {
	if minX >= maxX {
		slog.Error("PatrolBehavior: minX must be less than maxX")
		minX = maxX
	}
	if speed <= 0.0 {
		slog.Error("PatrolBehavior: speed must be greater than 0.0")
		speed = 50.0
	}
	return &PatrolBehavior{
		patrolMinX:  minX,
		patrolMaxX:  maxX,
		moveSpeed:   speed,
		movingRight: true,
	}
}

// 进入行为
func (pb *PatrolBehavior) Enter(aiComponent *component.AIComponent) {
	animComponent := aiComponent.GetOwner().GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
	if animComponent != nil {
		animComponent.PlayAnimation("walk")
	}
}

// 更新行为
func (pb *PatrolBehavior) Update(dt float64, aiComponent *component.AIComponent) {
	_ = dt

	// 获取必要组件
	physicComponent := aiComponent.Owner.GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent)
	transformComponent := aiComponent.Owner.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent)
	spriteComponent := aiComponent.Owner.GetComponent(def.ComponentTypeSprite).(*component.SpriteComponent)
	if physicComponent == nil || transformComponent == nil || spriteComponent == nil {
		slog.Error("patrol behavior owner must have physics, transform, sprite component")
		return
	}

	// 检查碰撞和边界
	currentX := transformComponent.GetPosition().X()
	// 撞右墙或到达设定目标则转向左
	if currentX >= pb.patrolMaxX || physicComponent.HasCollidedRight() {
		physicComponent.Velocity[0] = -pb.moveSpeed
		pb.movingRight = false
	}
	// 撞左墙或到达设定目标则转向右
	if currentX <= pb.patrolMinX || physicComponent.HasCollidedLeft() {
		physicComponent.Velocity[0] = pb.moveSpeed
		pb.movingRight = true
	}

	// 更新精灵方向
	spriteComponent.SetIsFliped(pb.movingRight)
}
