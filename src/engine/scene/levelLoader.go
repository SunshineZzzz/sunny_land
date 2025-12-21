package scene

import (
	"log/slog"
	"os"
	"path/filepath"
	"sunny_land/src/engine/component"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/bitly/go-simplejson"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/go-gl/mathgl/mgl32"
)

// 关卡加载器
type LevelLoader struct {
	// 地图路径
	mapPath string
	// 地图尺寸，瓦片数量
	mapSize mgl32.Vec2
	// 单个瓦片的尺寸
	tileSize mgl32.Vec2
	// 瓦片集数据
	tilesetsData *rbt.Tree
}

// 创建关卡加载器
func NewLevelLoader() *LevelLoader {
	slog.Debug("LevelLoader created")
	return &LevelLoader{
		tilesetsData: rbt.NewWithIntComparator(),
	}
}

// 加载关卡数据到指定的Scene对象中
func (ll *LevelLoader) LoadLevel(mapPath string, scene IScene) bool {
	// 加载JSON文件
	data, err := os.ReadFile(mapPath)
	if err != nil {
		slog.Error("Failed to read level file", slog.String("mapPath", mapPath), slog.Any("error", err))
		return false
	}

	// 解析JSON数据
	root, err := simplejson.NewJson(data)
	if err != nil {
		slog.Error("Failed to parse level file", slog.String("mapPath", mapPath), slog.Any("error", err))
		return false
	}

	ll.mapPath = mapPath
	ll.mapSize = mgl32.Vec2{
		float32(root.Get("width").MustInt(0)),
		float32(root.Get("height").MustInt(0)),
	}
	ll.tileSize = mgl32.Vec2{
		float32(root.Get("tilewidth").MustInt(0)),
		float32(root.Get("tileheight").MustInt(0)),
	}

	// 加载tilesets数据
	if root.Get("tilesets") == nil || root.Get("tilesets").MustArray() == nil {
		slog.Error("lack tilesets", slog.String("mapPath", ll.mapPath))
		return false
	}
	tilesets := root.Get("tilesets")
	// 遍历所有tileset
	for i := 0; i < len(tilesets.MustArray()); i++ {
		tileset := tilesets.GetIndex(i)
		if tileset.Get("firstgid") == nil {
			slog.Error("lack firstgid", slog.String("mapPath", ll.mapPath))
			continue
		}
		if tileset.Get("source") == nil {
			slog.Error("lack source", slog.String("mapPath", ll.mapPath))
			continue
		}
		tilesetPath := ll.resolvePath(tileset.Get("source").MustString(""), ll.mapPath)
		firstGId := tileset.Get("firstgid").MustInt(0)
		// 加载图块集
		ll.loadTileSet(tilesetPath, firstGId)
	}

	// 加载tiled图层数据
	if root.Get("layers") == nil || root.Get("layers").MustArray() == nil {
		slog.Error("lack layers", slog.String("mapPath", ll.mapPath))
		return false
	}
	layers := root.Get("layers")
	for i := 0; i < len(layers.MustArray()); i++ {
		layer := layers.GetIndex(i)
		// 获取个图层对象中的类型(type)字段
		layerType := layer.Get("type").MustString("none")
		if !layer.Get("visible").MustBool(false) {
			slog.Info("layer is not visible", slog.String("layerName", layer.Get("name").MustString("Unnamed")))
			continue
		}
		// 更具图层类型决定加载方法
		switch layerType {
		case "imagelayer":
			// 图像图层
			ll.loadImageLayer(layer, scene)
		case "tilelayer":
			// 图块图层，瓦片图层
			ll.loadTileLayer(layer, scene)
		case "objectgroup":
			// 对象图层
			ll.loadObjectLayer(layer, scene)
		default:
			slog.Warn("unknown layer type", slog.String("layerType", layerType))
		}
	}

	slog.Info("level loaded", slog.String("mapPath", ll.mapPath))
	return true
}

