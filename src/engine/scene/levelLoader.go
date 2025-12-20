package scene

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"sunny_land/src/engine/component"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/utils/math"

	"github.com/bitly/go-simplejson"
	"github.com/go-gl/mathgl/mgl32"
)

// 关卡加载器
type LevelLoader struct {
	// 地图路径
	mapPath string
}

// 创建关卡加载器
func NewLevelLoader() *LevelLoader {
	slog.Debug("LevelLoader created")
	return &LevelLoader{}
}

// 加载关卡数据到指定的Scene对象中
func (ll *LevelLoader) LoadLevel(mapPath string, scene IScene) bool {
	ll.mapPath = mapPath

	// 加载JSON文件
	data, err := os.ReadFile(ll.mapPath)
	if err != nil {
		slog.Error("Failed to read level file", slog.String("mapPath", ll.mapPath), slog.Any("error", err))
		return false
	}

	// 解析JSON数据
	json, err := simplejson.NewJson(data)
	if err != nil {
		slog.Error("Failed to parse level file", slog.String("mapPath", ll.mapPath), slog.Any("error", err))
		return false
	}

	// 加载tiled图层数据
	if json.Get("layers") == nil || json.Get("layers").MustArray() == nil {
		slog.Error("lack layers", slog.String("mapPath", ll.mapPath), slog.Any("error", err))
		return false
	}
	for _, layer := range json.Get("layers").MustArray() {
		layerMap := layer.(map[string]any)
		// 获取个图层对象中的类型(type)字段
		layerType := "none"
		if _, ok := layerMap["type"]; ok {
			layerType = layerMap["type"].(string)
		}
		if !layerMap["visible"].(bool) {
			slog.Info("layer is not visible", slog.String("layerName", layerMap["name"].(string)))
			continue
		}
		// 更具图层类型决定加载方法
		switch layerType {
		case "imagelayer":
			// 图像图层
			ll.loadImageLayer(layerMap, scene)
		case "tilelayer":
			// 图块图层，瓦片图层
			ll.loadTileLayer(layerMap, scene)
		case "objectgroup":
			// 对象图层
			ll.loadObjectLayer(layerMap, scene)
		default:
			slog.Warn("unknown layer type", slog.String("layerType", layerType))
		}
	}

	slog.Info("level loaded", slog.String("mapPath", ll.mapPath))
	return true
}

// 加载图像图层
func (ll *LevelLoader) loadImageLayer(layerMap map[string]any, scene IScene) {
	// 获取图像图层名称
	layerName := "Unnamed"
	if _, ok := layerMap["name"]; ok {
		layerName = layerMap["name"].(string)
	}

	// 获取纹理相对路径
	texturePath := ""
	if _, ok := layerMap["image"]; ok {
		texturePath = layerMap["image"].(string)
	}
	if texturePath == "" {
		slog.Error("image layer lack texture path", slog.String("layerName", layerName))
		return
	}
	// 解析纹理路径为干净的相对路径
	textureId := ll.resolvePath(texturePath)

	// 获取图像偏移量
	offset := mgl32.Vec2{
		0.0,
		0.0,
	}
	if _, ok := layerMap["offsetx"]; ok {
		offsetX, _ := layerMap["offsetx"].(json.Number).Float64()
		offset[0] = float32(offsetX)
	}
	if _, ok := layerMap["offsety"]; ok {
		offsetY, _ := layerMap["offsety"].(json.Number).Float64()
		offset[1] = float32(offsetY)
	}

	// 获取视差因子
	scrollFactor := mgl32.Vec2{
		1.0,
		1.0,
	}
	if _, ok := layerMap["parallaxx"]; ok {
		pallaxX, _ := layerMap["parallaxx"].(json.Number).Float64()
		scrollFactor[0] = float32(pallaxX)
	}
	if _, ok := layerMap["parallaxy"]; ok {
		pallaxY, _ := layerMap["parallaxy"].(json.Number).Float64()
		scrollFactor[1] = float32(pallaxY)
	}

	// 获取重复标志
	repeat := math.Vec2B{
		false,
		false,
	}
	if _, ok := layerMap["repeatx"]; ok {
		repeat[0] = layerMap["repeatx"].(bool)
	}
	if _, ok := layerMap["repeaty"]; ok {
		repeat[1] = layerMap["repeaty"].(bool)
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

// 加载图块图层
func (ll *LevelLoader) loadTileLayer(layerMap map[string]any, scene IScene) {
}

// 加载对象图层
func (ll *LevelLoader) loadObjectLayer(layerMap map[string]any, scene IScene) {
}

// 将相对于地图文件的图像路径转换为干净的相对路径
// path: 地图文件中定义的图片相对路径(例如:"../textures/Layers/middle.png")
func (ll *LevelLoader) resolvePath(path string) string {
	// 获取地图文件的父目录(相对于可执行文件)，"assets/maps/level1.tmj" -> "assets/maps"
	mapDir := filepath.Dir(ll.mapPath)
	// 合并路径并清理(处理"."和"..")，相当于(map_dir/image_path)
	fullPath := filepath.Join(mapDir, path)
	return fullPath
}
