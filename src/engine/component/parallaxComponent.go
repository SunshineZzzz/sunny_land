package component

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/utils/def"
	"sunny_land/src/engine/utils/math"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// 视差组件
// 该组件根据相机的位置和滚动因子来移动纹理。在背景中渲染可滚动纹理的组件，以创建视差效果。
type ParallaxComponent struct {
	// 继承基础组件
	Component
	// 缓存变换组件
	transformComponent *TransformComponent
	// 精灵图对象
	sprite *render.Sprite
	// 视差滚动因子，0.0=静止, 1.0=随相机移动, <1.0=比相机慢
	scrollFactor mgl32.Vec2
	// 是否沿着X和Y轴周期性重复
	repeat emath.Vec2B
	// 是否隐藏
	isHidden bool
}

// 确保ParallaxComponent实现了IComponent接口
var _ physics.IComponent = (*ParallaxComponent)(nil)

// 创建视差组件
func NewParallaxComponent(textureId string, scrollFactor mgl32.Vec2, repeat math.Vec2B) *ParallaxComponent {
	slog.Debug("create parallax component", slog.String("textureId", textureId), slog.Any("scrollFactor", scrollFactor), slog.Any("repeat", repeat))
	return &ParallaxComponent{
		Component: Component{
			componentType: def.ComponentTypeParallax,
		},
		// 视差背景默认为整张图片
		sprite:       render.NewSprite(textureId, nil, false),
		scrollFactor: scrollFactor,
		repeat:       repeat,
	}
}

// 初始化
func (pc *ParallaxComponent) Init() {
	if pc.owner == nil {
		slog.Error("parallax component owner is nil")
		return
	}
	pc.transformComponent = pc.owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if pc.transformComponent == nil {
		slog.Error("parallax component transform component is nil")
		return
	}
}

// 更新视差组件
func (pc *ParallaxComponent) Update(float64, physics.IContext) {
}

// 渲染视差组件
func (pc *ParallaxComponent) Render(context physics.IContext) {
	if pc.isHidden || pc.transformComponent == nil {
		return
	}

	// 直接调用视差滚动绘制函数
	context.GetRenderer().DrawSpriteWithParallax(
		context.GetCamera(),
		pc.sprite,
		pc.transformComponent.GetPosition(),
		pc.scrollFactor,
		mgl32.Vec2{1.0, 1.0},
		pc.repeat,
	)
}

// 设置精灵图对象
func (pc *ParallaxComponent) SetSprite(sprite *render.Sprite) {
	pc.sprite = sprite
}

// 设置视差滚动因子
func (pc *ParallaxComponent) SetScrollFactor(scrollFactor mgl32.Vec2) {
	pc.scrollFactor = scrollFactor
}

// 设置是否周期性重复
func (pc *ParallaxComponent) SetRepeat(repeat math.Vec2B) {
	pc.repeat = repeat
}

// 设置是否隐藏，不渲染
func (pc *ParallaxComponent) SetHidden(isHidden bool) {
	pc.isHidden = isHidden
}

// 获取精灵图对象
func (pc *ParallaxComponent) GetSprite() *render.Sprite {
	return pc.sprite
}

// 获取视差滚动因子
func (pc *ParallaxComponent) GetScrollFactor() mgl32.Vec2 {
	return pc.scrollFactor
}

// 获取是否周期性重复
func (pc *ParallaxComponent) GetRepeat() math.Vec2B {
	return pc.repeat
}

// 获取是否隐藏，不渲染
func (pc *ParallaxComponent) GetHidden() bool {
	return pc.isHidden
}
