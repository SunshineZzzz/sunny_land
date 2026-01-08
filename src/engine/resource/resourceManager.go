package resource

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/SunshineZzzz/purego-sdl3/ttf"
)

// 资源管理器
type ResourceManager struct {
	// 字体管理器
	fontManager *fontManager
	// 纹理管理器
	textureManager *textureManager
	// 音频管理器
	audioManager *audioManager
}

// 创建资源管理器
func NewResourceManager(renderer *sdl.Renderer) *ResourceManager {
	slog.Debug("resource manager init")
	return &ResourceManager{
		fontManager:    NewFontManager(),
		textureManager: NewTextureManager(renderer),
		audioManager:   NewAudioManager(),
	}
}

// 清理资源管理器
func (rm *ResourceManager) Clear() {
	rm.fontManager.Clear()
	rm.textureManager.Clear()
	rm.audioManager.Clear()

	slog.Debug("resource manager clear")
}

// 获取纹理
func (rm *ResourceManager) GetTexture(path string) *sdl.Texture {
	return rm.textureManager.GetTexture(path)
}

// 卸载纹理
func (rm *ResourceManager) UnloadTexture(path string) {
	rm.textureManager.UnloadTexture(path)
}

// 获取字体
func (rm *ResourceManager) GetFont(path string, size int) *ttf.Font {
	return rm.fontManager.GetFont(path, size)
}

// 卸载字体
func (rm *ResourceManager) UnloadFont(path string, size int) {
	rm.fontManager.UnloadFont(path, size)
}

// 获取音效
func (rm *ResourceManager) GetSound(path string) *[]IAudio {
	return rm.audioManager.GetSound(path)
}

// 加载音效
func (rm *ResourceManager) LoadSound(path string) *[]IAudio {
	return rm.audioManager.loadSound(path)
}

// 卸载音效
func (rm *ResourceManager) UnloadSound(path string) {
	rm.audioManager.UnloadSound(path)
}

// 获取音乐
func (rm *ResourceManager) GetMusic(path string) *[]IAudio {
	return rm.audioManager.GetMusic(path)
}

// 加载音乐
func (rm *ResourceManager) LoadMusic(path string) *[]IAudio {
	return rm.audioManager.loadMusic(path)
}

// 卸载音乐
func (rm *ResourceManager) UnloadMusic(path string) {
	rm.audioManager.UnloadMusic(path)
}

// 获取纹理大小
func (rm *ResourceManager) GetTextureSize(path string) *sdl.FRect {
	return rm.textureManager.GetTextureSize(path)
}

// 设置音效音量
func (rm *ResourceManager) SetSoundVolume(volume float32) {
	rm.audioManager.SetSoundVolume(volume)
}

// 获取音效音量
func (rm *ResourceManager) GetSoundVolume() float32 {
	return rm.audioManager.GetSoundVolume()
}

// 设置音乐音量
func (rm *ResourceManager) SetMusicVolume(volume float32) {
	rm.audioManager.SetMusicVolume(volume)
}

// 获取音乐音量
func (rm *ResourceManager) GetMusicVolume() float32 {
	return rm.audioManager.GetMusicVolume()
}
