package component

import (
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
)

// 管理游戏对象的生命值，处理伤害、治疗，并提供无敌帧功能
type HealthComponent struct {
	// 继承基础组件
	Component
	// 最大生命值
	maxHealth int
	// 当前生命值
	curHealth int
	// 是否处于无敌状态
	isInvincible bool
	// 受伤后无敌总时长(秒)
	invincibilityDuration float64
	// 无敌时间计时器
	invincibilityTimer float64
}

// 确保HealthComponent实现了IComponent接口
var _ physics.IComponent = (*HealthComponent)(nil)

// 确保HealthComponent实现了IHealthComponent接口
var _ physics.IHealthComponent = (*HealthComponent)(nil)

// 创建健康组件
func NewHealthComponent(maxHealth int, invincibilityDuration float64) *HealthComponent {
	hc := &HealthComponent{
		Component: Component{
			ComponentType: def.ComponentTypeHealth,
		},
		maxHealth:             max(1, maxHealth),
		invincibilityDuration: max(0.0, invincibilityDuration),
	}
	hc.curHealth = hc.maxHealth

	return hc
}

// 更新
func (hc *HealthComponent) Update(dt float64, ctx physics.IContext) {
	_ = ctx

	if hc.isInvincible {
		hc.invincibilityTimer -= dt
		if hc.invincibilityTimer <= 0.0 {
			hc.isInvincible = false
			hc.invincibilityTimer = 0.0
		}
	}
}

// 受到伤害
func (hc *HealthComponent) TakeDamage(damage int) bool {
	if damage <= 0 || !hc.IsAlive() {
		return false
	}

	if hc.isInvincible {
		return false
	}

	hc.curHealth = max(0, hc.curHealth-damage)
	// 如果受伤但是没有死亡，需要设置无敌时间
	if hc.IsAlive() && hc.invincibilityDuration > 0.0 {
		hc.setInvincible(hc.invincibilityDuration)
	}
	return true
}

// 是否存活(当前生命值大于 0)
func (hc *HealthComponent) IsAlive() bool {
	return hc.curHealth > 0
}

// 治疗
func (hc *HealthComponent) Heal(amount int) {
	if amount <= 0 || !hc.IsAlive() {
		return
	}
	hc.curHealth = min(hc.maxHealth, hc.curHealth+amount)
}

// 设置无敌状态
func (hc *HealthComponent) setInvincible(duration float64) {
	if duration > 0.0 {
		hc.isInvincible = true
		hc.invincibilityTimer = duration
		return
	}

	// 如果持续时间为 0.0 或负数，则立即取消无敌
	hc.isInvincible = false
	hc.invincibilityTimer = 0.0
}

// 设置最大生命值
func (hc *HealthComponent) SetMaxHealth(maxHealth int) {
	hc.maxHealth = max(1, maxHealth)
	hc.curHealth = min(hc.maxHealth, hc.curHealth)
}

// 设置当前生命值
func (hc *HealthComponent) SetCurHealth(curHealth int) {
	hc.curHealth = max(0, min(hc.curHealth, hc.maxHealth))
}

// 设置无敌状态持续时间
func (hc *HealthComponent) SetInvincibilityDuration(duration float64) {
	hc.invincibilityDuration = max(0.0, duration)
}

// 是否处于无敌状态
func (hc *HealthComponent) IsInvincible() bool {
	return hc.isInvincible
}

// 获取最大生命值
func (hc *HealthComponent) GetMaxHealth() int {
	return hc.maxHealth
}

// 获取当前生命值
func (hc *HealthComponent) GetCurHealth() int {
	return hc.curHealth
}
