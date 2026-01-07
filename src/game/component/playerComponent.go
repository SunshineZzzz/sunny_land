package component

import (
	"log/slog"
	"reflect"

	eComponent "sunny_land/src/engine/component"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
	"sunny_land/src/game/component/state"
)

// 玩家组件，处理玩家输入，状态和控制游戏对象移动的组件
type PlayerComponent struct {
	// 继承基础组件
	eComponent.Component
	// 物理组件
	physicsCom *eComponent.PhysicsComponent
	// 精灵图组件
	spriteCom *eComponent.SpriteComponent
	// 变换组件
	transformCom *eComponent.TransformComponent
	// 动画组件
	animationCom *eComponent.AnimationComponent
	// 生命值组件
	healthCom *eComponent.HealthComponent
	// 当前状态机
	currentState state.IPlayerState
	// 是否死亡
	isDead bool
	// 移动相关参数
	// 水平移动力(1kg像素每二次方秒)
	moveForce float32
	// 最大水平速度(像素每秒)
	maxSpeed float32
	// 摩擦系数
	frictionFactor float32
	// 跳跃速度(像素每秒)
	jumpSpeed float32
	// 攀爬速度(像素每秒)
	climbSpeed float32
	// 属性相关参数
	// 玩家被击中后的硬直时间(单位：秒)
	stunnedDuration float64
	// 土狼时间(Coyote Time): 允许玩家在离地后短暂时间内仍然可以跳跃，单位：秒
	coyoteTime float64
	// 土狼计时器
	coyoteTimeTimer float64
	// 无敌闪烁间隔时间，单位：秒
	flashInterval float64
	// 无敌闪烁计时器
	flashTimer float64
}

// 确保玩家组件实现了IPlayerComponent接口
var _ state.IPlayerComponent = (*PlayerComponent)(nil)

// 确保玩家组件实现了IComponent接口
var _ physics.IComponent = (*PlayerComponent)(nil)

// 创建玩家组件
func NewPlayerComponent() *PlayerComponent {
	return &PlayerComponent{
		Component: eComponent.Component{
			ComponentType: def.ComponentTypePlayer,
		},
		moveForce:       200.0,
		maxSpeed:        120.0,
		frictionFactor:  0.85,
		jumpSpeed:       350.0,
		climbSpeed:      100.0,
		stunnedDuration: 0.4,
		coyoteTime:      0.1,
		flashInterval:   0.1,
	}
}

// 初始化
func (p *PlayerComponent) Init() {
	if p.GetOwner() == nil {
		slog.Error("player component init failed, owner is nil")
		return
	}

	// 获得必要的组件
	p.physicsCom = p.GetOwner().GetComponent(def.ComponentTypePhysics).(*eComponent.PhysicsComponent)
	p.spriteCom = p.GetOwner().GetComponent(def.ComponentTypeSprite).(*eComponent.SpriteComponent)
	p.transformCom = p.GetOwner().GetComponent(def.ComponentTypeTransform).(*eComponent.TransformComponent)
	p.animationCom = p.GetOwner().GetComponent(def.ComponentTypeAnimation).(*eComponent.AnimationComponent)
	p.healthCom = p.GetOwner().GetComponent(def.ComponentTypeHealth).(*eComponent.HealthComponent)

	if p.physicsCom == nil || p.spriteCom == nil || p.transformCom == nil || p.animationCom == nil || p.healthCom == nil {
		slog.Error("player component init failed, physicsCom or spriteCom or transformCom or animationCom or healthCom is nil")
		return
	}

	// 初始化状态机
	p.SetState(state.NewIdleState(p))

	slog.Debug("player component init success")
}

// 试图造成伤害，返回是否成功
func (p *PlayerComponent) TakeDamage(damage int) bool {
	if p.isDead || p.healthCom == nil || damage <= 0 {
		return false
	}

	success := p.healthCom.TakeDamage(damage)
	if !success {
		return false
	}

	// 成功造成伤害了，根据是否存活决定状态切换
	if p.healthCom.IsAlive() {
		// 如果存活，切换到受伤状态
		p.SetState(state.NewHurtState(p))
	} else {
		p.isDead = true
		// 如果死亡了，切换到死亡状态
		p.SetState(state.NewDeadState(p))
	}

	return true
}

// 设置状态
func (p *PlayerComponent) SetState(state state.IPlayerState) {
	if state == nil {
		slog.Warn("player component set state failed, state is nil")
		return
	}

	// 退出当前状态
	if p.currentState != nil {
		p.currentState.Exit()
	}

	p.currentState = state
	slog.Debug("player component set state success", slog.String("state", reflect.TypeOf(state).String()))
	p.currentState.Enter()
}

