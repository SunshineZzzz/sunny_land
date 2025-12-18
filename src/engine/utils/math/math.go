package math

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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
