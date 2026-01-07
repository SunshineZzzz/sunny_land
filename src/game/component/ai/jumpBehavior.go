package ai

import (
	"log/slog"

	"sunny_land/src/engine/component"
	"sunny_land/src/engine/utils/def"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief AI 行为：在指定范围内周期性地跳跃。
 *
 * 在地面时等待，然后向当前方向跳跃。
 * 撞墙或到达边界时改变下次跳跃方向。
 */
type JumpBehavior struct {
	// 巡逻范围的左边界
	patrolMinX float32
	// 巡逻范围的右边界
	patrolMaxX float32
	// 跳跃速度
	jumpVel mgl32.Vec2
	// 跳跃间隔时间(秒)
	jumpInterval float64
	// 距离下次跳跃的计时器
	jumpTimer float64
	// 当前是否向右跳跃
	jumpingRight bool
}

// 确保JumpBehavior实现了IAIBehavior接口
var _ component.IAIBehavior = (*JumpBehavior)(nil)

// 创建跳跃行为
func NewJumpBehavior(patrolMinX, patrolMaxX float32, jumpVel mgl32.Vec2, jumpInterval float64) *JumpBehavior {
	// 检查参数是否有效
	if patrolMinX >= patrolMaxX {
		slog.Error("jump behavior patrol min x must be less than patrol max x")
		patrolMinX = patrolMaxX
	}
	if jumpInterval <= 0.0 {
		slog.Error("jump behavior jump interval must be greater than 0")
		jumpInterval = 2.0
	}
	if jumpVel.Y() > 0.0 {
		slog.Error("jump behavior jump vel must be less than 0")
		jumpVel[1] = -jumpVel[1]
	}
	return &JumpBehavior{
		patrolMinX:   patrolMinX,
		patrolMaxX:   patrolMaxX,
		jumpVel:      jumpVel,
		jumpInterval: jumpInterval,
		jumpTimer:    0.0,
		jumpingRight: true,
	}
}

// 进入行为
func (jb *JumpBehavior) Enter(*component.AIComponent) {
}

// 更新行为
func (jb *JumpBehavior) Update(dt float64, aiComponent *component.AIComponent) {
	// 获取必要组件
	physicComponent := aiComponent.Owner.GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent)
	transformComponent := aiComponent.Owner.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent)
	spriteComponent := aiComponent.Owner.GetComponent(def.ComponentTypeSprite).(*component.SpriteComponent)
	animationComponent := aiComponent.Owner.GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
	if physicComponent == nil || transformComponent == nil || spriteComponent == nil || animationComponent == nil {
		slog.Error("jump behavior owner must have physics, transform, sprite and animation component")
		return
	}

	// 着地标志
	isOnGround := physicComponent.HasCollidedBelow()
	if isOnGround {
		// 在地面上
		// 增加跳跃定时器
		jb.jumpTimer += dt
		// 停止水平移动
		physicComponent.Velocity[0] = 0.0
		// 检查是否到了跳跃时间
		if jb.jumpTimer >= jb.jumpInterval {
			// 重置跳跃定时器
			jb.jumpTimer = 0.0
			// 检查是否需要切换跳跃方向
			currentX := transformComponent.GetPosition().X()
			// 如果右边超限或者碰撞，向左跳
			if jb.jumpingRight && (physicComponent.HasCollidedRight() || currentX >= jb.patrolMaxX) {
				// 切换跳跃方向
				jb.jumpingRight = false
			} else if !jb.jumpingRight && (physicComponent.HasCollidedLeft() || currentX <= jb.patrolMinX) {
				// 切换跳跃方向
				jb.jumpingRight = true
			}
			// 确定水平跳跃方向
			if jb.jumpingRight {
				// 向右跳
				physicComponent.Velocity[0] = jb.jumpVel.X()
			} else {
				// 向左跳
				physicComponent.Velocity[0] = -jb.jumpVel.X()
			}
			// 设置速度
			physicComponent.Velocity[1] = jb.jumpVel.Y()
			// 播放跳跃动画
			animationComponent.PlayAnimation("jump")
			// 更新精灵图反转
			spriteComponent.SetIsFliped(jb.jumpingRight)
		} else {
			// 还在地面等待
			animationComponent.PlayAnimation("idle")
		}
	} else {
		// 在空中， 根据垂直速度判断是上升(jump)还是下落(fall)
		if physicComponent.Velocity[1] < 0.0 {
			// 上升中，播放跳跃动画
			animationComponent.PlayAnimation("jump")
		} else {
			// 下降中，播放下落动画
			animationComponent.PlayAnimation("fall")
		}
	}
}