// 处理输入
func (p *PlayerComponent) HandleInput(ctx physics.IContext) {
	if p.currentState == nil {
		return
	}

	nextState := p.currentState.HandleInput(ctx)
	if nextState != nil {
		p.SetState(nextState)
	}
}

// 更新
func (p *PlayerComponent) Update(dt float64, ctx physics.IContext) {
	if p.currentState == nil {
		return
	}

	// 一旦离地，开始计时 Coyote Timer
	if !p.physicsCom.HasCollidedBelow() {
		p.coyoteTimeTimer += dt
	} else {
		// 如果碰撞到地面，重置 Coyote Timer
		p.coyoteTimeTimer = 0.0
	}

	// 如果处于无敌状态，则进行闪烁
	if p.healthCom.IsInvincible() {
		// 累加帧时间
		p.flashTimer += dt

		// 手动实现类似%取余的循环逻辑，保证计时器永远在 [0, 0.2] 之间波动
		if p.flashTimer >= 2.0*p.flashInterval {
			p.flashTimer -= 2.0 * p.flashInterval
		}

		// 一半时间可见，一半时间不可见
		if p.flashTimer < p.flashInterval {
			// 前0.1秒内不可见
			p.spriteCom.SetHidden(true)
		} else {
			// 后0.1秒内可见
			p.spriteCom.SetHidden(false)
		}
	} else {
		// 不是无敌状态，确保精灵图可见
		p.spriteCom.SetHidden(false)
	}

	nextState := p.currentState.Update(dt, ctx)
	if nextState != nil {
		p.SetState(nextState)
	}
}

// 获取变换组件
func (p *PlayerComponent) GetTransformComponent() *eComponent.TransformComponent {
	return p.transformCom
}

// 获取精灵图组件
func (p *PlayerComponent) GetSpriteComponent() *eComponent.SpriteComponent {
	return p.spriteCom
}

// 获取物理组件
func (p *PlayerComponent) GetPhysicsComponent() *eComponent.PhysicsComponent {
	return p.physicsCom
}

// 获取动画组件
func (p *PlayerComponent) GetAnimationComponent() *eComponent.AnimationComponent {
	return p.animationCom
}

// 获取生命值组件
func (p *PlayerComponent) GetHealthComponent() *eComponent.HealthComponent {
	return p.healthCom
}

// 获取摩擦系数
func (p *PlayerComponent) GetFrictionFactor() float32 {
	return p.frictionFactor
}

// 设置摩擦系数
func (p *PlayerComponent) SetFrictionFactor(frictionFactor float32) {
	p.frictionFactor = frictionFactor
}

// 获取是否死亡
func (p *PlayerComponent) IsDead() bool {
	return p.isDead
}

// 设置是否死亡
func (p *PlayerComponent) SetIsDead(isDead bool) {
	p.isDead = isDead
}

// 获取水平移动力
func (p *PlayerComponent) GetMoveForce() float32 {
	return p.moveForce
}

// 设置水平移动力
func (p *PlayerComponent) SetMoveForce(moveForce float32) {
	p.moveForce = moveForce
}

// 获取最大水平速度
func (p *PlayerComponent) GetMaxSpeed() float32 {
	return p.maxSpeed
}

// 设置最大水平速度
func (p *PlayerComponent) SetMaxSpeed(maxSpeed float32) {
	p.maxSpeed = maxSpeed
}

// 设置跳跃速度
func (p *PlayerComponent) SetJumpSpeed(jumpSpeed float32) {
	p.jumpSpeed = jumpSpeed
}

// 获取跳跃速度
func (p *PlayerComponent) GetJumpSpeed() float32 {
	return p.jumpSpeed
}

// 获取攀爬速度
func (p *PlayerComponent) GetClimbSpeed() float32 {
	return p.climbSpeed
}

// 设置攀爬速度
func (p *PlayerComponent) SetClimbSpeed(climbSpeed float32) {
	p.climbSpeed = climbSpeed
}

// 获取玩家被击中后的硬直时间
func (p *PlayerComponent) GetStunnedDuration() float64 {
	return p.stunnedDuration
}

// 设置玩家被击中后的硬直时间
func (p *PlayerComponent) SetStunnedDuration(stunnedDuration float64) {
	p.stunnedDuration = stunnedDuration
}

// 检查玩家是否在地面上，考虑土狼时间
func (p *PlayerComponent) IsOnGround() bool {
	return p.coyoteTimeTimer <= p.coyoteTime || p.physicsCom.HasCollidedBelow()
}
