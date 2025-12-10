package core

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 帧率控制器
type FPS struct {
	// 上一帧时间
	lastFrameTime uint64
	// 当前帧开始时间
	currentFrameStartTime uint64
	//  delta time
	deltaTime float64
	// 时间缩放因子
	timeScale float64
	// 目标帧率
	targetFps int
	// 目标每帧时间
	targetFrameTime float64
}

// 创建帧率控制器
func NewFPS() *FPS {
	fps := &FPS{}
	fps.timeScale = 1.0
	fps.lastFrameTime = sdl.GetTicksNS()
	fps.currentFrameStartTime = fps.lastFrameTime

	slog.Debug("create fps controller", slog.Uint64("lastFrameTime", fps.lastFrameTime))

	return fps
}

// 更新帧率控制器
func (f *FPS) Update() {
	f.currentFrameStartTime = sdl.GetTicksNS()
	currentDeltaTime := float64(f.currentFrameStartTime-f.lastFrameTime) / 1e9
	if f.targetFrameTime > 0.0 {
		// 设置了帧率，则需要限制帧率
		f.limitFrameRate(currentDeltaTime)
	} else {
		f.deltaTime = currentDeltaTime
	}
	f.lastFrameTime = sdl.GetTicksNS()
}

// 限制帧率
func (f *FPS) limitFrameRate(currentDeltaTime float64) {
	if currentDeltaTime < f.targetFrameTime {
		// 当前帧耗费的时间小于目标每帧时间，需要等待
		timeToWait := f.targetFrameTime - currentDeltaTime
		nsToWait := uint64(timeToWait * 1e9)
		sdl.DelayNS(nsToWait)
		f.deltaTime = float64(sdl.GetTicksNS()-f.lastFrameTime) / 1e9
		return
	}
	f.deltaTime = currentDeltaTime
}

// 获取当前帧的delta time
func (f *FPS) GetDeltaTime() float64 {
	return f.deltaTime * f.timeScale
}

// 获取没有缩放的delta time
func (f *FPS) GetUnscaledDeltaTime() float64 {
	return f.deltaTime
}

// 设置时间缩放因子
func (f *FPS) SetTimeScale(timeScale float64) {
	f.timeScale = timeScale
}

// 获取时间缩放因子
func (f *FPS) GetTimeScale() float64 {
	return f.timeScale
}

// 设置目标帧率
func (f *FPS) SetTargetFps(targetFps int) {
	if targetFps < 0 {
		slog.Error("target fps must be greater than 0, set to 0")
		f.targetFps = 0
		return
	}

	f.targetFps = targetFps
	if f.targetFps > 0 {
		f.targetFrameTime = 1.0 / float64(f.targetFps)
	}

	slog.Info("set target fps", slog.Int("targetFps", f.targetFps), slog.Float64("targetFrameTime", f.targetFrameTime))
}

// 获取目标帧率
func (f *FPS) GetTargetFps() int {
	return f.targetFps
}
