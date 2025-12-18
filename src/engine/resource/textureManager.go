package resource

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/img"
	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 纹理管理器
type textureManager struct {
	// 渲染器
	sdlRenderer *sdl.Renderer
	// 纹理缓存
	textures map[string]*sdl.Texture
}

// 创建纹理管理器
func NewTextureManager(renderer *sdl.Renderer) *textureManager {
	if renderer == nil {
		panic("renderer is nil")
	}

	slog.Debug("texture manager init")
	return &textureManager{
		sdlRenderer: renderer,
		textures:    make(map[string]*sdl.Texture),
	}
}

// 清除
func (tm *textureManager) Clear() {
	for _, texture := range tm.textures {
		sdl.DestroyTexture(texture)
	}
	tm.textures = make(map[string]*sdl.Texture)

	slog.Debug("texture manager clear")
}

// 加载纹理
func (tm *textureManager) loadTexture(path string) *sdl.Texture {
	if texture, ok := tm.textures[path]; ok {
		return texture
	}

	texture := img.LoadTexture(tm.sdlRenderer, path)
	if texture == nil {
		slog.Error("load texture error", slog.String("path", path), slog.String("error", sdl.GetError()))
		return nil
	}

	// 载入纹理时，设置纹理缩放模式为最邻近插值(必不可少，否则TileLayer渲染中会出现边缘空隙/模糊)
	if !sdl.SetTextureScaleMode(texture, sdl.ScaleModeNearest) {
		slog.Warn("set texture scale mode error", slog.String("path", path), slog.String("error", sdl.GetError()))
	}

	tm.textures[path] = texture
	return texture
}

// 获取纹理
func (tm *textureManager) GetTexture(path string) *sdl.Texture {
	if texture, ok := tm.textures[path]; ok {
		return texture
	}

	slog.Debug("texture not in cache, try to load", slog.String("path", path))
	return tm.loadTexture(path)
}

// 卸载纹理
func (tm *textureManager) UnloadTexture(path string) {
	if texture, ok := tm.textures[path]; ok {
		sdl.DestroyTexture(texture)
		delete(tm.textures, path)
		slog.Debug("unload texture", slog.String("path", path))
		return
	}
	slog.Warn("texture not in cache , can not unload", slog.String("path", path))
}

// 获取纹理大小
func (tm *textureManager) GetTextureSize(path string) *sdl.FRect {
	texture := tm.GetTexture(path)
	if texture == nil {
		slog.Error("texture not found", slog.String("path", path))
		return nil
	}

	var w, h float32
	sdl.GetTextureSize(texture, &w, &h)
	return &sdl.FRect{X: 0.0, Y: 0.0, W: w, H: h}
}
