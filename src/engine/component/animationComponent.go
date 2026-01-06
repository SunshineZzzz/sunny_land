package component

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
)

// 动画组件
type AnimationComponent struct {
	// 继承组件基类
	Component
	// 动画名称<->动画对象映射
	animations map[string]physics.IAnimation
	// 指向必需的精灵图组件
	spriteComponent *SpriteComponent
	// 当前正在播放的动画
	currentAnimation physics.IAnimation
	// 动画播放计时器
	animationTimer float64
	// 当前是否有动画正在播放
	isPlaying bool
	// 是否在动画结束后删除整个游戏对象
	isOneShotRemoval bool
}

// 确保ColliderComponent实现了IComponent接口
var _ physics.IComponent = (*AnimationComponent)(nil)

// 创建动画组件
func NewAnimationComponent() *AnimationComponent {
	return &AnimationComponent{
		Component: Component{
			ComponentType: def.ComponentTypeAnimation,
		},
		animations:       make(map[string]physics.IAnimation),
		animationTimer:   0,
		isPlaying:        false,
		isOneShotRemoval: false,
	}
}

// 初始化
func (a *AnimationComponent) Init() {
	if a.Owner == nil {
		slog.Error("AnimationComponent requires an Owner")
		return
	}
	a.spriteComponent = a.Owner.GetComponent(def.ComponentTypeSprite).(*SpriteComponent)
	if a.spriteComponent == nil {
		slog.Error("AnimationComponent requires a SpriteComponent")
		return
	}
}

// 更新组件
func (a *AnimationComponent) Update(dt float64, ctx physics.IContext) {
	// 如果没有正在播放的动画，或者没有当前动画，或者没有精灵组件，或者当前动画没有帧，则直接返回
	if !a.isPlaying || a.currentAnimation == nil || a.spriteComponent == nil || a.currentAnimation.IsEmpty() {
		slog.Debug("AnimationComponent Update: currentAnimation is empty")
		return
	}

	// 推进动画计时器
	a.animationTimer += dt

	// 根据时间获取当前帧
	currentFrame := a.currentAnimation.GetFrameAtTime(a.animationTimer)

	// 更新精灵组件的源矩形(使用 SpriteComponent 的新方法)
	a.spriteComponent.SetSourceRect(currentFrame.SourceRect)

	// 检查非循环动画是否已结束
	if !a.currentAnimation.IsLooping() && a.animationTimer >= a.currentAnimation.GetTotalDuration() {
		a.isPlaying = false
		// 将时间限制在结束点
		a.animationTimer = a.currentAnimation.GetTotalDuration()
		if a.isOneShotRemoval {
			a.Owner.SetNeedRemove(true)
		}
	}
}

// 添加一个动画
func (a *AnimationComponent) AddAnimation(animation physics.IAnimation) {
	name := animation.GetName()
	a.animations[name] = animation
	slog.Debug("add animation to gameObject", slog.String("animation.name", name), slog.String("gameOject.name", a.Owner.GetName()))
}

// 播放指定名称的动画
func (a *AnimationComponent) PlayAnimation(name string) {
	animation, ok := a.animations[name]
	if !ok {
		slog.Error("animation name not found", slog.String("animation.name", name))
		return
	}

	// 如果已经在播放相同的动画，不重新开始，注释这一段则重新开始播放
	if a.currentAnimation == animation && a.isPlaying {
		return
	}

	a.currentAnimation = animation
	a.animationTimer = 0.0
	a.isPlaying = true

	// 立即将精灵更新到第一帧
	if a.spriteComponent != nil && !a.currentAnimation.IsEmpty() {
		currentFrame := a.currentAnimation.GetFrameAtTime(0.0)
		a.spriteComponent.SetSourceRect(currentFrame.SourceRect)
	}

	slog.Debug("play animation", slog.String("animation.name", name), slog.String("gameOject.name", a.Owner.GetName()))
}

// 停止播放当前动画
func (a *AnimationComponent) StopAnimation() {
	a.isPlaying = false
}

// 获取当前正在播放的动画名称
func (a *AnimationComponent) GetCurrentAnimationName() string {
	if a.currentAnimation == nil {
		return ""
	}
	return a.currentAnimation.GetName()
}

// 是否动画播放完毕
func (a *AnimationComponent) IsAnimationFinished() bool {
	// 如果没有当前动画(说明从未调用过playAnimation)，或者当前动画是循环的，则返回 false
	if a.currentAnimation == nil || a.currentAnimation.IsLooping() {
		return false
	}
	// 如果动画播放时间超过了总时长，则说明播放完毕
	return a.animationTimer >= a.currentAnimation.GetTotalDuration()
}

// 是否正在播放动画
func (a *AnimationComponent) IsPlaying() bool {
	return a.isPlaying
}

// 设置是否在动画结束后删除整个游戏对象
func (a *AnimationComponent) SetOneShotRemoval(removal bool) {
	a.isOneShotRemoval = removal
}

// 是否在动画结束后删除整个游戏对象
func (a *AnimationComponent) IsOneShotRemoval() bool {
	return a.isOneShotRemoval
}
