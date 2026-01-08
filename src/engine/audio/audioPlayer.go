package audio

import (
	"log/slog"
	"sunny_land/src/engine/resource"
)

/**
 * @brief 用于控制音频播放的单例类。
 *
 * 提供播放音效和音乐的方法，使用由 ResourceManager 管理的资源。
 * 必须使用有效的 ResourceManager 实例初始化。
 */
type AudioPlayer struct {
	// 资源管理器
	resourceManager *resource.ResourceManager
	// 当前正在播放的音乐路径，用于避免重复播放同一音乐
	currentMusicPath string
	// 当前正在播放音乐对象
	currentMusicAudio resource.IAudio
}

// 创建音乐播放器
func NewAudioPlayer(resourceManager *resource.ResourceManager) *AudioPlayer {
	if resourceManager == nil {
		panic("audio player resource manager is nil")
	}
	return &AudioPlayer{
		resourceManager:  resourceManager,
		currentMusicPath: "",
	}
}

// 播放音效
func (a *AudioPlayer) PlaySound(soundPath string) bool {
	sounds := a.resourceManager.GetSound(soundPath)
	if sounds == nil || len(*sounds) == 0 {
		slog.Error("audio player play sound error1", slog.String("soundPath", soundPath))
		return false
	}

	for _, sound := range *sounds {
		sound.SetLoop(false)
		if sound.Play() {
			return true
		}
	}

	// 走到这里肯定是所有池子都失败了
	if sounds = a.resourceManager.LoadSound(soundPath); sounds == nil || len(*sounds) == 0 {
		slog.Error("audio player play sound error2", slog.String("soundPath", soundPath))
		return false
	}

	(*sounds)[len(*sounds)-1].SetLoop(false)
	return (*sounds)[len(*sounds)-1].Play()
}

// 播放音乐
func (a *AudioPlayer) PlayMusic(musicPath string, loop bool) bool {
	// 如果当前音乐已经在播放，则不重复播放
	if a.currentMusicPath == musicPath && a.currentMusicAudio != nil && a.currentMusicAudio.IsPlaying() {
		return true
	}

	a.currentMusicPath = musicPath
	if a.currentMusicAudio != nil {
		a.currentMusicAudio.Close()
		a.currentMusicAudio = nil
	}

	musics := a.resourceManager.GetMusic(musicPath)
	if musics == nil || len(*musics) == 0 {
		slog.Error("audio player play music error1", slog.String("musicPath", musicPath))
		return false
	}

	for _, music := range *musics {
		music.SetLoop(loop)
		if music.Play() {
			a.currentMusicAudio = music
			return true
		}
	}

	// 走到这里肯定是所有池子都失败了
	if musics = a.resourceManager.LoadMusic(musicPath); musics == nil || len(*musics) == 0 {
		slog.Error("audio player play music error2", slog.String("musicPath", musicPath))
		return false
	}

	(*musics)[len(*musics)-1].SetLoop(true)
	if (*musics)[len(*musics)-1].Play() {
		a.currentMusicAudio = (*musics)[len(*musics)-1]
		return true
	}

	slog.Error("audio player play music error3", slog.String("musicPath", musicPath))
	return false
}

// 停止音乐播放
func (a *AudioPlayer) StopMusic() {
	if a.currentMusicAudio == nil {
		return
	}

	a.currentMusicAudio.Stop()
}

// 暂停音乐播放
func (a *AudioPlayer) PauseMusic() {
	if a.currentMusicAudio == nil {
		return
	}

	a.currentMusicAudio.Pause()
}

// 恢复音乐播放
func (a *AudioPlayer) ResumeMusic() {
	if a.currentMusicAudio == nil {
		return
	}

	a.currentMusicAudio.Resume()
}

// 设置音乐音量
func (a *AudioPlayer) SetMusicVolume(volume float32) {
	if a.currentMusicAudio == nil {
		a.resourceManager.SetMusicVolume(volume)
		return
	}

	a.currentMusicAudio.SetVolume(volume)
}

// 获取当前音乐音量
func (a *AudioPlayer) GetMusicVolume() float32 {
	if a.currentMusicAudio == nil {
		return a.resourceManager.GetMusicVolume()
	}

	return a.currentMusicAudio.GetVolume()
}

// 设置音效音量
func (a *AudioPlayer) SetSoundVolume(volume float32) {
	a.resourceManager.SetSoundVolume(volume)
}

// 获取当前音效音量
func (a *AudioPlayer) GetSoundVolume() float32 {
	return a.resourceManager.GetSoundVolume()
}
