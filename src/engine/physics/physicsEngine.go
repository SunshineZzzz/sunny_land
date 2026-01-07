package physics

import (
	"log/slog"
	"math"

	"sunny_land/src/engine/input"
	"sunny_land/src/engine/utils/def"
	emath "sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 渲染器抽象
type IRenderer interface {
	// 绘制视差精灵图
	DrawSpriteWithParallax(ICamera, ISprite, mgl32.Vec2, mgl32.Vec2, mgl32.Vec2, emath.Vec2B)
	// 绘制精灵图
	DrawSprite(ICamera, ISprite, mgl32.Vec2, mgl32.Vec2, float64)
}

// 摄像机抽象
type ICamera interface {
	// 世界坐标转换为屏幕坐标(视口坐标)，考虑视差
	// scrollFactor，视差系数，用于计算视差效果，0.0表示没有视差(固定背景)，1.0表示背景跟着相机移动，0.0~1.0之间视差
	// 移动得越快，看起来离玩家越近, 视差系数越接近1.0。移动得越慢，看起来离玩家越远，视差系数越接近0.0
	WorldToScreenWithParallax(mgl32.Vec2, mgl32.Vec2) mgl32.Vec2
	// 获取相机视口大小(屏幕大小)
	GetViewportSize() mgl32.Vec2
	// 世界坐标转换为屏幕坐标(视口坐标)
	WorldToScreen(mgl32.Vec2) mgl32.Vec2
}

// 精灵图抽象
type ISprite interface {
	// 获取要绘制的纹理部分
	GetSourceRect() *sdl.FRect
	// 获取纹理ID
	GetTextureId() string
	// 获取是否水平反转
	GetIsFlipped() bool
}

// 代表动画中的单个帧
type AnimationFrame struct {
	// 纹理图集上此帧的区域
	SourceRect *sdl.FRect
	// 此帧的显示时间(秒)
	Duration float64
}

// 动画抽象
type IAnimation interface {
	// 检查动画是否没有帧
	IsEmpty() bool
	// 获取在给定时间点应该显示的动画帧
	GetFrameAtTime(float64) *AnimationFrame
	// 是否循环播放
	IsLooping() bool
	// 获取动画总持续时间(秒)
	GetTotalDuration() float64
	// 获取动画名称
	GetName() string
}

// 上下文抽象
type IContext interface {
	// 获取渲染器
	GetRenderer() IRenderer
	// 获取摄像机
	GetCamera() ICamera
	// 获取输入管理器
	GetInputManager() input.InputManager
}

// 游戏对象抽象
type IGameObject interface {
	// 获取组件
	GetComponent(def.ComponentType) IComponent
	// 获取名字
	GetName() string
	// 获取标签
	GetTag() string
	// 设置是否需要删除
	SetNeedRemove(bool)
}

// 组件抽象
type IComponent interface {
	// 获取类型
	GetType() def.ComponentType
	// 初始化组件
	Init()
	// 更新组件
	Update(float64, IContext)
	// 处理输入
	HandleInput(IContext)
	// 渲染
	Render(IContext)
	// 清理
	Clean()
	// 设置组件所属的游戏对象
	SetOwner(IGameObject)
	// 获取组件所属的游戏对象
	GetOwner() IGameObject
}

// 碰撞器组件抽象
type IColliderComponent interface {
	// 继承组件接口
	IComponent
	// 获取碰撞器
	GetCollider() ICollider
	// 获取变换组件
	GetTransformComponent() ITransformComponent
	// 获取偏移量
	GetOffset() mgl32.Vec2
	// 是否激活
	IsActive() bool
	// 是否触发
	IsTrigger() bool
	// 获取世界AABB
	GetWorldAABB() emath.Rect
}

// 变换组件抽象
type ITransformComponent interface {
	// 继承组件接口
	IComponent
	// 平移
	Translate(mgl32.Vec2)
	// 获取缩放
	GetScale() mgl32.Vec2
	// 获取位置
	GetPosition() mgl32.Vec2
	// 设置位置
	SetPosition(mgl32.Vec2)
}

// 物理组件抽象
type IPhysicsComponent interface {
	// 继承组件接口
	IComponent
	// 组件是否启用
	IsEnabled() bool
	// 组件是否受重力影响
	IsUseGravity() bool
	// 获取质量
	GetMass() float32
	// 添加力
	AddForce(mgl32.Vec2)
	// 获取力
	GetForce() mgl32.Vec2
	// 清除力
	ClearForce()
	// 获取变换组件
	GetTransformComponent() ITransformComponent
	// 获取速度
	GetVelocity() mgl32.Vec2
	// 设置速度
	SetVelocity(mgl32.Vec2)
	// 重置所有碰撞标志位
	ResetCollisionFlags()
	// 设置下方碰撞标志位
	SetCollidedBelow(bool)
	// 设置上方碰撞标志位
	SetCollidedAbove(bool)
	// 设置左侧碰撞标志位
	SetCollidedLeft(bool)
	// 设置梯子顶层碰撞标志位
	SetCollidedLadderTop(bool)
	// 设置右侧碰撞标志位
	SetCollidedRight(bool)
	// 设置梯子碰撞标志位
	SetCollidedLadder(bool)
	// 检查是否与底部碰撞
	HasCollidedBelow() bool
	// 检查是否与顶部碰撞
	HasCollidedAbove() bool
	// 检查是否与左侧碰撞
	HasCollidedLeft() bool
	// 检查是否与右侧碰撞
	HasCollidedRight() bool
	// 检查是否与梯子碰撞
	HasCollidedLadder() bool
	// 检查是否与梯子顶层碰撞
	HasCollidedLadderTop() bool
}

// 健康组件抽象
type IHealthComponent interface {
	// 继承组件接口
	IComponent
}

// 瓦片类型
type TileType int

const (
	// 空白瓦片
	TileTypeEmpty TileType = iota
	// 普通瓦片
	TileTypeNormal
	// 静止可碰撞瓦片
	TileTypeSolid
	// 单向静止可碰撞瓦片
	TileTypeUniSolid
	// 斜坡瓦片，高度:左0 右1
	TileTypeSlope_0_1
	// 斜坡瓦片，高度:左1 右0
	TileTypeSlope_1_0
	// 斜坡瓦片，高度:左0 右1/2
	TileTypeSlope_0_2
	// 斜坡瓦片，高度:左1/2 右1
	TileTypeSlope_2_1
	// 斜坡瓦片，高度:左1 右1/2
	TileTypeSlope_1_2
	// 斜坡瓦片，高度:左1/2 右0
	TileTypeSlope_2_0
	// 危险瓦片，例如火焰、尖刺等
	TileTypeHazard
	// 梯子瓦片
	TileTypeLadder
)

// 单个瓦片信息
type TileInfo struct {
	// 精灵图
	Sprite ISprite
	// 瓦片类型
	Type TileType
}

// 瓦片图层组件抽象
type ITileLayerComponent interface {
	// 获取瓦片大小
	GetTileSize() mgl32.Vec2
	// 获取指定位置的瓦片类型，pos不是整数坐标
	GetTileTypeAt(int, int) TileType
}

// 碰撞组件对
type CollisionPair struct {
	A IGameObject
	B IGameObject
}

// 瓦片触发事件对
type TileTriggerEventPair struct {
	// 触发事件的游戏对象
	GameObject IGameObject
	// 触发的瓦片类型
	TileType TileType
}

// 物理引擎，负责管理和模拟物理行为，碰撞检测
type PhysicsEngine struct {
	// 注册的物理组件容器
	physicsComponents []IPhysicsComponent
	// 注册的瓦片图层组件容器
	tileLayerComponents []ITileLayerComponent
	// 默认重力加速度{0.0, 980.0}，单位：像素每二次方秒，现实中是，9.8米/s^2，游戏中是，100像素 * 9.8米/s^2 = 980.0像素/s^2
	gravity mgl32.Vec2
	// 最大速度值{-500.0, -500.0}/{500.0, 500.0}，单位：像素/秒
	maxSpeed float32
	// 存储本帧发生的碰撞组件对
	collisionPairs []CollisionPair
	// 存储本帧发生的瓦片触发事件
	tileTriggerEvents []TileTriggerEventPair
	// 世界边界，限制物体移动
	worldBounds *emath.Rect
}

// 创建物理引擎
func NewPhysicsEngine() *PhysicsEngine {
	slog.Debug("new physics engine")
	return &PhysicsEngine{
		physicsComponents: make([]IPhysicsComponent, 0),
		gravity:           mgl32.Vec2{0.0, 980.0},
		maxSpeed:          500.0,
		collisionPairs:    make([]CollisionPair, 0),
		tileTriggerEvents: make([]TileTriggerEventPair, 0),
	}
}

// 注册物理组件
func (pe *PhysicsEngine) RegisterPhysicsComponent(component IPhysicsComponent) {
	slog.Debug("register physics component")
	pe.physicsComponents = append(pe.physicsComponents, component)
}

// 移除注册物理组件
func (pe *PhysicsEngine) UnregisterComponent(component IPhysicsComponent) {
	slog.Debug("remove physics component")
	for i, comp := range pe.physicsComponents {
		if comp == component {
			pe.physicsComponents = append(pe.physicsComponents[:i], pe.physicsComponents[i+1:]...)
			return
		}
	}
}

// 注册瓦片图层组件
func (pe *PhysicsEngine) RegisterTileLayerComponent(component ITileLayerComponent) {
	slog.Debug("register tile layer component")
	pe.tileLayerComponents = append(pe.tileLayerComponents, component)
}

// 移除注册瓦片图层组件
func (pe *PhysicsEngine) UnregisterTileLayerComponent(component ITileLayerComponent) {
	slog.Debug("remove tile layer component")
	for i, comp := range pe.tileLayerComponents {
		if comp == component {
			pe.tileLayerComponents = append(pe.tileLayerComponents[:i], pe.tileLayerComponents[i+1:]...)
			return
		}
	}
}

// 更新
func (pe *PhysicsEngine) Update(deltaTime float64) {
	// 每帧开始时先清空碰撞对列表和瓦片触发事件列表
	pe.collisionPairs = pe.collisionPairs[:0]
	pe.tileTriggerEvents = pe.tileTriggerEvents[:0]

	// 遍历所有注册的物理组件，更新他们的物理状态
	for _, pc := range pe.physicsComponents {
		if pc == nil || !pc.IsEnabled() {
			continue
		}

		// 重置碰撞标志位
		pc.ResetCollisionFlags()

		// 是否使用重力，如果组件接受重力影响，F = m * a
		if pc.IsUseGravity() {
			pc.AddForce(pe.gravity.Mul(pc.GetMass()))
		}
		// 还可以添加其他影响，比如风力，摩擦力，目前不考虑

		// 更新速度，v += a * dt，其中 a = F / m
		pc.SetVelocity(
			pc.GetVelocity().Add(
				pc.GetForce().Mul(1.0 / pc.GetMass()).Mul(float32(deltaTime)),
			),
		)
		// 清除当前帧的力
		pc.ClearForce()

		// 处理瓦片层(图块层)碰撞
		pe.resolveTileLayerCollisions(pc, deltaTime)

		// 应用世界边界，限制物体移动范围
		pe.ApplyWorldBounds(pc)
	}

	// 处理对象间的碰撞
	pe.checkObjectCollisions()
	// 检测瓦片触发事件，检测前已经处理完位移
	pe.checkTileTriggers()
}

// 检查对象间的碰撞
func (pe *PhysicsEngine) checkObjectCollisions() {
	// 两层循环判断所有包含碰撞组件的GameObject对是否发生碰撞
	for i, pca := range pe.physicsComponents {
		// 物理组件如果都没有启用，不考虑碰撞
		if pca == nil || !pca.IsEnabled() {
			continue
		}

		// 获取碰撞组件，如果都没有启用，不考虑碰撞
		cca := pca.GetOwner().GetComponent(def.ComponentTypeCollider).(IColliderComponent)
		if cca == nil || !cca.IsActive() {
			continue
		}

		for j := i + 1; j < len(pe.physicsComponents); j++ {
			pcb := pe.physicsComponents[j]
			if pcb == nil || !pcb.IsEnabled() {
				continue
			}

			ccb := pcb.GetOwner().GetComponent(def.ComponentTypeCollider).(IColliderComponent)
			if ccb == nil || !ccb.IsActive() {
				continue
			}

			// 检查碰撞
			if checkCollision(cca, ccb) {
				// TODO: 并不是所有碰撞都需要插入切片，比如触发器，未来会添加过滤条件
				// 如果是可移动物体与SOLID静态物体碰撞，直接处理位置变化，不用记录碰撞
				if pca.GetOwner().GetTag() != "solid" && pcb.GetOwner().GetTag() == "solid" {
					pe.resolveSolidObjectCollisions(pca.GetOwner(), pcb.GetOwner())
				} else if pca.GetOwner().GetTag() == "solid" && pcb.GetOwner().GetTag() != "solid" {
					pe.resolveSolidObjectCollisions(pcb.GetOwner(), pca.GetOwner())
				} else {
					// 碰撞对加入切片
					pe.collisionPairs = append(pe.collisionPairs, CollisionPair{pca.GetOwner(), pcb.GetOwner()})
				}
			}
		}
	}
}

// 获取本帧检测到的所有碰撞对，此列表在每次update开始时清空
func (pe *PhysicsEngine) GetCollisionPairs() []CollisionPair {
	return pe.collisionPairs
}

// 获取本帧检测到的所有瓦片触发事件，此列表在每次update开始时清空
func (pe *PhysicsEngine) GetTileTriggerEvents() []TileTriggerEventPair {
	return pe.tileTriggerEvents
}

// 处理瓦片层(图块层)碰撞
func (pe *PhysicsEngine) resolveTileLayerCollisions(pc IPhysicsComponent, deltaTime float64) {
	if pc.GetOwner() == nil {
		return
	}
	// 获取变换组件
	tc := pc.GetTransformComponent()
	// 获取碰撞组件
	cca := pc.GetOwner().GetComponent(def.ComponentTypeCollider).(IColliderComponent)
	if tc == nil || cca == nil || cca.IsTrigger() {
		return
	}
	// 使用最小包围盒进行碰撞检测，不考虑圆形碰撞器啥的，简化
	worldAABB := cca.GetWorldAABB()
	// 物体的当前位置
	objPos := worldAABB.Position
	objSize := worldAABB.Size
	if objSize.X() <= 0.0 || objSize.Y() <= 0.0 {
		return
	}

	// 右下瓦片y，左下瓦片y，右下瓦片x，右上瓦片x，这些情况需要减去1个像素，避免检查到下一行/列的瓦片
	tolerance := float32(1.0)
	// 计算物体在dt时间内的位移
	ds := pc.GetVelocity().Mul(float32(deltaTime))
	// 计算物体在dt时间内的目标(期望)位置
	newObjPos := objPos.Add(ds)

	// 如果碰撞器未激活，直接让物体正常移动，然后返回
	if !cca.IsActive() {
		tc.Translate(ds)
		pc.SetVelocity(
			emath.Mgl32Vec2Clamp(
				pc.GetVelocity(),
				mgl32.Vec2{-pe.maxSpeed, -pe.maxSpeed},
				mgl32.Vec2{pe.maxSpeed, pe.maxSpeed},
			),
		)
		return
	}

	// 遍历所有注册的碰撞瓦片层
	for _, tl := range pe.tileLayerComponents {
		if tl == nil {
			continue
		}

		// 获取瓦片大小
		tileSize := tl.GetTileSize()
		// 采用轴分离碰撞检测，如果不这样做就会出现问题，比如：
		// 我想往右走1像素，同时往上走1像素。计算目标位置：现在的坐标(x, y)变成(x+1, y+1)。刚好(x+1, y+1)有一堵墙(碰撞)。这次移动是非法的，程序把你的坐标锁定在原地(x, y)
		// 玩家的感受：我按着“右”和“上”，角色动都不动，像被胶水粘在了墙上。所以需要分离。处理X轴(向右走1像素)，发现(x+1, y)确实撞墙了。取消这次X轴位移，保持x不变。
		// 处理Y轴(向上走1像素)，程序计算新位置(x, y+1)，发现上方是空的，没墙，成功移动，坐标更新为(x, y+1)。和顺序无关，甚至可以并行判断。
		//
		//  x轴，x坐标为期望位置，y坐标为当前位置，检测x方向是否碰撞，如果碰撞，取消X轴位移和速度
		if ds.X() > 0.0 {
			// 检测右侧碰撞
			// 右上x期望坐标，需要检测右上瓦片和右下瓦片是否碰撞
			rightTopX := newObjPos.X() + objSize.X()
			// 右上x期望坐标对应的右上x瓦片坐标，这里不要减去tolerance，实际计算是哪个就需要判断那个瓦片x位置
			rightTopTileX := int(math.Floor(float64(rightTopX / tileSize.X())))
			// 右上瓦片y，这个也是不要减tolerance，实际计算是哪个就需要判断那个瓦片y位置
			rightTopTileY := int(math.Floor(float64(objPos.Y() / tileSize.Y())))
			// 获取右上瓦片类型
			rightTopTileType := tl.GetTileTypeAt(rightTopTileX, rightTopTileY)
			// 右下瓦片y，这个需要减tolerance，否则会检查到下一行的瓦片
			rightBottomtileY := int(math.Floor(float64((objPos.Y() + objSize.Y() - tolerance) / tileSize.Y())))
			// 获取右下瓦片类型
			rightBottomTileType := tl.GetTileTypeAt(rightTopTileX, rightBottomtileY)

			if rightTopTileType == TileTypeSolid || rightBottomTileType == TileTypeSolid {
				// 撞墙了，速度归0
				pc.SetVelocity(mgl32.Vec2{0.0, pc.GetVelocity().Y()})
				// x方向移动到贴着墙壁的位置
				newObjPos[0] = float32(rightTopTileX)*tileSize.X() - objSize.X()
				pc.SetCollidedRight(true)
			} else {
				// 没有碰撞，需要检测右下角是否是斜坡瓦片
				// 计算斜坡碰撞x的偏移量，可以利用相似三角形计算y方向的偏移量
				widthRight := newObjPos.X() + objSize.X() - float32(rightTopTileX)*tileSize.X()
				heightRight := pe.getTileHeightAtWidth(widthRight, rightBottomTileType, tileSize)
				if heightRight > 0.0 {
					// 右下瓦片下标*瓦片高度-物体高度-斜坡高度=物体在斜坡上的y坐标
					// 假设没有斜坡，物体应该在的y坐标，再减去heightRight，就是物体在斜坡上的y坐标
					targetY := float32(rightBottomtileY+1)*tileSize.Y() - objSize.Y() - heightRight
					if targetY < newObjPos.Y() {
						// 说明碰撞了
						newObjPos[1] = targetY
						pc.SetCollidedBelow(true)
					}
				}
			}
		} else if ds.X() < 0.0 {
			// 检测左侧碰撞
			// 左上x期望坐标，需要检测左上瓦片和左下瓦片是否碰撞
			leftTopX := newObjPos.X()
			// 左上x期望坐标对应的左上x瓦片坐标，这里不要减去tolerance，实际计算是哪个就需要判断那个瓦片x位置
			leftTopTileX := int(math.Floor(float64(leftTopX / tileSize.X())))
			// 左上瓦片y，这个也是不要减tolerance，实际计算是哪个就需要判断那个瓦片y位置
			leftTopTileY := int(math.Floor(float64(objPos.Y() / tileSize.Y())))
			// 获取左上瓦片类型
			leftTopTileType := tl.GetTileTypeAt(leftTopTileX, leftTopTileY)
			// 左下瓦片y，这个需要减tolerance，否则会检查到下一行的瓦片
			leftBottomtileY := int(math.Floor(float64((objPos.Y() + objSize.Y() - tolerance) / tileSize.Y())))
			// 获取左下瓦片类型
			leftBottomTileType := tl.GetTileTypeAt(leftTopTileX, leftBottomtileY)

			if leftTopTileType == TileTypeSolid || leftBottomTileType == TileTypeSolid {
				// 撞墙了，速度归0
				pc.SetVelocity(mgl32.Vec2{0.0, pc.GetVelocity().Y()})
				// x方向移动到贴着墙壁的位置
				newObjPos[0] = float32(leftTopTileX+1) * tileSize.X()
				pc.SetCollidedLeft(true)
			} else {
				// 没有碰撞，需要检测左下角是否是斜坡瓦片
				// 计算斜坡碰撞x的偏移量，可以利用相似三角形计算y方向的偏移量
				widthLeft := newObjPos.X() - float32(leftTopTileX)*tileSize.X()
				heightLeft := pe.getTileHeightAtWidth(widthLeft, leftBottomTileType, tileSize)
				if heightLeft > 0.0 {
					// 左下瓦片下标*瓦片高度-物体高度-斜坡高度=物体在斜坡上的y坐标
					// 假设没有斜坡，物体应该在的y坐标，再减去heightLeft，就是物体在斜坡上的y坐标
					targetY := float32(leftBottomtileY+1)*tileSize.Y() - objSize.Y() - heightLeft
					if targetY < newObjPos.Y() {
						// 说明碰撞了
						newObjPos[1] = targetY
						pc.SetCollidedBelow(true)
					}
				}
			}
		}

		//  y轴，y坐标为期望位置，x坐标为当前位置，检测y方向是否碰撞，如果碰撞，取消Y轴位移和速度
		if ds.Y() > 0.0 {
			// 检测底部碰撞
			// 左下y期望坐标，需要检测左下瓦片和右下瓦片是否碰撞
			leftBottomY := newObjPos.Y() + objSize.Y()
			// 左下y期望坐标对应的左下瓦片坐标，这里不要减去tolerance，实际计算是哪个就需要判断那个瓦片y位置
			leftBottomTileY := int(math.Floor(float64(leftBottomY / tileSize.Y())))
			// 左下瓦片x，这个也是不要减tolerance，实际计算是哪个就需要判断那个瓦片x位置
			leftBottomTileX := int(math.Floor(float64(objPos.X() / tileSize.X())))
			// 获取左下瓦片类型
			leftBottomTileType := tl.GetTileTypeAt(leftBottomTileX, leftBottomTileY)
			// 右下瓦片x，这个需要减tolerance，否则会检查到下一列的瓦片
			rightBottomTileX := int(math.Floor(float64((objPos.X() + objSize.X() - tolerance) / tileSize.X())))
			// 获取右下瓦片类型
			rightBottomTileType := tl.GetTileTypeAt(rightBottomTileX, leftBottomTileY)

			if leftBottomTileType == TileTypeSolid || rightBottomTileType == TileTypeSolid ||
				leftBottomTileType == TileTypeUniSolid || rightBottomTileType == TileTypeUniSolid {
				// 撞墙了，速度归0
				pc.SetVelocity(mgl32.Vec2{pc.GetVelocity().X(), 0.0})
				// y方向移动到贴着墙壁的位置
				newObjPos[1] = float32(leftBottomTileY)*tileSize.Y() - objSize.Y()
				pc.SetCollidedBelow(true)
			} else if leftBottomTileType == TileTypeLadder && rightBottomTileType == TileTypeLadder {
				// 如果两个角点都位于梯子上，则判断是不是处在梯子顶层
				// 左角点上方瓦片类型
				leftBottomTileTypeUp := tl.GetTileTypeAt(leftBottomTileX, leftBottomTileY-1)
				// 右角点上方瓦片类型
				rightBottomTileTypeUp := tl.GetTileTypeAt(rightBottomTileX, leftBottomTileY-1)
				// 如果上方不是梯子，证明处在梯子顶层
				if leftBottomTileTypeUp != TileTypeLadder && rightBottomTileTypeUp != TileTypeLadder {
					// 通过是否使用重力来区分是否处于攀爬状态
					if pc.IsUseGravity() {
						// 非攀爬状态，证明处在梯子顶层
						pc.SetCollidedLadderTop(true)
						// 重置梯子顶层碰撞标志位
						pc.SetCollidedBelow(true)
						// 让物体贴着梯子顶层位置(与SOLID情况相同)
						// 撞墙了，速度归0
						pc.SetVelocity(mgl32.Vec2{pc.GetVelocity().X(), 0.0})
						// y方向移动到贴着墙壁的位置
						newObjPos[1] = float32(leftBottomTileY)*tileSize.Y() - objSize.Y()
					} else {
						// 攀爬状态，不做任何处理
					}
				}
			} else {
				// 下方两个斜坡都需要检测
				widthLeft := newObjPos.X() - float32(leftBottomTileX)*tileSize.X()
				widthRight := newObjPos.X() + objSize.X() - float32(rightBottomTileX)*tileSize.X()
				heightLeft := pe.getTileHeightAtWidth(widthLeft, leftBottomTileType, tileSize)
				heightRight := pe.getTileHeightAtWidth(widthRight, rightBottomTileType, tileSize)
				// 找到两个角点的最高点进行检测
				height := max(heightLeft, heightRight)
				if height > 0.0 {
					// 左下瓦片下标*瓦片高度-物体高度-斜坡高度=物体在斜坡上的y坐标
					// 假设没有斜坡，物体应该在的y坐标，再减去height，就是物体在斜坡上的y坐标
					targetY := float32(leftBottomTileY+1)*tileSize.Y() - objSize.Y() - height
					if targetY < newObjPos.Y() {
						// 说明碰撞了
						newObjPos[1] = targetY
						// 只有向下运动时才需要让y速度归零
						pc.SetVelocity(mgl32.Vec2{pc.GetVelocity().X(), 0.0})
						pc.SetCollidedBelow(true)
					}
				}
			}
		} else if ds.Y() < 0.0 {
			// 检测顶部碰撞
			// 左上y期望坐标，需要检测左上瓦片和右上瓦片是否碰撞
			topTopY := newObjPos.Y()
			// 左上y期望坐标对应的左上瓦片坐标，这里不要减去tolerance，实际计算是哪个就需要判断那个瓦片y位置
			topTopTileY := int(math.Floor(float64(topTopY / tileSize.Y())))
			// 左上瓦片x，这个也是不要减tolerance，实际计算是哪个就需要判断那个瓦片x位置
			topTopTileX := int(math.Floor(float64(objPos.X() / tileSize.X())))
			// 获取左上瓦片类型
			topTopTileType := tl.GetTileTypeAt(topTopTileX, topTopTileY)
			// 右上瓦片x，这个需要减tolerance，否则会检查到下一列的瓦片
			rightTopTileX := int(math.Floor(float64((objPos.X() + objSize.X() - tolerance) / tileSize.X())))
			// 获取右上瓦片类型
			rightTopTileType := tl.GetTileTypeAt(rightTopTileX, topTopTileY)

			if topTopTileType == TileTypeSolid || rightTopTileType == TileTypeSolid {
				// 撞墙了，速度归0
				pc.SetVelocity(mgl32.Vec2{pc.GetVelocity().X(), 0.0})
				// y方向移动到贴着墙壁的位置
				newObjPos[1] = float32(topTopTileY+1) * tileSize.Y()
				pc.SetCollidedAbove(true)
			}
		}

		// 更新位置
		// tc.SetPosition(newObjPos)
		// 不可以使用SetPosition，因为有的物体碰撞盒是有偏移量的，使用SetPosition会导致碰撞盒偏移量失效
		tc.Translate(newObjPos.Sub(objPos))
		// 限制最大速度
		pc.SetVelocity(
			emath.Mgl32Vec2Clamp(
				pc.GetVelocity(),
				mgl32.Vec2{-pe.maxSpeed, -pe.maxSpeed},
				mgl32.Vec2{pe.maxSpeed, pe.maxSpeed},
			),
		)
	}
}

// 处理移动物理组件(游戏对象)和固体(静态)物理组件(游戏对象)的碰撞
func (pe *PhysicsEngine) resolveSolidObjectCollisions(moveObj, solidObj IGameObject) {
	// 进入这个函数前，已经检查了各个组件的有效性，因此直接计算
	moveTC := moveObj.GetComponent(def.ComponentTypeTransform).(ITransformComponent)
	moveCC := moveObj.GetComponent(def.ComponentTypeCollider).(IColliderComponent)
	movePC := moveObj.GetComponent(def.ComponentTypePhysics).(IPhysicsComponent)
	solidCC := solidObj.GetComponent(def.ComponentTypeCollider).(IColliderComponent)
	// 这里只能获取期望位置，因为这个函数前已经处理过了物理组件和瓦片碰撞，无法获取当前帧初始位置，因此无法进行轴分离检测
	// 未来可以重构，这里使用长宽最小平移向量解决碰撞
	moveAABB := moveCC.GetWorldAABB()
	solidAABB := solidCC.GetWorldAABB()

	// 计算移动物理组件中心位置
	moveCenter := moveAABB.Position.Add(moveAABB.Size.Mul(0.5))
	solidCenter := solidAABB.Position.Add(solidAABB.Size.Mul(0.5))
	// 计算两个包围盒的重叠部分
	overlap := moveAABB.Size.Mul(0.5).Add(solidAABB.Size.Mul(0.5)).Sub(emath.Mgl32Vec2ABS(moveCenter, solidCenter))
	// 如果重叠部分太小了，就认为没有碰撞
	if overlap.X() < 0.1 && overlap.Y() < 0.1 {
		return
	}

	if overlap.X() < overlap.Y() {
		// 如果重叠部分在x方向上更小，则认为碰撞发生在x方向上，推出x方向平移向量最小
		if moveCenter.X() < solidCenter.X() {
			// 移动物体在固体物体的左边，让移动物体贴着固体物体的右边，y方向正常移动
			moveTC.Translate(mgl32.Vec2{-overlap.X(), 0.0})
			// 如果速度为正(向右移动)，则速度归0，万一物体在固体物体的左边，速度为负，也归零就会被吸附
			if movePC.GetVelocity().X() > 0.0 {
				movePC.SetVelocity(mgl32.Vec2{0.0, movePC.GetVelocity().Y()})
				movePC.SetCollidedRight(true)
			}
		} else {
			// 移动物体在固体物体的右边，让移动物体体贴着固体物体的左边，y方向正常移动
			moveTC.Translate(mgl32.Vec2{overlap.X(), 0.0})
			// 如果速度为负(向左移动)，则速度归0，万一物体在固体物体的右边，速度为正，也归零就会被吸附
			if movePC.GetVelocity().X() < 0.0 {
				movePC.SetVelocity(mgl32.Vec2{0.0, movePC.GetVelocity().Y()})
				movePC.SetCollidedLeft(true)
			}
		}
	} else {
		// y轴方向碰撞
		if moveCenter.Y() < solidCenter.Y() {
			// 移动物体在固体物体的上面，让移动物体体贴着固体物体的下面，x方向正常移动
			moveTC.Translate(mgl32.Vec2{0.0, -overlap.Y()})
			// 如果速度为正(向上移动)，则速度归0，万一物体在固体物体的上面，速度为负，也归零就会被吸附
			if movePC.GetVelocity().Y() > 0.0 {
				movePC.SetVelocity(mgl32.Vec2{movePC.GetVelocity().X(), 0.0})
				movePC.SetCollidedBelow(true)
			}
		} else {
			// 移动物体在固体物体的下面，让移动物体体贴着固体物体的上面，x方向正常移动
			moveTC.Translate(mgl32.Vec2{0.0, overlap.Y()})
			// 如果速度为负(向下移动)，则速度归0，万一物体在固体物体的下面，速度为正，也归零就会被吸附
			if movePC.GetVelocity().Y() < 0.0 {
				movePC.SetVelocity(mgl32.Vec2{movePC.GetVelocity().X(), 0.0})
				movePC.SetCollidedAbove(true)
			}
		}
	}
}

// 根据宽度获取斜坡瓦片的高度
func (pe *PhysicsEngine) getTileHeightAtWidth(width float32, tileType TileType, tileSize mgl32.Vec2) float32 {
	relX := mgl32.Clamp(width/tileSize.X(), 0.0, 1.0)
	switch tileType {
	case TileTypeSlope_0_1:
		return relX * tileSize.Y()
	case TileTypeSlope_0_2:
		return relX * tileSize.Y() * 0.5
	case TileTypeSlope_2_1:
		return relX*tileSize.Y()*0.5 + tileSize.Y()*0.5
	case TileTypeSlope_1_0:
		return (1.0 - relX) * tileSize.Y()
	case TileTypeSlope_2_0:
		return (1.0 - relX) * tileSize.Y() * 0.5
	case TileTypeSlope_1_2:
		return (1.0-relX)*tileSize.Y()*0.5 + tileSize.Y()*0.5
	default:
		return 0.0
	}
}

// 设置世界边界
func (pe *PhysicsEngine) SetWorldBounds(bounds *emath.Rect) {
	pe.worldBounds = bounds
}

// 获取世界边界
func (pe *PhysicsEngine) GetWorldBounds() *emath.Rect {
	return pe.worldBounds
}

// 应用世界边界，限制物体移动范围
func (pe *PhysicsEngine) ApplyWorldBounds(pc IPhysicsComponent) {
	if pe.worldBounds == nil || pc == nil {
		return
	}

	// 只限定左，上，右边界，不限定下边界，以碰撞盒作为判断依据
	obj := pc.GetOwner()
	cc := obj.GetComponent(def.ComponentTypeCollider).(IColliderComponent)
	tc := obj.GetComponent(def.ComponentTypeTransform).(ITransformComponent)
	worldAABB := cc.GetWorldAABB()
	objPos := worldAABB.Position
	objSize := worldAABB.Size

	// 限制左边界
	if objPos.X() < pe.worldBounds.Position.X() {
		pc.SetVelocity(mgl32.Vec2{0.0, pc.GetVelocity().Y()})
		objPos[0] = pe.worldBounds.Position.X()
		pc.SetCollidedLeft(true)
	}
	// 限制上边界
	if objPos.Y() < pe.worldBounds.Position.Y() {
		pc.SetVelocity(mgl32.Vec2{pc.GetVelocity().X(), 0.0})
		objPos[1] = pe.worldBounds.Position.Y()
		pc.SetCollidedAbove(true)
	}
	// 限制右边界
	if objPos.X()+objSize.X() > pe.worldBounds.Position.X()+pe.worldBounds.Size.X() {
		pc.SetVelocity(mgl32.Vec2{0.0, pc.GetVelocity().Y()})
		objPos[0] = pe.worldBounds.Position.X() + pe.worldBounds.Size.X() - objSize.X()
		pc.SetCollidedRight(true)
	}
	// 更新物体位置，因为物体碰撞盒和实际大小不一定一样，所以采用平移
	tc.Translate(objPos.Sub(worldAABB.Position))
}

// 检测所有游戏对象与瓦片层的触发器类型瓦片碰撞，并记录触发事件。位移处理完毕后再调用
func (pe *PhysicsEngine) checkTileTriggers() {
	for _, pc := range pe.physicsComponents {
		if pc == nil || !pc.IsEnabled() {
			continue
		}
		obj := pc.GetOwner()
		if obj == nil {
			continue
		}
		cc := obj.GetComponent(def.ComponentTypeCollider).(IColliderComponent)
		// 如果游戏对象本就是触发器，则不需要检查瓦片触发事件
		if cc == nil || !cc.IsActive() || cc.IsTrigger() {
			continue
		}

		// 获取游戏对象的世界AABB
		worldAABB := cc.GetWorldAABB()
		// 使用set来跟踪循环遍历中已经触发过的瓦片类型，防止重复添加，例如，玩家同时踩到两个尖刺，只需要受到一次伤害
		triggeredTypes := make(map[TileType]bool)

		// 遍历所有注册的碰撞瓦片层分别进行检测
		for _, tileLayerComp := range pe.tileLayerComponents {
			if tileLayerComp == nil {
				continue
			}

			tileSize := tileLayerComp.GetTileSize()
			// 检查右边缘和下边缘时，需要减1像素，否则会检查到下一行/列的瓦片
			tolernance := float32(1.0)
			// 获取瓦片坐标范围
			startX := int(math.Floor(float64(worldAABB.Position.X() / tileSize.X())))
			endX := int(math.Ceil(float64((worldAABB.Position.X() + worldAABB.Size.X() - tolernance) / tileSize.X())))
			startY := int(math.Floor(float64(worldAABB.Position.Y() / tileSize.Y())))
			endY := int(math.Ceil(float64((worldAABB.Position.Y() + worldAABB.Size.Y() - tolernance) / tileSize.Y())))

			// 遍历瓦片坐标范围，检查是否有触发器类型的瓦片
			for x := startX; x < endX; x++ {
				for y := startY; y < endY; y++ {
					tileType := tileLayerComp.GetTileTypeAt(x, y)
					// 未来可以添加更多触发器类型的瓦片，目前只有HAZARD类型
					if tileType == TileTypeHazard && !triggeredTypes[tileType] {
						triggeredTypes[tileType] = true
					} else if tileType == TileTypeLadder && !triggeredTypes[tileType] {
						// 梯子类型不必记录到事件容器，物理引擎自己处理
						pc.SetCollidedLadder(true)
					}
				}
			}
			// 遍历触发事件集合，添加到tileTriggerEvents中
			for tileType := range triggeredTypes {
				pe.tileTriggerEvents = append(pe.tileTriggerEvents, TileTriggerEventPair{GameObject: obj, TileType: tileType})
			}
		}
	}
}
