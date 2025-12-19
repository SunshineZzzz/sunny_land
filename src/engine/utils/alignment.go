package utils

type Alignment int

const (
	// 不指定对齐方式，偏移量通常为 (0,0) 或手动设置
	AlignNone Alignment = iota
	// 左上角
	AlignTopLeft
	// 顶部中心
	AlignTopCenter
	// 右上角
	AlignTopRight
	// 中心左侧
	AlignCenterLeft
	// 正中心(几何中心)
	AlignCenter
	// 中心右侧
	AlignCenterRight
	// 左下角
	AlignBottomLeft
	// 底部中心
	AlignBottomCenter
	// 右下角
	AlignBottomRight
)
