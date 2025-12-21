package component

import (
	"log/slog"
	"math"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/render"

	"github.com/go-gl/mathgl/mgl32"
)

// 瓦片类型
type TileType int

const (
	// 空白瓦片
	TileTypeEmpty TileType = iota
	// 普通瓦片
	TileTypeNormal
	// 静止可放置瓦片
	TileTypeSolid
)

// 单个瓦片信息
type TileInfo struct {
	// 精灵图
	Sprite *render.Sprite
	// 瓦片类型
	Type TileType
}

// 创建单个瓦片信息
func NewTileInfo(sprite *render.Sprite, tileType TileType) *TileInfo {
	return &TileInfo{
		Sprite: sprite,
		Type:   tileType,
	}
}

// 瓦片图层组件
type TileLayerComponent struct {
	// 继承基础组件
	Component
	// 单个瓦片的尺寸
	tileSize mgl32.Vec2
	// 地图尺寸，以瓦片为单位
	mapSize mgl32.Vec2
	// 存储所有瓦片信息，按行主序存储，index = y * mapSize.X + x
	tiles []*TileInfo
	// 瓦片层在世界中的偏移量，瓦片层通常不需要缩放及旋转，因此不引入Transform组件
	// offset 最好也保持默认的0，以免增加不必要的复杂性
	offset mgl32.Vec2
	// 是否隐藏
	isHidden bool
}

// 创建瓦片图层组件
func NewTileLayerComponent(tileSize mgl32.Vec2, mapSize mgl32.Vec2, tiles []*TileInfo) *TileLayerComponent {
	slog.Debug("create tile layer component", slog.Any("tileSize", tileSize), slog.Any("mapSize", mapSize), slog.Int("tileCount", len(tiles)))
	return &TileLayerComponent{
		tileSize: tileSize,
		mapSize:  mapSize,
		tiles:    tiles,
		offset:   mgl32.Vec2{0.0, 0.0},
		isHidden: false,
	}
}

// 初始化
func (tlc *TileLayerComponent) Init() {
	if tlc.owner == nil {
		slog.Error("owner is nil")
		return
	}
	slog.Debug("init tile layer component")
}

// 渲染
func (tlc *TileLayerComponent) Render(context *econtext.Context) {
	if tlc.isHidden || tlc.tileSize.X() <= 0.0 || tlc.tileSize.Y() <= 0.0 {
		return
	}

	// 遍历所有瓦片
	for y := 0; y < int(tlc.mapSize.Y()); y++ {
		for x := 0; x < int(tlc.mapSize.X()); x++ {
			// 获取索引
			index := y*int(tlc.mapSize.X()) + x
			// 检查索引有效性和瓦片是否需要渲染
			if index >= len(tlc.tiles) || tlc.tiles[index].Type == TileTypeEmpty {
				continue
			}
			tileInfo := tlc.tiles[index]
			// 计算该瓦片在世界中左上角位置
			leftTopPos := mgl32.Vec2{
				tlc.offset.X() + float32(x)*tlc.tileSize.X(),
				tlc.offset.Y() + float32(y)*tlc.tileSize.Y(),
			}
			// 如果图片大小和瓦片大小不一致，需要调整Y坐标
			if tileInfo.Sprite.GetSourceRect().H != tlc.tileSize.Y() {
				// 目的就是让图片从左下角往上绘制
				leftTopPos[1] -= (tileInfo.Sprite.GetSourceRect().H - tlc.tileSize.Y())
			}
			// 执行绘制
			context.Renderer.DrawSprite(context.Camera, tileInfo.Sprite, leftTopPos, mgl32.Vec2{1.0, 1.0}, 0.0)
		}
	}
}

// 获取指定位置的瓦片信息
func (tlc *TileLayerComponent) GetTileInfoAt(posX, posY int) *TileInfo {
	if posX < 0 || posX >= int(tlc.mapSize.X()) || posY < 0 || posY >= int(tlc.mapSize.Y()) {
		slog.Warn("pos out of range", slog.Int("posX", posX), slog.Int("posY", posY))
		return nil
	}

	index := posY*int(tlc.mapSize.X()) + posX
	if index >= len(tlc.tiles) {
		slog.Warn("index out of range", slog.Int("index", index))
		return nil
	}
	return tlc.tiles[index]
}

// 获取指定位置的瓦片类型，pos不是整数坐标
func (tlc *TileLayerComponent) GetTileTypeAt(posX, posY int) TileType {
	tileInfo := tlc.GetTileInfoAt(posX, posY)
	if tileInfo == nil {
		return TileTypeEmpty
	}
	return tileInfo.Type
}

// 获取指定世界位置的瓦片类型
func (tlc *TileLayerComponent) GetTileTypeAtWorldPos(posXF, posYF float32) TileType {
	// 先将世界位置转换为瓦片位置
	posX := int(math.Floor(float64(posXF) / float64(tlc.tileSize.X())))
	posY := int(math.Floor(float64(posYF) / float64(tlc.tileSize.Y())))
	return tlc.GetTileTypeAt(posX, posY)
}
