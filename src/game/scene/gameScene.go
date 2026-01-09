package scene

import (
	"log/slog"
	"strconv"

	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	escene "sunny_land/src/engine/scene"
	"sunny_land/src/engine/ui"
	"sunny_land/src/engine/utils"
	"sunny_land/src/engine/utils/def"
	emath "sunny_land/src/engine/utils/math"
	gcomponent "sunny_land/src/game/component"
	"sunny_land/src/game/component/ai"
	"sunny_land/src/game/data"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 游戏场景，主要的游戏场景，包含玩家、敌人、关卡元素等。
type GameScene struct {
	// 继承基础场景
	escene.Scene
	// 保存玩家对象指针，方便访问
	playerObject *object.GameObject
	// 场景共享数据
	sessionData *data.SessionData
	// 得分标签
	scoreLabel *ui.UILabel
	// 生命值面板
	healthPanel *ui.UIPanel
}

// 确保GameScene实现IScene接口
var _ escene.IScene = (*GameScene)(nil)

// 创建游戏场景
func NewGameScene(ctx *econtext.Context, sceneManager *escene.SceneManager, sd *data.SessionData) *GameScene {
	gs := &GameScene{}
	escene.BuildScene(&gs.Scene, "GameScene", ctx, sceneManager)
	gs.sessionData = sd
	if sd == nil {
		gs.sessionData = data.NewSessionData()
	}
	slog.Debug("GameScene created")
	return gs
}

// 初始化游戏场景
func (gs *GameScene) Init() {
	gs.Scene.Init()

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

	// 初始化UI
	if !gs.InitUI() {
		slog.Error("ui init failed")
		gs.GetContext().InputManager.SetShouldQuit(true)
		return
	}

	// 设置音量
	// 设置背景音乐音量为20%
	gs.GetContext().AudioPlayer.SetMusicVolume(0.2)
	// 设置音效音量为50%
	gs.GetContext().AudioPlayer.SetSoundVolume(0.5)
	// 播放背景音乐
	gs.GetContext().AudioPlayer.PlayMusic("assets/audio/hurry_up_and_run.ogg", true)

	slog.Debug("GameScene initialized", slog.String("sceneName", gs.GetName()))
}

// 初始化关卡
func (gs *GameScene) InitLevel() bool {
	// 加载关卡
	levelPath := gs.sessionData.GetMapPath()
	if !escene.NewLevelLoader().LoadLevel(levelPath, gs) {
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

	gs.GetContext().PhysicsEngine.RegisterTileLayerComponent(tileLayerComp)
	slog.Info("main layer registered to physics engine")

	// 世界大小
	worldSize := mainLayer.GetComponent(def.ComponentTypeTileLayer).(*component.TileLayerComponent).GetWorldSize()
	// 设置相机限制范围
	gs.GetContext().Camera.SetLimitBounds(&emath.Rect{Position: mgl32.Vec2{0.0, 0.0}, Size: worldSize})
	// 开始时重置相机位置，以免切换场景时晃动
	gs.GetContext().Camera.SetPosition(mgl32.Vec2{0.0, 0.0})

	// 设置世界边界
	gs.GetContext().PhysicsEngine.SetWorldBounds(&emath.Rect{Position: mgl32.Vec2{0.0, 0.0}, Size: worldSize})

	slog.Debug("GameScene level initialized", slog.String("sceneName", gs.GetName()))
	return true
}

// 根据关卡名称获取对应的地图文件路径
func (gs *GameScene) levelNameToPath(levelName string) string {
	return "assets/maps/" + levelName + ".tmj"
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
	gs.GetContext().Camera.SetTargetTC(transformComp)

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
			aiCom := component.NewAIComponent()
			if gt.AddComponent(aiCom) != nil {
				yMax := gt.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent).GetPosition().Y()
				yMin := yMax - 80.0
				aiCom.SetBehavior(ai.NewUpDownBehavior(yMin, yMax, 50.0))
			}
		case "frog":
			aiCom := component.NewAIComponent()
			if gt.AddComponent(aiCom) != nil {
				xMax := gt.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent).GetPosition().X() - 10.0
				xMin := xMax - 90.0
				aiCom.SetBehavior(ai.NewJumpBehavior(xMin, xMax, mgl32.Vec2{100.0, -300.0}, 2.0))
			}
		case "opossum":
			aiCom := component.NewAIComponent()
			if gt.AddComponent(aiCom) != nil {
				xMax := gt.GetComponent(def.ComponentTypeTransform).(*component.TransformComponent).GetPosition().X()
				xMin := xMax - 200.0
				aiCom.SetBehavior(ai.NewPatrolBehavior(xMin, xMax, 50.0))
			}
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

// 初始化UI
func (gs *GameScene) InitUI() bool {
	if !gs.UIManager.Init(mgl32.Vec2{640.0, 360.0}) {
		return false
	}

	gs.createScoreUI()
	gs.createHealthPanel()

	return true
}

// 更新
func (gs *GameScene) Update(dt float64) {
	gs.Scene.Update(dt)
	// 处理游戏对象间的碰撞事件
	gs.handleObjectCollisions()
	// 处理瓦片触发事件
	gs.handleTileTriggers()
}

// 渲染
func (gs *GameScene) Render() {
	gs.Scene.Render()
}

// 处理事件
func (gs *GameScene) HandleInput() {
	gs.Scene.HandleInput()
}

// 清理
func (gs *GameScene) Clean() {
	gs.Scene.Clean()
}

// 处理瓦片触发事件，从PhysicsEngine获取信息
func (gs *GameScene) handleTileTriggers() {
	// 从物理引擎获取触发事件
	triggerEvents := gs.GetContext().PhysicsEngine.GetTileTriggerEvents()
	for _, event := range triggerEvents {
		// 瓦片触发事件的对象
		obj := event.GameObject.(*object.GameObject)
		// 瓦片类型
		tileType := event.TileType
		if tileType == physics.TileTypeHazard {
			// 处理玩家与"hazard"对象的碰撞
			if obj.GetName() == "player" {
				gs.handlePlayerDamage(1)
			}
			// TODO: 其他对象类型的处理，目前让敌人无视瓦片伤害
		}
	}
}

// 处理玩家受伤，更新得分、UI等
func (gs *GameScene) handlePlayerDamage(damage int) {
	playerComAny := gs.playerObject.GetComponent(def.ComponentTypePlayer)
	if playerComAny == nil {
		slog.Error("player component not found")
		return
	}
	playerCom := playerComAny.(*gcomponent.PlayerComponent)
	if !playerCom.TakeDamage(damage) {
		// 没有受伤，直接返回
		return
	}
	if playerCom.IsDead() {
		slog.Info("player dead", slog.String("name", gs.playerObject.GetName()))
		// TODO: 可能的死亡逻辑处理
	}
	// 更新生命值及HealthUI
	gs.updateHealthWithUI()
}

// 处理游戏对象间的碰撞逻辑，从PhysicsEngine获取信息
func (gs *GameScene) handleObjectCollisions() {
	// 从物理引擎获取碰撞对
	collisionPairs := gs.GetContext().PhysicsEngine.GetCollisionPairs()
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
			gs.handlePlayerDamage(1)
		} else if obj1.GetTag() == "hazard" && obj2.GetTag() == "player" {
			gs.handlePlayerDamage(1)
		}
		// 处理玩家与关底触发器碰撞
		if obj1.GetName() == "player" && obj2.GetTag() == "next_level" {
			gs.toNextLevel(obj2)
		} else if obj1.GetTag() == "next_level" && obj2.GetName() == "player" {
			gs.toNextLevel(obj1)
		}
	}
}

