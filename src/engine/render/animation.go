package render

import (
	"log/slog"

	"sunny_land/src/engine/physics"
	emath "sunny_land/src/engine/utils/math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 管理一系列动画帧
type Animation struct {
	// 动画的名称
	name string
	// 动画中的所有帧
	frames []*physics.AnimationFrame
	// 动画总持续时间(秒)
	totalDuration float64
	// 是否循环播放
	loop bool
}

// 确保Animation实现了IAnimation接口
var _ physics.IAnimation = (*Animation)(nil)

// 创建动画
func NewAnimation(name string, loop bool) *Animation {
	return &Animation{
		name:          name,
		frames:        make([]*physics.AnimationFrame, 0),
		totalDuration: 0.0,
		loop:          loop,
	}
}

// 向动画添加一帧
func (a *Animation) AddFrame(rect *sdl.FRect, duration float64) {
	a.frames = append(a.frames, &physics.AnimationFrame{
		SourceRect: rect,
		Duration:   duration,
	})
	a.totalDuration += duration
}

// 获取在给定时间点应该显示的动画帧
func (a *Animation) GetFrameAtTime(time float64) *physics.AnimationFrame {
	if len(a.frames) == 0 {
		slog.Error("animation has no frames", slog.String("name", a.name))
		return nil
	}

	currentTime := time
	if a.loop && a.totalDuration > 0.0 {
		// 对循环动画使用取模运算获取有效时间
		currentTime = emath.Mod(time, a.totalDuration)
	} else {
		// 对于非循环动画，如果时间超过总时长，返回最后一帧
		if time >= a.totalDuration {
			return a.frames[len(a.frames)-1]
		}
	}

	// 遍历所有帧，找到当前时间点所属的帧
	var accumulatedTime float64
	for _, frame := range a.frames {
		accumulatedTime += frame.Duration
		if currentTime < accumulatedTime {
			return frame
		}
	}

	slog.Warn("time error of animation duration", slog.String("name", a.name), slog.Float64("time", time))
	return a.frames[len(a.frames)-1]
}

// 获取动画名称
func (a *Animation) GetName() string {
	return a.name
}

// 设置动画名称
func (a *Animation) SetName(name string) {
	a.name = name
}

// 获取动画帧列表
func (a *Animation) GetFrames() []*physics.AnimationFrame {
	return a.frames
}

// 获取帧数量
func (a *Animation) GetFrameCount() int {
	return len(a.frames)
}

// 获取动画总持续时间(秒)
func (a *Animation) GetTotalDuration() float64 {
	return a.totalDuration
}

// 是否循环播放
func (a *Animation) IsLooping() bool {
	return a.loop
}

// 设置是否循环播放
func (a *Animation) SetLooping(loop bool) {
	a.loop = loop
}

// 检查动画是否没有帧
func (a *Animation) IsEmpty() bool {
	return len(a.frames) == 0
}
