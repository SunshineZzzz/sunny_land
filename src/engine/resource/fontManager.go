package resource

import (
	"fmt"
	"hash/fnv"
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/SunshineZzzz/purego-sdl3/ttf"
)

// 字体键类型
type fontKey struct {
	// 字体文件路径
	path string
	// 字体大小
	size int
}

// 字体键哈希函数
func (fk fontKey) hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fk.path))
	sizeHash := uint64(fk.size)
	return h.Sum64() ^ sizeHash
}

// 字体管理器
type fontManager struct {
	fonts map[uint64]*ttf.Font
}

// 创建字体管理器
func NewFontManager() *fontManager {
	// 初始化TTF
	if !ttf.Init() {
		panic(fmt.Sprintf("ttf init error,%s", sdl.GetError()))
	}

	slog.Debug("font manager init")
	return &fontManager{
		fonts: make(map[uint64]*ttf.Font),
	}
}

// 清理字体管理器
func (fm *fontManager) Clear() {
	for _, font := range fm.fonts {
		ttf.CloseFont(font)
	}
	fm.fonts = make(map[uint64]*ttf.Font)

	ttf.Quit()
	slog.Debug("font manager clear")
}

// 加载字体
func (fm *fontManager) loadFont(path string, size int) *ttf.Font {
	if size <= 0 {
		slog.Error("font size must be greater than 0")
		return nil
	}

	fontKey := fontKey{
		path: path,
		size: size,
	}
	font, ok := fm.fonts[fontKey.hash()]
	if ok {
		return font
	}

	font = ttf.OpenFont(path, float32(size))
	if font == nil {
		slog.Error("open font error", slog.String("path", path), slog.Int("size", size), slog.String("error", sdl.GetError()))
		return nil
	}

	fm.fonts[fontKey.hash()] = font
	slog.Debug("load font size", slog.String("path", path), slog.Int("size", size))
	return font
}

// 获取字体
func (fm *fontManager) GetFont(path string, size int) *ttf.Font {
	fontKey := fontKey{
		path: path,
		size: size,
	}
	font, ok := fm.fonts[fontKey.hash()]
	if ok {
		return font
	}
	slog.Debug("font size not in cache, try to load", slog.String("path", path), slog.Int("size", size))
	return fm.loadFont(path, size)
}

// 卸载字体
func (fm *fontManager) UnloadFont(path string, size int) {
	fontKey := fontKey{
		path: path,
		size: size,
	}
	font, ok := fm.fonts[fontKey.hash()]
	if !ok {
		slog.Warn("font size not in cache, can not unload", slog.String("path", path), slog.Int("size", size))
		return
	}
	ttf.CloseFont(font)
	delete(fm.fonts, fontKey.hash())
	slog.Debug("unload font size", slog.String("path", path), slog.Int("size", size))
}
