package physics

import (
	"sunny_land/src/engine/utils/math"
	emath "sunny_land/src/engine/utils/math"

	"github.com/go-gl/mathgl/mgl32"
)

// 碰撞器组件抽象
type IColliderComponent interface {
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

// 检查两个碰撞器组件是否碰撞
func checkCollision(a, b IColliderComponent) bool {
	// 获取碰撞器
	aCollider := a.GetCollider()
	bCollider := b.GetCollider()
	aTransform := a.GetTransformComponent()
	bTransform := b.GetTransformComponent()

	// 先计算最小包围盒是否碰撞，如果没有碰撞，直接返回false
	aSize := emath.Mgl32Vec2MulElem(aCollider.GetAABBSize(), aTransform.GetScale())
	bSize := emath.Mgl32Vec2MulElem(bCollider.GetAABBSize(), bTransform.GetScale())
	aPos := aTransform.GetPosition().Add(a.GetOffset())
	bPos := bTransform.GetPosition().Add(b.GetOffset())
	if !checkAABBOverlap(aPos, aSize, bPos, bSize) {
		return false
	}

	// 最小碰撞盒有碰撞，再进行更加细致的判断
	// AABB vs AABB, 直接返回true
	if aCollider.GetType() == ColliderTypeAABB && bCollider.GetType() == ColliderTypeAABB {
		return true
	}

	// circle vs circle，判断两个圆心之间的距离是否小于半径之和
	if aCollider.GetType() == ColliderTypeCircle && bCollider.GetType() == ColliderTypeCircle {
		aCenter := aPos.Add(aSize.Mul(0.5))
		bCenter := bPos.Add(bSize.Mul(0.5))
		aRadius := aSize.X() * 0.5
		bRadius := bSize.X() * 0.5
		return checkCircleOverlap(aCenter, aRadius, bCenter, bRadius)
	}

	// AABB vs circle，判断圆心到AABB的最邻近是否在圆内
	if aCollider.GetType() == ColliderTypeAABB && bCollider.GetType() == ColliderTypeCircle {
		// 圆心位置
		bCenter := bPos.Add(bSize.Mul(0.5))
		// 半径
		bRadius := bSize.X() * 0.5
		// 计算圆心到AABB的最邻近点，如果最邻近点到圆心的距离 小于 半径，说明撞上了，如果最邻近点到圆心的距离 大于 半径，说明没撞上。
		nearestPoint := math.Mgl32Vec2Clamp(bCenter, aPos, aPos.Add(aSize))
		return checkPointInCircle(nearestPoint, bCenter, bRadius)
	}

	// circle vs AABB，判断圆心到AABB的最邻近是否在圆内
	if aCollider.GetType() == ColliderTypeCircle && bCollider.GetType() == ColliderTypeAABB {
		// 圆心位置
		aCenter := aPos.Add(aSize.Mul(0.5))
		// 半径
		aRadius := aSize.X() * 0.5
		// 计算圆心到AABB的最邻近点，
		nearestPoint := math.Mgl32Vec2Clamp(aCenter, bPos, bPos.Add(bSize))
		return checkPointInCircle(nearestPoint, aCenter, aRadius)
	}

	return false
}

// 检查点是否在圆内
func checkPointInCircle(point, center mgl32.Vec2, radius float32) bool {
	return point.Sub(center).Len() <= radius
}

// 检查两个圆是否重叠
func checkCircleOverlap(aCenter mgl32.Vec2, aRadius float32, bCenter mgl32.Vec2, bRadius float32) bool {
	// 圆心之间的距离是否小于半径之和
	return aCenter.Sub(bCenter).Len() <= aRadius+bRadius
}

// 检查两个AABB是否重叠
func checkAABBOverlap(aPos, aSize, bPos, bSize mgl32.Vec2) bool {
	if aPos.X()+aSize.X() <= bPos.X() || aPos.X() >= bPos.X()+bSize.X() ||
		aPos.Y()+aSize.Y() <= bPos.Y() || aPos.Y() >= bPos.Y()+bSize.Y() {
		return false
	}
	return true
}
