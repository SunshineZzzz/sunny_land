package render

import (
	"log/slog"

	"sunny_land/src/engine/resource"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/SunshineZzzz/purego-sdl3/ttf"
	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 使用 SDL_ttf 和 TTF_Text 对象处理文本渲染。
 *
 * 封装 TTF_TextEngine 并提供创建和绘制 TTF_Text 对象的方法，
 * 管理字体加载和颜色设置。
 */
type TextRenderer struct {
	// 持有渲染器的非拥有指针
	sdlRenderer *sdl.Renderer
	// 持有资源管理器的非拥有指针
	resourceManager *resource.ResourceManager
	// 使用SDL3引入的 TTF_TextEngine 来进行绘制
	ttfEngine *ttf.TextEngine
}

// 创建文本渲染器
func NewTextRenderer(sdlRenderer *sdl.Renderer, resourceManager *resource.ResourceManager) *TextRenderer {
	if sdlRenderer == nil {
		panic("sdlRenderer is nil")
	}
	if resourceManager == nil {
		panic("resourceManager is nil")
	}
	return &TextRenderer{
		sdlRenderer:     sdlRenderer,
		resourceManager: resourceManager,
		ttfEngine:       ttf.CreateRendererTextEngine(sdlRenderer),
	}
}

// 关闭
func (tr *TextRenderer) Close() {
	if tr.ttfEngine != nil {
		ttf.DestroyRendererTextEngine(tr.ttfEngine)
		tr.ttfEngine = nil
	}
	// 一定要确保在ResourceManager销毁之后调用
	ttf.Quit()
}

/**
* @brief 绘制UI上的字符串。
*
* @param text UTF-8 字符串内容。
* @param font_id 字体 ID。
* @param font_size 字体大小。
* @param position 左上角屏幕位置。
* @param color 文本颜色。(默认为白色)
 */
func (tr *TextRenderer) DrawUIText(text string, fontId string, fontSize int,
	position mgl32.Vec2, color sdl.FColor) {
	// 从资源管理器获取字体
	font := tr.resourceManager.GetFont(fontId, fontSize)
	if font == nil {
		slog.Warn("font not found", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
		return
	}

	// 创建 TTF_Text 对象
	ttfText := ttf.CreateText(tr.ttfEngine, font, text, 0)
	if ttfText == nil {
		slog.Error("failed to create TTF_Text", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
		return
	}

	// 先渲染一次黑色文字模拟阴影
	ttf.SetTextColorFloat(ttfText, 0.0, 0.0, 0.0, 1.0)
	if !ttf.DrawRendererText(ttfText, position.X()+2.0, position.Y()+2.0) {
		slog.Error("failed to draw TTF_Text", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
	}

	// 然后正常绘制文字
	ttf.SetTextColorFloat(ttfText, color.R, color.G, color.B, color.A)
	if !ttf.DrawRendererText(ttfText, position.X(), position.Y()) {
		slog.Error("failed to draw TTF_Text", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
	}

	// 销毁 TTF_Text 对象
	ttf.DestroyText(ttfText)
}

/**
 * @brief 绘制地图上的字符串。
 *
 * @param camera 相机
 * @param text UTF-8 字符串内容。
 * @param font_id 字体 ID。
 * @param font_size 字体大小。
 * @param position 左上角屏幕位置。
 * @param color 文本颜色。
 */
func (tr *TextRenderer) DrawText(camera *Camera, text string, fontId string, fontSize int,
	position mgl32.Vec2, color sdl.FColor) {
	// 应用相机变换
	screenPosition := camera.WorldToScreen(position)
	// 绘制字符串
	tr.DrawUIText(text, fontId, fontSize, screenPosition, color)
}

/**
 * @brief 获取文本的尺寸。
 *
 * @param text 要测量的文本。
 * @param font_id 字体 ID。
 * @param font_size 字体大小。
 * @return 文本的尺寸。
 */
func (tr *TextRenderer) GetTextSize(text string, fontId string, fontSize int) mgl32.Vec2 {
	// 从资源管理器获取字体
	font := tr.resourceManager.GetFont(fontId, fontSize)
	if font == nil {
		slog.Warn("font not found", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
		return mgl32.Vec2{}
	}

	// 创建 TTF_Text 对象
	ttfText := ttf.CreateText(tr.ttfEngine, font, text, 0)
	if ttfText == nil {
		slog.Error("failed to create TTF_Text", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
		return mgl32.Vec2{}
	}

	// 获取文本尺寸
	var w, h int32
	if !ttf.GetTextSize(ttfText, &w, &h) {
		slog.Error("failed to get TTF_Text size", slog.String("fontId", fontId), slog.Int("fontSize", fontSize))
		return mgl32.Vec2{}
	}

	// 销毁 TTF_Text 对象
	ttf.DestroyText(ttfText)

	return mgl32.Vec2{float32(w), float32(h)}
}
