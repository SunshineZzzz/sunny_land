package scene

import (
	"log/slog"

	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/utils"
	"sunny_land/src/engine/utils/def"
	emath "sunny_land/src/engine/utils/math"
	gcomponent "sunny_land/src/game/component"

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
	// 处理游戏对象间的碰撞事件
	gs.handleObjectCollisions()
	// 处理瓦片触发事件
	gs.handleTileTriggers()
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

// 处理瓦片触发事件，从PhysicsEngine获取信息
func (gs *GameScene) handleTileTriggers() {
	// 从物理引擎获取触发事件
	triggerEvents := gs.ctx.PhysicsEngine.GetTileTriggerEvents()
	for _, event := range triggerEvents {
		// 瓦片触发事件的对象
		obj := event.GameObject.(*object.GameObject)
		// 瓦片类型
		tileType := event.TileType
		if tileType == physics.TileTypeHazard {
			// 处理玩家与"hazard"对象的碰撞
			if obj.GetName() == "player" {
				obj.GetComponent(def.ComponentTypePlayer).(*gcomponent.PlayerComponent).TakeDamage(1)
			}
			// TODO: 其他对象类型的处理，目前让敌人无视瓦片伤害
		}
	}
}

// 处理游戏对象间的碰撞逻辑，从PhysicsEngine获取信息
func (gs *GameScene) handleObjectCollisions() {
	// 从物理引擎获取碰撞对
	collisionPairs := gs.ctx.PhysicsEngine.GetCollisionPairs()
	for _, pair := range collisionPairs {
		obj1 := pair.A.(*object.GameObject)
		obj2 := pair.B.(*object.GameObject)

		// 处理玩家与敌人的碰撞
		if obj1.GetTag() == "player" && obj2.GetTag() == "enemy" {
			gs.playerVSEnemyCollision(obj1, obj2)
		} else if obj1.GetTag() == "enemy" && obj2.GetTag() == "player" {
			gs.playerVSEnemyCollision(obj2, obj1)
		}
		// 处理玩家与道具的碰撞
		if obj1.GetTag() == "player" && obj2.GetTag() == "item" {
			gs.playerVSItemCollision(obj1, obj2)
		} else if obj1.GetTag() == "item" && obj2.GetTag() == "player" {
			gs.playerVSItemCollision(obj2, obj1)
		}
		// 处理玩家与"hazard"对象的碰撞
		if obj1.GetTag() == "player" && obj2.GetTag() == "hazard" {
			obj1.GetComponent(def.ComponentTypePlayer).(*gcomponent.PlayerComponent).TakeDamage(1)
		} else if obj1.GetTag() == "hazard" && obj2.GetTag() == "player" {
			obj2.GetComponent(def.ComponentTypePlayer).(*gcomponent.PlayerComponent).TakeDamage(1)
		}
	}
}

// 玩家与敌人碰撞处理
func (gs *GameScene) playerVSEnemyCollision(player, enemy *object.GameObject) {
	// 踩踏判断逻辑：1. 玩家中心点在敌人上方  2. 重叠区域：overlap.x > overlap.y
	playerAABB := player.GetComponent(def.ComponentTypeCollider).(*component.ColliderComponent).GetWorldAABB()
	enemyAABB := enemy.GetComponent(def.ComponentTypeCollider).(*component.ColliderComponent).GetWorldAABB()
	playerCenter := playerAABB.Position.Add(playerAABB.Size.Mul(0.5))
	enemyCenter := enemyAABB.Position.Add(enemyAABB.Size.Mul(0.5))
	overlap := playerAABB.Size.Mul(0.5).Add(enemyAABB.Size.Mul(0.5)).Sub(emath.Mgl32Vec2ABS(playerCenter, enemyCenter))

	// 踩踏判断成功，敌人受伤
	if playerCenter.Y() < enemyCenter.Y() && overlap.X() > overlap.Y() {
		slog.Info("player stomped on enemy", slog.String("playerName", player.GetName()), slog.String("enemyName", enemy.GetName()))
		enemyHealthComp := enemy.GetComponent(def.ComponentTypeHealth).(*component.HealthComponent)
		if enemyHealthComp == nil {
			slog.Error("enemy health component not found", slog.String("enemyName", enemy.GetName()))
			return
		}
		// 造成1点伤害
		enemyHealthComp.TakeDamage(1)
		if !enemyHealthComp.IsAlive() {
			slog.Info("enemy is dead", slog.String("enemyName", enemy.GetName()))
			enemy.SetNeedRemove(true)
			// 敌人死亡后，创建死亡特效
			gs.createEffect(enemyCenter, enemy.GetTag())
		}
		// 玩家跳起效果
		player.GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent).Velocity[1] = -300.0
	} else {
		// 踩踏失败，玩家受伤
		slog.Info("player failed to stomp on enemy", slog.String("playerName", player.GetName()), slog.String("enemyName", enemy.GetName()))
		player.GetComponent(def.ComponentTypePlayer).(*gcomponent.PlayerComponent).TakeDamage(1)
		// TODO 其他受伤逻辑
	}
}

