package render

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/resource"
	emath "sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 渲染器
type Renderer struct {
	// SDL渲染器
	sdlRenderer *sdl.Renderer
	// 资源管理器
	resourceManager *resource.ResourceManager
}

// 确保Renderer实现了IRenderer接口
var _ physics.IRenderer = (*Renderer)(nil)

// 创建渲染器
func NewRenderer(sdlRenderer *sdl.Renderer, resourceManager *resource.ResourceManager) *Renderer {
	if sdlRenderer == nil {
		panic("sdlRenderer is nil")
	}

	if resourceManager == nil {
		panic("resourceManager is nil")
	}

	renderer := &Renderer{
		sdlRenderer:     sdlRenderer,
		resourceManager: resourceManager,
	}
	renderer.SetDrawColor(0, 0, 0, 255)

	slog.Debug("create renderer")

	return renderer
}

// 绘制精灵图
func (r *Renderer) DrawSprite(camera physics.ICamera, sprite physics.ISprite, position, scale mgl32.Vec2, angle float64) {
	texture := r.resourceManager.GetTexture(sprite.GetTextureId())
	if texture == nil {
		slog.Error("texture is nil", slog.String("textureID", sprite.GetTextureId()))
		return
	}

	// 获取精灵图源矩形
	srcRect := r.GetSpriteSrcRect(sprite)
	if srcRect == nil {
		slog.Error("sourceRect is nil", slog.String("textureID", sprite.GetTextureId()))
		return
	}

	// 应用相机转化
	positionScreen := camera.WorldToScreen(position)
	// 计算目标矩形
	dstRect := sdl.FRect{
		X: positionScreen.X(),
		Y: positionScreen.Y(),
		W: srcRect.W * scale.X(),
		H: srcRect.H * scale.Y(),
	}

	// 视口裁剪
	if !r.IsInViewport(camera, dstRect) {
		return
	}

	flipMode := sdl.FlipNone
	if sprite.GetIsFlipped() {
		flipMode = sdl.FlipHorizontal
	}
	// 执行绘制，默认旋转中心为精灵的中心点
	if !sdl.RenderTextureRotated(r.sdlRenderer, texture, srcRect, &dstRect, angle, nil, flipMode) {
		slog.Error("render texture rotated failed", slog.String("textureID", sprite.GetTextureId()), slog.Any("srcRect", srcRect),
			slog.Float64("angle", angle), slog.Any("flipMode", flipMode))
	}
}

// 绘制视差精灵图
func (r *Renderer) DrawSpriteWithParallax(camera physics.ICamera, sprite physics.ISprite, position, scrollFactor, scale mgl32.Vec2, repeat emath.Vec2B) {
	texture := r.resourceManager.GetTexture(sprite.(*Sprite).textureId)
	if texture == nil {
		slog.Error("texture is nil", slog.String("textureID", sprite.(*Sprite).textureId))
		return
	}

	// 获取精灵图源矩形
	srcRect := r.GetSpriteSrcRect(sprite)
	if srcRect == nil {
		slog.Error("sourceRect is nil", slog.String("textureID", sprite.GetTextureId()))
		return
	}

	// 应用相机转化
	positionScreen := camera.WorldToScreenWithParallax(position, scrollFactor)

	// 计算缩放后的纹理尺寸
	scaledTexW := srcRect.W * scale.X()
	scaledTexH := srcRect.H * scale.Y()

	start := mgl32.Vec2{}
	stop := mgl32.Vec2{}
	viewportSize := camera.GetViewportSize()

	if repeat.X() {
		// 取模保证了背景的起始绘制点永远在[0, width)之间循环
		// 减去一个宽度确保了从屏幕左边界之外就开始绘制，从而完美遮盖屏幕左侧
		start[0] = emath.Mod(positionScreen.X(), scaledTexW) - scaledTexW
		// 从左侧负数坐标开始，一块接一块地向右画，直到覆盖掉屏幕最右侧的一像素为止
		stop[0] = viewportSize.X()
	} else {
		// 一般来说就是普通纹理绘制
		start[0] = positionScreen.X()
		stop[0] = min(positionScreen.X()+scaledTexW, viewportSize.X())
	}
	// 同上
	if repeat.Y() {
		start[1] = emath.Mod(positionScreen.Y(), scaledTexH) - scaledTexH
		stop[1] = viewportSize.Y()
	} else {
		start[1] = positionScreen.Y()
		stop[1] = min(positionScreen.Y()+scaledTexH, viewportSize.Y())
	}

	// 开始绘制
	for y := start.Y(); y < stop.Y(); y += scaledTexH {
		for x := start.X(); x < stop.X(); x += scaledTexW {
			destRect := sdl.FRect{
				X: x,
				Y: y,
				W: scaledTexW,
				H: scaledTexH,
			}
			if !sdl.RenderTexture(r.sdlRenderer, texture, srcRect, &destRect) {
				slog.Error("render parallax texture failed", slog.String("textureID", sprite.GetTextureId()), slog.Any("srcRect", srcRect),
					slog.Any("destRect", destRect))
				return
			}
		}
	}
}

