package math

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// float32 color
type FColor struct {
	R, G, B, A float32
}

// 数值类型
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// 2维向量的布尔值
type Vec2B [2]bool

// 获取向量的X分量
func (v Vec2B) X() bool {
	return v[0]
}

// 获取向量的Y分量
func (v Vec2B) Y() bool {
	return v[1]
}

// 矩形
type Rect struct {
	// 矩形左上角的世界坐标
	Position mgl32.Vec2
	// 矩形的大小
	Size mgl32.Vec2
}

// 限制向量在min向量和max向量之间
func Mgl32Vec2Clamp(vec, min, max mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		mgl32.Clamp(vec.X(), min.X(), max.X()),
		mgl32.Clamp(vec.Y(), min.Y(), max.Y()),
	}
}

// 向量分量乘法
func Mgl32Vec2MulElem(src, factor mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		src.X() * factor.X(),
		src.Y() * factor.Y(),
	}
}

// 取模运算
func Mod[T Number](x, y T) T {
	res := float64(x) - float64(y)*math.Floor(float64(x)/float64(y))
	return T(res)
}

// 2维向量的绝对值
func Mgl32Vec2ABS(a, b mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		mgl32.Abs(a.X() - b.X()),
		mgl32.Abs(a.Y() - b.Y()),
	}
}

// 2维向量的线性插值
func Mgl32Vec2Mix(current, target mgl32.Vec2, t float32) mgl32.Vec2 {
	// 确保 t 在 0.0 到 1.0 之间，防止过冲（可选，取决于你是否想要弹性效果）
	if t > 1.0 {
		t = 1.0
	}
	if t < 0.0 {
		t = 0.0
	}

	// 公式：A + (B - A) * t
	return current.Add(target.Sub(current).Mul(t))
}

// clamp
// 限制值在min和max之间
func Clamp[T Number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