// 加载图像图层
func (ll *LevelLoader) loadImageLayer(layer *simplejson.Json, scene IScene) {
	// 获取图像图层名称
	layerName := layer.Get("name").MustString("Unnamed")

	// 获取纹理相对路径
	texturePath := layer.Get("image").MustString("")
	if texturePath == "" {
		slog.Error("image layer lack texture path", slog.String("layerName", layerName))
		return
	}
	// 解析纹理路径为干净的相对路径
	textureId := ll.resolvePath(texturePath, ll.mapPath)

	// 获取图像偏移量
	offset := mgl32.Vec2{
		float32(layer.Get("offsetx").MustFloat64(0.0)),
		float32(layer.Get("offsety").MustFloat64(0.0)),
	}
	// 获取视差因子
	scrollFactor := mgl32.Vec2{
		float32(layer.Get("parallaxx").MustFloat64(1.0)),
		float32(layer.Get("parallaxy").MustFloat64(1.0)),
	}

	// 获取重复标志
	repeat := math.Vec2B{
		layer.Get("repeatx").MustBool(false),
		layer.Get("repeaty").MustBool(false),
	}

	// 创建游戏对象
	gameObject := object.NewGameObject(layerName, layerName)
	// 依次添加变换组件，视差组件
	transformComp := component.NewTransformComponent(offset, mgl32.Vec2{1.0, 1.0}, 0.0)
	parallaxComp := component.NewParallaxComponent(textureId, scrollFactor, repeat)
	gameObject.AddComponent(transformComp)
	gameObject.AddComponent(parallaxComp)
	// 添加到场景中
	scene.AddGameObject(gameObject)
	slog.Info("image layer loaded", slog.String("layerName", layerName))
}

// 加载图块集
func (ll *LevelLoader) loadTileSet(tilesetPath string, firstGId int) {
	data, err := os.ReadFile(tilesetPath)
	if err != nil {
		slog.Error("Failed to read tileset file", slog.String("tilesetPath", tilesetPath), slog.Any("error", err))
		return
	}

	root, err := simplejson.NewJson(data)
	if err != nil {
		slog.Error("Failed to parse tileset file", slog.String("tilesetPath", tilesetPath), slog.Any("error", err))
		return
	}

	root.Set("file_path", tilesetPath)
	ll.tilesetsData.Put(firstGId, root)
	slog.Info("tileset loaded", slog.String("tilesetPath", tilesetPath), slog.Int("firstGId", firstGId))
}

// 加载图块图层
func (ll *LevelLoader) loadTileLayer(layer *simplejson.Json, scene IScene) {
	if layer == nil || layer.Get("data") == nil || layer.Get("data").MustArray() == nil {
		slog.Error("lack tile layer data", slog.String("layerName", layer.Get("name").MustString("Unnamed")))
		return
	}

	// 准备瓦片信息切片，瓦片数量 = 地图宽度 * 地图高度
	tileInfos := make([]*component.TileInfo, 0, int(ll.mapSize.X()*ll.mapSize.Y()))
	// 遍历图块数据数组
	datas := layer.Get("data")
	for i := 0; i < len(datas.MustArray()); i++ {
		data := datas.GetIndex(i)
		tileInfos = append(tileInfos, ll.getTileInfoByGId(data.MustInt()))
	}

	// 获取图块图层名称
	layerName := layer.Get("name").MustString("Unnamed")
	// 创建游戏对象
	gameObject := object.NewGameObject(layerName, layerName)
	// 创建瓦片图层组件
	tileLayerComp := component.NewTileLayerComponent(ll.tileSize, ll.mapSize, tileInfos)
	// 游戏对象添加组件
	gameObject.AddComponent(tileLayerComp)
	// 游戏对象添加到场景中
	scene.AddGameObject(gameObject)
	slog.Info("tile layer loaded", slog.String("layerName", layerName))
}