// 玩家与道具碰撞处理
func (gs *GameScene) playerVSItemCollision(player, item *object.GameObject) {
	if item.GetName() == "fruit" {
		// 加血
		player.GetComponent(def.ComponentTypeHealth).(*component.HealthComponent).Heal(1)
	} else if item.GetName() == "gem" {
		// TODO: 加分
	}
	// 标记道具为待删除状态
	item.SetNeedRemove(true)
	itemAABB := item.GetComponent(def.ComponentTypeCollider).(*component.ColliderComponent).GetWorldAABB()
	// 创建特效
	gs.createEffect(itemAABB.Position.Add(itemAABB.Size.Mul(0.5)), item.GetTag())
}

/**
 * @brief 创建一个特效对象，一次性
 * @param centerPos 特效中心位置
 * @param tag 特效标签（决定特效类型,例如"enemy","item"）
 */
func (gs *GameScene) createEffect(centerPos mgl32.Vec2, tag string) {
	// 创建游戏对象和变换组件
	effectObj := object.NewGameObject("effect_"+tag, "effect_"+tag)
	effectTransform := component.NewTransformComponent(centerPos, mgl32.Vec2{1.0, 1.0}, 0.0)
	effectObj.AddComponent(effectTransform)

	// 根据标签创建不同的精灵组件和动画
	animation := render.NewAnimation("effect", false)
	switch tag {
	case "enemy":
		effectSprite := component.NewSpriteComponent("assets/textures/FX/enemy-deadth.png", gs.ctx.ResourceManager, utils.AlignCenter, nil, false)
		effectObj.AddComponent(effectSprite)
		for i := range 5 {
			animation.AddFrame(&sdl.FRect{X: float32(i * 40), Y: 0.0, W: 40.0, H: 40.0}, 0.1)
		}
	case "item":
		effectSprite := component.NewSpriteComponent("assets/textures/FX/item-feedback.png", gs.ctx.ResourceManager, utils.AlignCenter, nil, false)
		effectObj.AddComponent(effectSprite)
		for i := range 4 {
			animation.AddFrame(&sdl.FRect{X: float32(i * 32), Y: 0.0, W: 32.0, H: 32.0}, 0.1)
		}
	default:
		slog.Error("createEffect: unknown tag", slog.String("tag", tag))
		return
	}

	// 根据创建的动画，添加动画组件，并设置为单次播放
	animationComponent := component.NewAnimationComponent()
	effectObj.AddComponent(animationComponent)
	animationComponent.AddAnimation(animation)
	animationComponent.SetOneShotRemoval(true)
	animationComponent.PlayAnimation("effect")
	gs.SafeAddGameObject(effectObj)
}
