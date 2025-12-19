package scene

import (
	"log/slog"
	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/utils"

	"github.com/go-gl/mathgl/mgl32"
)

// 游戏场景，主要的游戏场景，包含玩家、敌人、关卡元素等。
type GameScene struct {
	// 继承基础场景
	scene
}

// 创建游戏场景
func NewGameScene(sceneName string, ctx *econtext.Context, sceneManager *SceneManager) *GameScene {
	gs := &GameScene{}
	buildScene(&gs.scene, sceneName, ctx, sceneManager)
	slog.Debug("GameScene created", slog.String("sceneName", sceneName))
	return gs
}

// 初始化游戏场景
func (gs *GameScene) Init() {
	gs.scene.Init()
	slog.Debug("GameScene initialized", slog.String("sceneName", gs.sceneName))

	gs.createTestObject()
}

// 更新
func (gs *GameScene) Update(dt float64) {
	gs.scene.Update(dt)
}

// 渲染
func (gs *GameScene) Render() {
	gs.scene.Render()
}

// 处理事件
func (gs *GameScene) HandleInput() {
	gs.scene.HandleInput()
}

// 清理
func (gs *GameScene) Clean() {
	gs.scene.Clean()
}

// 创建测试游戏对象
func (gs *GameScene) createTestObject() {
	// 创建一个测试游戏对象
	gt := object.NewGameObject("TestObject", "TestObject")
	transform := component.NewTransformComponent(mgl32.Vec2{100.0, 100.0}, mgl32.Vec2{1.0, 1.0}, 0.0)
	sprite := component.NewSpriteComponent("assets/textures/Props/big-crate.png", gs.ctx.ResourceManager, utils.AlignCenter, nil, false)
	gt.AddComponent(transform)
	gt.AddComponent(sprite)
	// 将创建好的游戏对象添加到场景中
	gs.scene.AddGameObject(gt)
}