// 进入下一个关卡
func (gs *GameScene) toNextLevel(trigger *object.GameObject) {
	sceneName := trigger.GetName()
	mapPath := gs.levelNameToPath(sceneName)
	// 设置下一个关卡信息
	gs.sessionData.SetNextLevel(mapPath)
	nextScene := NewGameScene(gs.GetContext(), gs.SceneManager, gs.sessionData)
	gs.SceneManager.RequestReplaceScene(nextScene)
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
		enemyAIComp := enemy.GetComponent(def.ComponentTypeAI).(*component.AIComponent)
		if enemyAIComp == nil {
			slog.Error("enemy AI component not found", slog.String("enemyName", enemy.GetName()))
			return
		}
		// 造成1点伤害
		enemyAIComp.TakeDamage(1)
		if !enemyAIComp.IsAlive() {
			slog.Info("enemy is dead", slog.String("enemyName", enemy.GetName()))
			enemy.SetNeedRemove(true)
			// 敌人死亡后，创建死亡特效
			gs.createEffect(enemyCenter, enemy.GetTag())
			// 加分
			gs.sessionData.AddScore(10)
		}
		// 玩家跳起效果
		player.GetComponent(def.ComponentTypePhysics).(*component.PhysicsComponent).Velocity[1] = -300.0
		// 播放音效，此音效完全可以放在玩家的音频组件中，这里示例另一种用法：直接用AudioPlayer播放，传入文件路径
		gs.GetContext().GetAudioPlayer().PlaySound("assets/audio/punch2a.mp3")
		// 加分
		gs.addScoreWithUI(10)
	} else {
		// 踩踏失败，玩家受伤
		slog.Info("player failed to stomp on enemy", slog.String("playerName", player.GetName()), slog.String("enemyName", enemy.GetName()))
		gs.handlePlayerDamage(1)
	}
}

