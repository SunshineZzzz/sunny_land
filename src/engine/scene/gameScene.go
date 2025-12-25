package scene

import (
	"log/slog"
	"sunny_land/src/engine/component"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 游戏场景，主要的游戏场景，包含玩家、敌人、关卡元素等。
type GameScene struct {
	// 继承基础场景
	scene
	// 测试游戏对象
	testObject *object.GameObject
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

	// 创建测试游戏对象
	gs.createTestObject()

	slog.Debug("GameScene initialized", slog.String("sceneName", gs.sceneName))
}

// 更新
func (gs *GameScene) Update(dt float64) {
	// 测试更新摄像机
	// gs.testCamera()
	// 测试更新测试游戏对象
	gs.TestObject()
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

// 创建测试游戏对象
func (gs *GameScene) createTestObject() {
	// 创建一个测试游戏对象
	gt := object.NewGameObject("TestObject", "TestObject")
	transformComp := component.NewTransformComponent(mgl32.Vec2{100.0, 100.0}, mgl32.Vec2{1.0, 1.0}, 0.0)
	spriteComp := component.NewSpriteComponent("assets/textures/Props/big-crate.png", gs.ctx.ResourceManager, utils.AlignCenter, nil, false)
	physicsComp := component.NewPhysicsComponent(gs.ctx.PhysicsEngine, 1.0, true)
	colliderComp := component.NewColliderComponent(physics.NewAABBCollider(mgl32.Vec2{32.0, 32.0}), utils.AlignNone, mgl32.Vec2{0.0, 0.0}, false, true)
	gt.AddComponent(transformComp)
	gt.AddComponent(spriteComp)
	gt.AddComponent(physicsComp)
	gt.AddComponent(colliderComp)

	// 添加第二个游戏对象
	gt2 := object.NewGameObject("TestObject2", "TestObject2")
	transformComp2 := component.NewTransformComponent(mgl32.Vec2{50.0, 50.0}, mgl32.Vec2{1.0, 1.0}, 0.0)
	spriteComp2 := component.NewSpriteComponent("assets/textures/Props/big-crate.png", gs.ctx.ResourceManager, utils.AlignCenter, nil, false)
	physicsComp2 := component.NewPhysicsComponent(gs.ctx.PhysicsEngine, 1.0, false)
	colliderComp2 := component.NewColliderComponent(physics.NewCircleCollider(16.0), utils.AlignNone, mgl32.Vec2{0.0, 0.0}, false, true)
	gt2.AddComponent(transformComp2)
	gt2.AddComponent(spriteComp2)
	gt2.AddComponent(physicsComp2)
	gt2.AddComponent(colliderComp2)

	// 将创建好的游戏对象添加到场景中
	gs.scene.AddGameObject(gt)
	gs.scene.AddGameObject(gt2)

	// 保存测试游戏对象
	gs.testObject = gt
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
func (gs *GameScene) TestObject() {
	if gs.testObject == nil {
		return
	}
	inputManager := gs.ctx.InputManager

	if inputManager.IsActionDown("move_left") {
		gs.testObject.GetComponent(&component.TransformComponent{}).(*component.TransformComponent).Translate(mgl32.Vec2{-1, 0})
	}
	if inputManager.IsActionDown("move_right") {
		gs.testObject.GetComponent(&component.TransformComponent{}).(*component.TransformComponent).Translate(mgl32.Vec2{1, 0})
	}
	if inputManager.IsActionPressed("jump") {
		gs.testObject.GetComponent(&component.PhysicsComponent{}).(*component.PhysicsComponent).SetVelocity(mgl32.Vec2{0, -400})
	}
}

// 测试碰撞组件对
func (gs *GameScene) TestCollisionPairs() {
	collision_pairs := gs.ctx.PhysicsEngine.GetCollisionPairs()
	for _, pair := range collision_pairs {
		slog.Info("碰撞对", slog.String("a", pair.A.(*object.GameObject).GetName()), slog.String("b", pair.B.(*object.GameObject).GetName()))
	}
}