// 绘制用户界面精灵图
func (r *Renderer) DrawUISprite(sprite physics.ISprite, position mgl32.Vec2, size *mgl32.Vec2) {
	texture := r.resourceManager.GetTexture(sprite.GetTextureId())
	if texture == nil {
		slog.Error("texture is nil", slog.String("textureID", sprite.GetTextureId()))
		return
	}

	// 获取精灵图源矩形
	srcRect := r.GetSpriteSrcRect(sprite)
	if srcRect == nil {
		slog.Error("sourceRect is nil", slog.String("textureID", sprite.GetTextureId()))
		return
	}

	// 目标矩形
	dstRect := sdl.FRect{
		X: position.X(),
		Y: position.Y(),
		W: srcRect.W,
		H: srcRect.H,
	}
	if size != nil {
		dstRect.W = size.X()
		dstRect.H = size.Y()
	}

	flipMode := sdl.FlipNone
	if sprite.GetIsFlipped() {
		flipMode = sdl.FlipHorizontal
	}
	// 执行绘制
	if !sdl.RenderTextureRotated(r.sdlRenderer, texture, srcRect, &dstRect, 0.0, nil, flipMode) {
		slog.Error("render ui texture failed", slog.String("textureID", sprite.GetTextureId()), slog.Any("srcRect", srcRect),
			slog.Any("destRect", dstRect))
		return
	}
}

// 设置绘制颜色
func (r *Renderer) SetDrawColorFloat(rc, gc, bc, a float32) {
	if !sdl.SetRenderDrawColorFloat(r.sdlRenderer, rc, gc, bc, a) {
		slog.Error("set render draw color float failed", slog.Any("rc", rc), slog.Any("gc", gc), slog.Any("bc", bc), slog.Any("a", a))
	}
}
func (r *Renderer) SetDrawColor(rc, gc, bc, a uint8) {
	if !sdl.SetRenderDrawColor(r.sdlRenderer, rc, gc, bc, a) {
		slog.Error("set render draw color failed", slog.Any("rc", rc), slog.Any("gc", gc), slog.Any("bc", bc), slog.Any("a", a))
	}
}

// 清屏
func (r *Renderer) ClearScreen() {
	if !sdl.RenderClear(r.sdlRenderer) {
		slog.Error("render clear failed")
	}
}

// 渲染
func (r *Renderer) Present() {
	// 交换缓冲区，将渲染结果显示到屏幕上
	if !sdl.RenderPresent(r.sdlRenderer) {
		slog.Error("render present failed")
	}
}

// 获取精灵图源大小
func (r *Renderer) GetSpriteSrcRect(sprite physics.ISprite) *sdl.FRect {
	srcRect := sprite.GetSourceRect()

	if srcRect != nil {
		// 如果Sprite中存在指定rect，则判断尺寸是否有效
		if srcRect.W <= 0.0 || srcRect.H <= 0.0 {
			slog.Error("sourceRect size is invalid", slog.String("textureID", sprite.GetTextureId()), slog.Any("sourceRect", srcRect))
			return nil
		}
		return srcRect
	}

	// 否则获取纹理尺寸并返回整个纹理大小
	return r.resourceManager.GetTextureSize(sprite.GetTextureId())
}

// 是否在视口内
func (r *Renderer) IsInViewport(camera physics.ICamera, rect sdl.FRect) bool {
	viewportSize := camera.GetViewportSize()
	return rect.X+rect.W >= 0.0 && rect.X <= viewportSize.X() &&
		rect.Y+rect.H >= 0.0 && rect.Y <= viewportSize.Y()
}

/**
 * @brief 绘制填充矩形
 *
 * @param rect 矩形区域
 * @param color 填充颜色
 */
func (r *Renderer) DrawUIFilledRect(rect emath.Rect, color emath.FColor) {
	r.SetDrawColorFloat(color.R, color.G, color.B, color.A)
	sdlRect := sdl.FRect{
		X: rect.Position.X(),
		Y: rect.Position.Y(),
		W: rect.Size.X(),
		H: rect.Size.Y(),
	}
	if !sdl.RenderFillRect(r.sdlRenderer, &sdlRect) {
		slog.Error("render fill rect failed", slog.Any("rect", sdlRect))
	}
	r.SetDrawColorFloat(0.0, 0.0, 0.0, 1.0)
}