// 根据GId获取瓦片信息
func (ll *LevelLoader) getTileInfoByGId(gId int) *component.TileInfo {
	if gId == 0 {
		// 空白瓦片
		return &component.TileInfo{
			Type: component.TileTypeEmpty,
		}
	}

	// 返回第一个小于等于gId的元素
	entry, found := ll.tilesetsData.Floor(gId)
	if !found {
		slog.Error("no tileset data for gId", slog.Int("gId", gId))
		return &component.TileInfo{
			Type: component.TileTypeEmpty,
		}
	}

	// 对应的是tiled中整张图块集中数组下标(行主序)，多张图片集合对应json数据中tiles数据对象中id等于localId的内容
	localId := gId - entry.Key.(int)
	// 对应的图集(整张或者多张图集)json数据
	root := entry.Value.(*simplejson.Json)
	tilesetPath := root.Get("file_path").MustString("")
	if tilesetPath == "" {
		slog.Error("tileset path is empty", slog.Int("gId", gId))
		return &component.TileInfo{
			Type: component.TileTypeEmpty,
		}
	}
	// 图块集分为两种，整张图块集和多张图块集
	if _, ok := root.CheckGet("image"); ok {
		// 整张图块集
		// 获取图片路径
		textureId := ll.resolvePath(root.Get("image").MustString(""), tilesetPath)
		// 计算瓦片在图片网格中的坐标
		coordinateX := localId % root.Get("columns").MustInt(0)
		coordinateY := localId / root.Get("columns").MustInt(0)
		// 计算瓦片在图片中的像素坐标
		textureRect := &sdl.FRect{
			X: ll.tileSize.X() * float32(coordinateX),
			Y: ll.tileSize.Y() * float32(coordinateY),
			W: ll.tileSize.X(),
			H: ll.tileSize.Y(),
		}
		// 目前只完成渲染，以后再考虑瓦片类型
		return &component.TileInfo{
			Type:   component.TileTypeNormal,
			Sprite: render.NewSprite(textureId, textureRect, false),
		}
	} else {
		// 多张图块集
		if root.Get("tiles") == nil {
			slog.Error("tileset data for gId has no tiles", slog.Int("gId", gId))
			return &component.TileInfo{
				Type: component.TileTypeEmpty,
			}
		}
		// 遍历图块数据数组
		tiles := root.Get("tiles")
		for i := 0; i < len(tiles.MustArray()); i++ {
			tile := tiles.GetIndex(i)
			if tile.Get("id").MustInt(0) == localId {
				if tile.Get("image") == nil {
					slog.Error("tile data for gId has no image", slog.Int("gId", gId))
					return &component.TileInfo{
						Type: component.TileTypeEmpty,
					}
				}
				// 获取图片路径
				textureId := ll.resolvePath(tile.Get("image").MustString(""), tilesetPath)
				// 确认图片尺寸
				imageWidth := tile.Get("imagewidth").MustInt(0)
				imageHeight := tile.Get("imageheight").MustInt(0)
				// 从json中获取源矩形信息
				textureRect := &sdl.FRect{
					X: float32(tile.Get("x").MustFloat64(0.0)),
					Y: float32(tile.Get("y").MustFloat64(0.0)),
					W: float32(tile.Get("width").MustFloat64(float64(imageWidth))),
					H: float32(tile.Get("height").MustFloat64(float64(imageHeight))),
				}
				// 目前只完成渲染，以后再考虑瓦片类型
				return &component.TileInfo{
					Type:   component.TileTypeNormal,
					Sprite: render.NewSprite(textureId, textureRect, false),
				}
			}
		}
	}
	slog.Error("not find tile data for gId", slog.Int("gId", gId))
	return &component.TileInfo{
		Type: component.TileTypeEmpty,
	}
}

// 加载对象图层
func (ll *LevelLoader) loadObjectLayer(layer *simplejson.Json, scene IScene) {
}

/**
 * @brief 解析图片路径，合并地图路径和相对路径。例如：
 * 1. 文件路径："assets/maps/level1.tmj"
 * 2. 相对路径："../textures/Layers/back.png"
 * 3. 最终路径："assets/textures/Layers/back.png"
 * @param relativePath 相对路径（相对于文件）
 * @param filePath 文件路径
 * @return string 解析后的完整路径。
 */
func (ll *LevelLoader) resolvePath(relativePath, filePath string) string {
	// 获取地图文件的父目录(相对于可执行文件)，"assets/maps/level1.tmj" -> "assets/maps"
	mapDir := filepath.Dir(filePath)
	// 合并路径并清理(处理"."和"..")，相当于(map_dir/image_path)
	fullPath := filepath.Join(mapDir, relativePath)
	return fullPath
}
