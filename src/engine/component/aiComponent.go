package component

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
)

// AI行为策略抽象
type IAIBehavior interface {
	// 进入行为
	Enter(*AIComponent)
	// 更新行为
	Update(float64, *AIComponent)
}

// AI组件，负责管理游戏对象的AI行为
type AIComponent struct {
	// 继承组件基类
	Component
	// 当前AI行为策略
	currentBehavior IAIBehavior
	// 缓存组件
	transformComponent *TransformComponent
	physicsComponent   *PhysicsComponent
	spriteComponent    *SpriteComponent
	animationComponent *AnimationComponent
	audioComponent     *AudioComponent
}

// 确保AIComponent实现了IComponent接口
var _ physics.IComponent = (*AIComponent)(nil)

// 创建AI组件
func NewAIComponent() *AIComponent {
	return &AIComponent{
		Component: Component{
			ComponentType: def.ComponentTypeAI,
		},
	}
}

// 初始化
func (ac *AIComponent) Init() {
	if ac.Owner == nil {
		slog.Error("owner is nil")
		return
	}

	ac.transformComponent = ac.Owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if ac.transformComponent == nil {
		slog.Warn("transform component is nil", slog.String("owner", ac.Owner.GetName()))
		return
	}

	ac.physicsComponent = ac.Owner.GetComponent(def.ComponentTypePhysics).(*PhysicsComponent)
	if ac.physicsComponent == nil {
		slog.Warn("physics component is nil", slog.String("owner", ac.Owner.GetName()))
		return
	}

	ac.spriteComponent = ac.Owner.GetComponent(def.ComponentTypeSprite).(*SpriteComponent)
	if ac.spriteComponent == nil {
		slog.Warn("sprite component is nil", slog.String("owner", ac.Owner.GetName()))
		return
	}

	ac.animationComponent = ac.Owner.GetComponent(def.ComponentTypeAnimation).(*AnimationComponent)
	if ac.animationComponent == nil {
		slog.Warn("animation component is nil", slog.String("owner", ac.Owner.GetName()))
		return
	}

	// 音频组件不一定存在
	if ac.Owner.GetComponent(def.ComponentTypeAudio) != nil {
		ac.audioComponent = ac.Owner.GetComponent(def.ComponentTypeAudio).(*AudioComponent)
	}
}

// 更新
func (ac *AIComponent) Update(dt float64, ctx physics.IContext) {
	_ = ctx

	// 将更新委托给当前的行为策略
	if ac.currentBehavior != nil {
		ac.currentBehavior.Update(dt, ac)
	} else {
		slog.Warn("current behavior is nil", slog.String("owner", ac.Owner.GetName()))
	}
}

// 设置当前行为策略
func (ac *AIComponent) SetBehavior(behavior IAIBehavior) {
	ac.currentBehavior = behavior
	if behavior != nil {
		behavior.Enter(ac)
	}
}

// 处理伤害逻辑，返回是否造成伤害
func (ac *AIComponent) TakeDamage(damage int) bool {
	success := false
	healthComponent := ac.GetOwner().GetComponent(def.ComponentTypeHealth).(*HealthComponent)
	if healthComponent != nil {
		success = healthComponent.TakeDamage(damage)
		// TODO: 处理伤害/死亡后的行为
	}
	return success
}

// 是否活着
func (ac *AIComponent) IsAlive() bool {
	healthComponent := ac.GetOwner().GetComponent(def.ComponentTypeHealth).(*HealthComponent)
	if healthComponent != nil {
		return healthComponent.IsAlive()
	}
	// 如果没有生命组件，默认返回存活状态
	return true
}
