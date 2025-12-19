package render

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 精灵图
type Sprite struct {
	// 纹理资源标识符
	textureId string
	// 要绘制的纹理部分
	sourceRect *sdl.FRect
	// 是否水平反转
	isFlipped bool
}

// 创建精灵图
func NewSprite(textureId string, sourceRect *sdl.FRect, isFlipped bool) *Sprite {
	slog.Debug("create sprite", slog.Any("textureId", textureId), slog.Any("sourceRect", sourceRect), slog.Any("isFlipped", isFlipped))
	return &Sprite{
		textureId:  textureId,
		sourceRect: sourceRect,
		isFlipped:  isFlipped,
	}
}

// 获取纹理资源标识符
func (s *Sprite) GetTextureId() string {
	return s.textureId
}

// 设置纹理资源标识符
func (s *Sprite) SetTextureId(textureId string) {
	s.textureId = textureId
}

// 获取要绘制的纹理部分
func (s *Sprite) GetSourceRect() *sdl.FRect {
	return s.sourceRect
}

// 设置要绘制的纹理部分
func (s *Sprite) SetSourceRect(sourceRect *sdl.FRect) {
	s.sourceRect = sourceRect
}

// 获取是否水平反转
func (s *Sprite) GetIsFlipped() bool {
	return s.isFlipped
}

// 设置是否水平反转
func (s *Sprite) SetIsFlipped(isFlipped bool) {
	s.isFlipped = isFlipped
}
