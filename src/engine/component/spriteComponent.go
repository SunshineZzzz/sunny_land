package component

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/resource"
	"sunny_land/src/engine/utils"
	"sunny_land/src/engine/utils/def"
	emath "sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 精灵图组件
type SpriteComponent struct {
	// 继承基础组件
	Component
	// 资源管理器
	resourceManager *resource.ResourceManager
	// 缓存变换组件
	transformComponent *TransformComponent
	// 精灵图对象
	sprite *render.Sprite
	// 对齐方式
	alignment utils.Alignment
	// 精灵图尺寸
	spriteSize mgl32.Vec2
	// 偏移量
	offset mgl32.Vec2
	// 是否隐藏
	isHidden bool
}

// 确保SpriteComponent实现了IComponent接口
var _ physics.IComponent = (*SpriteComponent)(nil)

// 创建精灵图组件
func NewSpriteComponent(textureId string, resourceManager *resource.ResourceManager, alignment utils.Alignment,
	sourceRect *sdl.FRect, isFlipped bool) *SpriteComponent {
	if resourceManager == nil {
		slog.Error("resourceManager is nil")
	}
	slog.Debug("create sprite component", slog.String("textureId", textureId), slog.Any("sourceRect", sourceRect), slog.Bool("isFlipped", isFlipped),
		slog.Any("alignment", alignment))
	return &SpriteComponent{
		Component: Component{
			componentType: def.ComponentTypeSprite,
		},
		resourceManager: resourceManager,
		sprite:          render.NewSprite(textureId, sourceRect, isFlipped),
		alignment:       alignment,
	}
}

// 根据精灵图对象创建精灵图组件
func NewSpriteComponentFromSprite(sprite *render.Sprite, resourceManager *resource.ResourceManager, alignment utils.Alignment) *SpriteComponent {
	if sprite == nil {
		slog.Error("sprite is nil")
		return nil
	}
	if resourceManager == nil {
		slog.Error("resourceManager is nil")
	}
	slog.Debug("create sprite component from sprite", slog.String("textureId", sprite.GetTextureId()),
		slog.Any("sourceRect", sprite.GetSourceRect()), slog.Bool("isFlipped", sprite.GetIsFlipped()))
	return &SpriteComponent{
		Component: Component{
			componentType: def.ComponentTypeSprite,
		},
		resourceManager: resourceManager,
		sprite:          sprite,
		alignment:       alignment,
	}
}

// 初始化
func (sc *SpriteComponent) Init() {
	if sc.owner == nil {
		slog.Error("owner is nil")
		return
	}
	sc.transformComponent = sc.owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if sc.transformComponent == nil {
		slog.Warn("transform component is nil", slog.String("owner", sc.owner.GetName()))
		return
	}

	// 获取大小和偏移
	sc.updateSpriteSize()
	sc.updateOffset()
}

// 更新精灵图尺寸
func (sc *SpriteComponent) updateSpriteSize() {
	if sc.resourceManager == nil {
		slog.Error("resourceManager is nil")
		return
	}

	if sc.sprite.GetSourceRect() != nil {
		sc.spriteSize = mgl32.Vec2{sc.sprite.GetSourceRect().W, sc.sprite.GetSourceRect().H}
	} else {
		textureSize := sc.resourceManager.GetTextureSize(sc.sprite.GetTextureId())
		sc.spriteSize = mgl32.Vec2{textureSize.W, textureSize.H}
	}
}

// 更新偏移量
func (sc *SpriteComponent) updateOffset() {
	// 如果尺寸无效，偏移量为0.0
	if sc.spriteSize[0] <= 0.0 || sc.spriteSize[1] <= 0.0 {
		sc.offset = mgl32.Vec2{0.0, 0.0}
		return
	}
	// 获取缩放
	scale := sc.transformComponent.GetScale()
	// 计算精灵图左上角对于变换位置的偏移
	switch sc.alignment {
	case utils.AlignTopLeft:
		sc.offset = mgl32.Vec2{0.0, 0.0}
	case utils.AlignTopCenter:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X() * 0.5, 0.0}, scale)
	case utils.AlignTopRight:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X(), 0.0}, scale)
	case utils.AlignCenterLeft:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{0.0, -sc.spriteSize.Y() * 0.5}, scale)
	case utils.AlignCenter:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X() * 0.5, -sc.spriteSize.Y() * 0.5}, scale)
	case utils.AlignCenterRight:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X(), -sc.spriteSize.Y() * 0.5}, scale)
	case utils.AlignBottomLeft:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{0.0, -sc.spriteSize.Y()}, scale)
	case utils.AlignBottomCenter:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X() * 0.5, -sc.spriteSize.Y()}, scale)
	case utils.AlignBottomRight:
		sc.offset = emath.Mgl32Vec2MulElem(mgl32.Vec2{-sc.spriteSize.X(), -sc.spriteSize.Y()}, scale)
	}
}

// 设置对齐方式
func (sc *SpriteComponent) SetAlignment(anchor utils.Alignment) {
	sc.alignment = anchor
	sc.updateOffset()
}

// 设置源矩形
func (sc *SpriteComponent) SetSourceRect(sourceRect *sdl.FRect) {
	sc.sprite.SetSourceRect(sourceRect)
	sc.updateSpriteSize()
	sc.updateOffset()
}

// 根据Id设置精灵图
func (sc *SpriteComponent) SetSpriteById(textureId string, sourceRect *sdl.FRect) {
	sc.sprite.SetTextureId(textureId)
	sc.sprite.SetSourceRect(sourceRect)

	sc.updateSpriteSize()
	sc.updateOffset()
}

// 渲染
func (sc *SpriteComponent) Render(context physics.IContext) {
	if sc.isHidden || sc.transformComponent == nil || sc.resourceManager == nil {
		return
	}

	// 获取变换信息，并且考虑偏移量
	transform := sc.transformComponent.GetPosition().Add(sc.offset)
	scale := sc.transformComponent.GetScale()
	rotationDegrees := sc.transformComponent.GetRotation()

	// 执行绘制
	context.GetRenderer().DrawSprite(context.GetCamera(), sc.sprite, transform, scale, rotationDegrees)
}

// 更新组件状态
func (sc *SpriteComponent) Update(float64, physics.IContext) {
}
