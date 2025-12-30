package scene

import (
	"log/slog"
	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 游戏场景，主要的游戏场景，包含玩家、敌人、关卡元素等。
type GameScene struct {
	// 继承基础场景
	scene
	// 保存玩家对象指针，方便访问
	playerObject *object.GameObject
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

	// 加载关卡
	NewLevelLoader().LoadLevel("assets/maps/level1.tmj", gs)
	// 注册场景中的main层到物理引擎，main层会有物理属性
	mainLayer := gs.FindGameObjectByName("main")
	if mainLayer != nil {
		tileLayerComp := mainLayer.GetComponent(&component.TileLayerComponent{}).(*component.TileLayerComponent)
		if tileLayerComp != nil {
			gs.ctx.PhysicsEngine.RegisterTileLayerComponent(tileLayerComp)
			slog.Info("main layer registered to physics engine")
		}
	}

	// 获取玩家对象
	gs.playerObject = gs.FindGameObjectByName("player")
	if gs.playerObject == nil {
		slog.Error("player object not found")
		return
	}

	slog.Debug("GameScene initialized", slog.String("sceneName", gs.sceneName))
}

// 更新
func (gs *GameScene) Update(dt float64) {
	// 测试更新摄像机
	// gs.testCamera()
	// 测试更新玩家
	gs.TestPlayer()
	// 测试碰撞对
	gs.TestCollisionPairs()

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

// 测试摄像机
func (gs *GameScene) testCamera() {
	key_state := sdl.GetKeyboardState()
	if key_state[sdl.ScancodeW] {
		// 摄像机向上运动
		gs.ctx.Camera.Move(mgl32.Vec2{0, -1})
	}
	if key_state[sdl.ScancodeS] {
		// 摄像机向下运动
		gs.ctx.Camera.Move(mgl32.Vec2{0, 1})
	}
	if key_state[sdl.ScancodeA] {
		// 摄像机向左运动
		gs.ctx.Camera.Move(mgl32.Vec2{-1, 0})
	}
	if key_state[sdl.ScancodeD] {
		// 摄像机向右运动
		gs.ctx.Camera.Move(mgl32.Vec2{1, 0})
	}
}

// 测试游戏对象
func (gs *GameScene) TestPlayer() {
	if gs.playerObject == nil {
		return
	}
	inputManager := gs.ctx.InputManager
	physicsComp := gs.playerObject.GetComponent(&component.PhysicsComponent{}).(*component.PhysicsComponent)
	if physicsComp == nil {
		return
	}
	if inputManager.IsActionDown("move_left") {
		physicsComp.Velocity[0] = -100.0
	} else {
		physicsComp.Velocity[0] *= 0.9
	}

	if inputManager.IsActionDown("move_right") {
		physicsComp.Velocity[0] = 100.0
	} else {
		physicsComp.Velocity[0] *= 0.9
	}

	if inputManager.IsActionPressed("jump") {
		physicsComp.Velocity[1] = -400.0
	}
}

// 测试碰撞组件对
func (gs *GameScene) TestCollisionPairs() {
	collision_pairs := gs.ctx.PhysicsEngine.GetCollisionPairs()
	for _, pair := range collision_pairs {
		slog.Info("碰撞对", slog.String("a", pair.A.(*object.GameObject).GetName()), slog.String("b", pair.B.(*object.GameObject).GetName()))
	}
}
