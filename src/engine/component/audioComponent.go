package component

import (
	"log/slog"

	"sunny_land/src/engine/audio"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/utils/def"
)

/**
 * @brief 音频组件，用于处理音频播放和管理。
 */
type AudioComponent struct {
	// 继承组件基类
	Component
	// 音频播放器的非拥有指针
	audioPlayer *audio.AudioPlayer
	// 相机的非拥有指针，用于音频空间定位
	camera *render.Camera
	// 缓存变换组件的非拥有指针，用于音频空间定位
	transformComponent *TransformComponent
	// 音效Id到路径的映射表
	soundIdToPathMap map[string]string
}

// 确保AudioComponent实现了IComponent接口
var _ physics.IComponent = (*AudioComponent)(nil)

// 创建音频组件
func NewAudioComponent(audioPlayer *audio.AudioPlayer, camera *render.Camera) *AudioComponent {
	if audioPlayer == nil || camera == nil {
		slog.Error("AudioComponent Init: audioPlayer or camera is nil")
		return nil
	}

	return &AudioComponent{
		Component: Component{
			ComponentType: def.ComponentTypeAudio,
		},
		audioPlayer:      audioPlayer,
		camera:           camera,
		soundIdToPathMap: make(map[string]string),
	}
}

// 初始化
func (a *AudioComponent) Init() {
	if a.Owner == nil {
		slog.Error("AudioComponent Init: owner is nil")
		return
	}
	a.transformComponent = a.Owner.GetComponent(def.ComponentTypeTransform).(*TransformComponent)
	if a.transformComponent == nil {
		slog.Error("AudioComponent Init: transform component is nil")
		return
	}
}

// 播放音效
func (a *AudioComponent) PlaySound(soundId string, useSpatial bool) {
	soundPath, ok := a.soundIdToPathMap[soundId]
	if !ok {
		soundPath = soundId
	}

	// 是否使用空间定位
	if useSpatial && a.transformComponent != nil {
		// 这里给一个简单的功能：150像素范围内播放，否则不播放
		// 相机中心
		cameraCenter := a.camera.GetPosition().Add(a.camera.GetViewportSize().Mul(0.5))
		objPos := a.transformComponent.GetPosition()
		// 距离
		dist := cameraCenter.Sub(objPos).Len()
		if dist > 150.0 {
			return
		}
		a.audioPlayer.PlaySound(soundPath)
	} else {
		a.audioPlayer.PlaySound(soundPath)
	}
}

/**
 * @brief 添加音效到映射表。
 * @param sound_id 音效的标识符（针对本组件唯一即可）。
 * @param sound_path 音效文件的路径。
 */
func (a *AudioComponent) AddSound(soundId string, soundPath string) {
	a.soundIdToPathMap[soundId] = soundPath
}
