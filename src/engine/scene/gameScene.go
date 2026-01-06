package scene

import (
	"log/slog"

	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/utils/def"
	emath "sunny_land/src/engine/utils/math"
	gcomponent "sunny_land/src/game/component"

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

	// 初始化关卡
	if !gs.InitLevel() {
		slog.Error("level init failed")
		gs.GetContext().InputManager.SetShouldQuit(true)
		return
	}

	// 初始化玩家
	if !gs.InitPlayer() {
		slog.Error("player init failed")
		gs.GetContext().InputManager.SetShouldQuit(true)
		return
	}

	// 初始化敌人和道具
	if !gs.InitEnemiesAndItem() {
		slog.Error("enemies and item init failed")
		gs.GetContext().InputManager.SetShouldQuit(true)
		return
	}

	slog.Debug("GameScene initialized", slog.String("sceneName", gs.sceneName))
}

// 初始化关卡
func (gs *GameScene) InitLevel() bool {
	// 加载关卡
	if !NewLevelLoader().LoadLevel("assets/maps/level1.tmj", gs) {
		slog.Error("level1.tmj load failed")
		return false
	}

	// 注册场景中的main层到物理引擎，main层会有物理属性
	mainLayer := gs.FindGameObjectByName("main")
	if mainLayer == nil {
		slog.Error("main layer not found")
		return false
	}

	tileLayerComp := mainLayer.GetComponent(def.ComponentTypeTileLayer).(*component.TileLayerComponent)
	if tileLayerComp == nil {
		slog.Error("main layer tile layer component not found")
		return false
	}

	gs.ctx.PhysicsEngine.RegisterTileLayerComponent(tileLayerComp)
	slog.Info("main layer registered to physics engine")

	// 世界大小
	worldSize := mainLayer.GetComponent(def.ComponentTypeTileLayer).(*component.TileLayerComponent).GetWorldSize()
	// 设置相机限制范围
	gs.ctx.Camera.SetLimitBounds(&emath.Rect{Position: mgl32.Vec2{0.0, 0.0}, Size: worldSize})

	// 设置世界边界
	gs.ctx.PhysicsEngine.SetWorldBounds(&emath.Rect{Position: mgl32.Vec2{0.0, 0.0}, Size: worldSize})

	slog.Debug("GameScene level initialized", slog.String("sceneName", gs.sceneName))
	return true
}

// 初始化玩家
func (gs *GameScene) InitPlayer() bool {
	// 获取玩家对象
	gs.playerObject = gs.FindGameObjectByName("player")
	if gs.playerObject == nil {
		slog.Error("player object not found")
		return false
	}

	// 添加PlayerComponent到玩家对象
	playerCom := gcomponent.NewPlayerComponent()
	if gs.playerObject.AddComponent(playerCom) == nil {
		slog.Error("player component init failed, playerComponent is nil")
		return false
	}

	// 相机目标跟踪玩家
	transformComp := gs.playerObject.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent)
	if transformComp == nil {
		slog.Error("player object transform component not found")
		return false
	}
	gs.ctx.Camera.SetTargetTC(transformComp)

	slog.Debug("player object transform component set to camera target")
	return true
}

// 初始化敌人和道具
func (gs *GameScene) InitEnemiesAndItem() bool {
	success := true
	for e := gs.GameObjects.Front(); e != nil; e = e.Next() {
		gt := e.Value.(*object.GameObject)
		switch gt.GetName() {
		case "eagle":
			ac := gt.GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
			if ac == nil {
				slog.Error("eagle object animation component not found")
				success = false
			}
			ac.PlayAnimation("fly")
		case "frog":
			ac := gt.GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
			if ac == nil {
				slog.Error("frog object animation component not found")
				success = false
			}
			ac.PlayAnimation("idle")
		case "opossum":
			ac := gt.GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
			if ac == nil {
				slog.Error("opossum object animation component not found")
				success = false
			}
			ac.PlayAnimation("walk")
		}

		if gt.GetTag() == "item" {
			ac := gt.GetComponent(def.ComponentTypeAnimation).(*component.AnimationComponent)
			if ac == nil {
				slog.Error("item object animation component not found")
				success = false
			}
			ac.PlayAnimation("idle")
		}
	}
	return success
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
	// 测试玩家健康值
	gs.testHealth()
}

// 清理
func (gs *GameScene) Clean() {
	gs.scene.Clean()
}

// 测试玩家健康值
func (gs *GameScene) testHealth() {
	inputManager := gs.GetContext().InputManager
	if inputManager.IsActionPressed("attack") {
		playerCom := gs.playerObject.GetComponent(def.ComponentTypePlayer).(*gcomponent.PlayerComponent)
		playerCom.TakeDamage(1)
	}
}