// 玩家与道具碰撞处理
func (gs *GameScene) playerVSItemCollision(player, item *object.GameObject) {
	_ = player

	if item.GetName() == "fruit" {
		// 加血
		gs.healWithUI(1)
	} else if item.GetName() == "gem" {
		// 加分
		gs.addScoreWithUI(5)
	}
	// 标记道具为待删除状态
	item.SetNeedRemove(true)
	itemAABB := item.GetComponent(def.ComponentTypeCollider).(*component.ColliderComponent).GetWorldAABB()
	// 创建特效
	gs.createEffect(itemAABB.Position.Add(itemAABB.Size.Mul(0.5)), item.GetTag())
	// 播放音效，此音效完全可以放在玩家的音频组件中，这里示例另一种用法：直接用AudioPlayer播放，传入文件路径
	gs.GetContext().GetAudioPlayer().PlaySound("assets/audio/poka01.mp3")
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
		effectSprite := component.NewSpriteComponent("assets/textures/FX/enemy-deadth.png", gs.GetContext().ResourceManager, utils.AlignCenter, nil, false)
		effectObj.AddComponent(effectSprite)
		for i := range 5 {
			animation.AddFrame(&sdl.FRect{X: float32(i * 40), Y: 0.0, W: 40.0, H: 40.0}, 0.1)
		}
	case "item":
		effectSprite := component.NewSpriteComponent("assets/textures/FX/item-feedback.png", gs.GetContext().ResourceManager, utils.AlignCenter, nil, false)
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

// 测试保存和加载
func (gs *GameScene) testSaveAndLoad() {
	inputManager := gs.GetContext().GetInputManager()
	if inputManager.IsActionPressed("attack") {
		gs.sessionData.SaveToFile("assets/save.json")
	}
	if inputManager.IsActionPressed("pause") {
		gs.sessionData.LoadFromFile("assets/save.json")
		slog.Info("current health", slog.Int("health", gs.sessionData.GetCurrentHealth()))
		slog.Info("current score", slog.Int("score", gs.sessionData.GetCurrentScore()))
	}
}

// 测试文本渲染
func (gs *GameScene) testTextRenderer() {
	gs.GetContext().GetTextRenderer().DrawUIText("Hello, World!", "assets/fonts/VonwaonBitmap-16px.ttf", 32,
		mgl32.Vec2{100.0, 100.0}, emath.FColor{R: 1.0, G: 0.0, B: 0.0, A: 1.0})
	gs.GetContext().GetTextRenderer().DrawText(gs.GetContext().GetCamera().(*render.Camera), "Map Text", "assets/fonts/VonwaonBitmap-16px.ttf", 32,
		mgl32.Vec2{200.0, 200.0}, emath.FColor{R: 1.0, G: 1.0, B: 1.0, A: 1.0})
}

// 创建得分UI
func (gs *GameScene) createScoreUI() {
	// 创建得分标签
	scoreText := "Score: " + strconv.Itoa(gs.sessionData.GetCurrentScore())
	gs.scoreLabel = ui.NewUILabel(gs.GetContext().TextRenderer, scoreText, "assets/fonts/VonwaonBitmap-16px.ttf", 16,
		emath.FColor{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, mgl32.Vec2{0.0, 0.0})
	// 获取根UI面板尺寸
	screenSize := gs.UIManager.GetRootElement().GetSize()
	gs.scoreLabel.SetPosition(mgl32.Vec2{screenSize.X() - 100.0, 10.0})
	gs.UIManager.AddElement(gs.scoreLabel)
}

// 创建生命值UI (或最大生命值改变时重设)
func (gs *GameScene) createHealthPanel() {
	maxHealth := gs.sessionData.GetMaxHealth()
	curHealth := gs.sessionData.GetCurrentHealth()
	stratX := float32(10.0)
	startY := float32(10.0)
	iconWidth := float32(20.0)
	iconHeight := float32(18.0)
	spacing := float32(5.0)
	fullHeartTex := "assets/textures/UI/Heart.png"
	emptyHeartTex := "assets/textures/UI/Heart-bg.png"

	// 创建一个默认的UIPanel (不需要背景色，因此大小无所谓，只用于定位)
	gs.healthPanel = ui.NewUIPanel(mgl32.Vec2{0.0, 0.0}, mgl32.Vec2{0.0, 0.0}, &emath.FColor{R: 0.0, G: 0.0, B: 0.0, A: 1.0})

	// 根据最大生命值，循环创建生命值图标(添加到UIPanel中)
	// 创建背景图标
	for i := 0; i < maxHealth; i++ {
		iconPos := mgl32.Vec2{stratX + float32(i)*(iconWidth+spacing), startY}
		iconSize := mgl32.Vec2{iconWidth, iconHeight}
		bgIcon := ui.NewUIImage(emptyHeartTex, iconPos, iconSize, nil, false)
		gs.healthPanel.AddChild(bgIcon)
	}
	// 创建前景图标
	for i := 0; i < curHealth; i++ {
		iconPos := mgl32.Vec2{stratX + float32(i)*(iconWidth+spacing), startY}
		iconSize := mgl32.Vec2{iconWidth, iconHeight}
		fgIcon := ui.NewUIImage(fullHeartTex, iconPos, iconSize, nil, false)
		gs.healthPanel.AddChild(fgIcon)
	}
	// 添加到根UI面板
	gs.UIManager.AddElement(gs.healthPanel)
}

// 增加得分，同时更新UI
func (gs *GameScene) addScoreWithUI(score int) {
	gs.sessionData.AddScore(score)
	gs.scoreLabel.SetText("Score: " + strconv.Itoa(gs.sessionData.GetCurrentScore()))
	slog.Info("add score", slog.Int("score", score))
}

// 增加生命，同时更新UI
func (gs *GameScene) healWithUI(health int) {
	gs.playerObject.GetComponent(def.ComponentTypeHealth).(*component.HealthComponent).Heal(health)
	// 更新生命值与UI
	gs.updateHealthWithUI()
}

// 更新生命值UI(只适用最大生命值不变的情况)
func (gs *GameScene) updateHealthWithUI() {
	if gs.playerObject == nil || gs.healthPanel == nil {
		slog.Error(" playerObject or healthPanel is nil")
		return
	}

	// 获取当前生命值并更新游戏数据
	curHealth := gs.playerObject.GetComponent(def.ComponentTypeHealth).(*component.HealthComponent).GetCurHealth()
	gs.sessionData.SetCurrentHealth(curHealth)
	maxHealth := gs.sessionData.GetMaxHealth()

	// 前景图标是后添加的，因此设置后半段的可见性即可
	for i := maxHealth; i < maxHealth*2; i++ {
		childs := gs.healthPanel.GetChildren()
		j := 0
		for child := childs.Front(); child != nil; child = child.Next() {
			if j == i {
				child.Value.(*ui.UIImage).SetVisible(i-maxHealth < curHealth)
				goto NEXT
			}
			j++
		}
	NEXT:
	}
}
