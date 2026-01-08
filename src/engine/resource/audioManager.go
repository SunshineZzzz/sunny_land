package resource

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"sunny_land/src/engine/utils"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/hajimehoshi/go-mp3"
	"github.com/jfreymuth/oggvorbis"
)

var (
	// 音效音量大小
	GSoundVolume = float32(1.0)
	// 音乐音量大小
	GMusicVolume = float32(1.0)
)

// 声音类型
type AudioType int

const (
	// 音效
	AudioTypeEffect AudioType = iota
	// 音乐
	AudioTypeMusic
)

// 音频抽象
type IAudio interface {
	// 播放声音
	Play() bool
	// 暂停声音
	Pause()
	// 恢复声音
	Resume()
	// 停止声音
	Stop()
	// 设置是否循环播放
	SetLoop(loop bool)
	// 关闭声音
	Close()
	// 获取音频类型
	GetAudioType() AudioType
	// 设置音量，范围0.0-1.0
	SetVolume(volume float32)
	// 获取当前音量
	GetVolume() float32
	// 是否正在播放
	IsPlaying() bool
}

// 全局音频句柄管理
var audioHandles = struct {
	sync.RWMutex
	handles map[uint32]IAudio
	nextID  uint32
}{

	handles: make(map[uint32]IAudio),
	nextID:  1,
}

// 注册音频
func registerAudio(audio IAudio) uint32 {
	audioHandles.Lock()
	defer audioHandles.Unlock()

	id := audioHandles.nextID
	if _, ok := audioHandles.handles[id]; ok {
		tryCount := 0
		tryMax := 10000
		//lint:ignore S1006 循环查找可用ID
		for true {
			tryCount++
			if tryCount >= tryMax {
				panic("register audio failed, max try count reached")
			}
			if _, ok := audioHandles.handles[id]; ok {
				audioHandles.nextID++
				id = audioHandles.nextID
				continue
			}
			break
		}
	}
	audioHandles.handles[id] = audio
	audioHandles.nextID++
	return id
}

// 获取音频
func getAudio(id uint32) IAudio {
	audioHandles.RLock()
	defer audioHandles.RUnlock()

	return audioHandles.handles[id]
}

// 注销音频
func unregisterAudio(id uint32) {
	audioHandles.Lock()
	defer audioHandles.Unlock()

	delete(audioHandles.handles, id)
}

// 创建音频
func newAudio(audioFilePath string, audioType AudioType) (IAudio, error) {
	extWithDot := filepath.Ext(audioFilePath)
	ext := strings.ToLower(extWithDot[1:])

	switch ext {
	case "ogg":
		return newOggAudio(audioFilePath, audioType)
	case "wav":
		return newWavAudio(audioFilePath, audioType)
	case "mp3":
		return newMp3Audio(audioFilePath, audioType)
	default:
		return nil, fmt.Errorf("unsupported audio file format: %s", extWithDot)
	}
}

// 音频管理器
type audioManager struct {
	// 存储所有加载的音效
	sounds map[string]*[]IAudio
	// 存储所有加载的音乐
	musics map[string]*[]IAudio
}

// 创建音频管理器
func NewAudioManager() *audioManager {
	slog.Debug("New AudioManager")
	return &audioManager{
		sounds: make(map[string]*[]IAudio),
		musics: make(map[string]*[]IAudio),
	}
}

// 清理
func (am *audioManager) Clear() {
	for _, sounds := range am.sounds {
		for _, sound := range *sounds {
			sound.Close()
		}
	}
	am.sounds = make(map[string]*[]IAudio)

	for _, musics := range am.musics {
		for _, music := range *musics {
			music.Close()
		}
	}
	am.musics = make(map[string]*[]IAudio)

	slog.Debug("audio manager clear")
}

// 加载音效
func (am *audioManager) loadSound(filePath string) *[]IAudio {
	audio, err := newAudio(filePath, AudioTypeEffect)
	if err != nil {
		slog.Error("load sound error", slog.String("path", filePath), slog.String("error", err.Error()))
		return nil
	}
	if _, ok := am.sounds[filePath]; !ok {
		am.sounds[filePath] = new([]IAudio)
	}
	*am.sounds[filePath] = append(*am.sounds[filePath], audio)
	return am.sounds[filePath]
}

// 获取音效
func (am *audioManager) GetSound(filePath string) *[]IAudio {
	sound, ok := am.sounds[filePath]
	if ok {
		return sound
	}
	slog.Debug("sound not in cache, try to load", slog.String("path", filePath))
	return am.loadSound(filePath)
}

// 卸载音效
func (am *audioManager) UnloadSound(filePath string) {
	if _, ok := am.sounds[filePath]; !ok {
		slog.Warn("sound not in cache, can not unload", slog.String("path", filePath))
		return
	}
	for _, sound := range *am.sounds[filePath] {
		sound.Close()
	}
	slog.Debug("unload sound", slog.String("path", filePath))
	delete(am.sounds, filePath)
}

// 加载音乐
func (am *audioManager) loadMusic(filePath string) *[]IAudio {
	audio, err := newAudio(filePath, AudioTypeMusic)
	if err != nil {
		slog.Error("load music error", slog.String("path", filePath), slog.String("error", err.Error()))
		return nil
	}
	if _, ok := am.musics[filePath]; !ok {
		am.musics[filePath] = new([]IAudio)
	}
	*am.musics[filePath] = append(*am.musics[filePath], audio)
	return am.musics[filePath]
}

// 获取音乐
func (am *audioManager) GetMusic(filePath string) *[]IAudio {
	music, ok := am.musics[filePath]
	if ok {
		return music
	}
	slog.Debug("music not in cache, try to load", slog.String("path", filePath))
	return am.loadMusic(filePath)
}

// 卸载音乐
func (am *audioManager) UnloadMusic(filePath string) {
	if _, ok := am.musics[filePath]; !ok {
		slog.Warn("music not in cache, can not unload", slog.String("path", filePath))
		return
	}
	for _, music := range *am.musics[filePath] {
		music.Close()
	}
	slog.Debug("unload music", slog.String("path", filePath))
	delete(am.musics, filePath)
}

// 设置音效音量
func (am *audioManager) SetSoundVolume(volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	GSoundVolume = volume
}

// 设置音乐音量
func (am *audioManager) SetMusicVolume(volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	GMusicVolume = volume
}

// 获取音效音量
func (am *audioManager) GetSoundVolume() float32 {
	return GSoundVolume
}

// 获取音乐音量
func (am *audioManager) GetMusicVolume() float32 {
	return GMusicVolume
}

// OGG格式音频
type oggAudio struct {
	// 锁
	sync.Mutex
	// 音频类型
	audioType AudioType
	// SDL音频流
	stream *sdl.AudioStream
	// ogg音频数据
	audioData []byte
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// id
	id uint32
	// 音频规格
	sampleRate int32
	channels   int32
	// 音量
	volume float32
}

var _ IAudio = (*oggAudio)(nil)

func newOggAudio(audioFilePath string, audioType AudioType) (*oggAudio, error) {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file, %v, %v", audioFilePath, err)
	}
	defer file.Close()

	oggReader, err := oggvorbis.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create oggvorbis reader, %v, %v", audioFilePath, err)
	}

	pcmData := make([]float32, 1024*1024)
	totalSamples := 0
	for {
		n, err := oggReader.Read(pcmData[totalSamples:])
		if err != nil && err.Error() != "EOF" {
			return nil, fmt.Errorf("failed to read oggvorbis data, %v, %v", audioFilePath, err)
		}
		if n == 0 {
			break
		}
		totalSamples += n
	}

	spec := &sdl.AudioSpec{
		Freq:     int32(oggReader.SampleRate()),
		Channels: int32(oggReader.Channels()),
		Format:   sdl.AudioF32,
	}

	volume := GSoundVolume
	if audioType == AudioTypeMusic {
		volume = GMusicVolume
	}

	callback := sdl.NewAudioStreamCallback(oggAudioCallback)
	ogg := &oggAudio{
		audioType:  audioType,
		audioData:  utils.Float32ToBytes(pcmData[:totalSamples]),
		dataPos:    0,
		isPlaying:  false,
		loop:       false,
		sampleRate: int32(oggReader.SampleRate()),
		channels:   int32(oggReader.Channels()),
		volume:     volume,
	}
	ogg.id = registerAudio(ogg)

	ogg.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(ogg.id)),
	)

	if ogg.stream == nil {
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return ogg, nil
}

// 音频回调函数
func oggAudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	ogg := getAudio(id).(*oggAudio)

	// 安全检查
	if ogg == nil {
		return
	}

	ogg.Lock()
	defer ogg.Unlock()

	if ogg.id != id || !ogg.isPlaying || len(ogg.audioData) == 0 {
		return
	}

	// 计算剩余数据量
	remaining := len(ogg.audioData) - ogg.dataPos
	if remaining <= 0 {
		if ogg.loop {
			ogg.dataPos = 0
			remaining = len(ogg.audioData)
		} else {
			ogg.isPlaying = false
			sdl.PauseAudioStreamDevice(stream)
			sdl.ClearAudioStream(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := ogg.audioData[ogg.dataPos : ogg.dataPos+dataToSend]
		sdl.PutAudioStreamData(stream, (*uint8)(unsafe.Pointer(&data[0])), int32(dataToSend))
		ogg.dataPos += dataToSend
	}

	// 再次检查循环（如果刚好发送完所有数据）
	if ogg.loop && ogg.dataPos >= len(ogg.audioData) {
		ogg.dataPos = 0
	}
}

// 播放
func (o *oggAudio) Play() bool {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return false
	}

	if o.isPlaying {
		// 这里肯定会进来，但是断不到点
		return false
	}

	o.isPlaying = true
	o.dataPos = 0
	sdl.ClearAudioStream(o.stream)
	sdl.SetAudioStreamGain(o.stream, o.volume)
	sdl.ResumeAudioStreamDevice(o.stream)

	return true
}

// 暂停
func (o *oggAudio) Pause() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = false
	sdl.PauseAudioStreamDevice(o.stream)
}

// 恢复
func (o *oggAudio) Resume() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = true
	sdl.ResumeAudioStreamDevice(o.stream)
}

// 停止
func (o *oggAudio) Stop() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = false
	o.dataPos = 0
	sdl.PauseAudioStreamDevice(o.stream)
	sdl.ClearAudioStream(o.stream)
}

// 设置循环播放
func (o *oggAudio) SetLoop(loop bool) {
	o.Lock()
	defer o.Unlock()

	o.loop = loop
}

// 关闭播放器，释放资源
func (o *oggAudio) Close() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	if o.stream != nil {
		o.isPlaying = false
		o.dataPos = 0
		sdl.PauseAudioStreamDevice(o.stream)
		sdl.ClearAudioStream(o.stream)
		sdl.DestroyAudioStream(o.stream)
		o.stream = nil
	}

	unregisterAudio(o.id)
}

// 获取音频类型
func (o *oggAudio) GetAudioType() AudioType {
	o.Lock()
	defer o.Unlock()

	return o.audioType
}

// 设置音量，范围0.0-1.0
func (o *oggAudio) SetVolume(volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.volume = volume
	// 设置音量
	sdl.SetAudioStreamGain(o.stream, o.volume)
}

// 获取当前音量
func (o *oggAudio) GetVolume() float32 {
	o.Lock()
	defer o.Unlock()

	return o.volume
}

// 是否正在播放
func (o *oggAudio) IsPlaying() bool {
	o.Lock()
	defer o.Unlock()

	return o.isPlaying
}

// wav格式音频
type wavAudio struct {
	// 锁
	sync.Mutex
	// 音频类型
	audioType AudioType
	// SDL音频流
	stream *sdl.AudioStream
	// wav音频数据
	audioBuf *uint8
	// wav音频数据长度
	audioLen int
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// 音频规格
	spec *sdl.AudioSpec
	// id
	id uint32
	// 音量
	volume float32
}

// 创建wav音频
func newWavAudio(soundFilePath string, audioType AudioType) (*wavAudio, error) {
	// 打开文件IO流
	ioStream := sdl.IOFromFile(soundFilePath, "rb")
	if ioStream == nil {
		return nil, fmt.Errorf("failed to open WAV file: %s", sdl.GetError())
	}
	// 自动释放了
	// defer sdl.CloseIO(ioStream)

	// 使用SDL直接加载WAV文件
	var audioBuf *uint8
	var audioLen uint32
	spec := &sdl.AudioSpec{}
	// 加载WAV数据
	success := sdl.LoadWAVIO(ioStream, true, spec, &audioBuf, &audioLen)
	if !success {
		return nil, fmt.Errorf("failed to load WAV data: %s", sdl.GetError())
	}

	volume := GSoundVolume
	if audioType == AudioTypeMusic {
		volume = GMusicVolume
	}

	wav := &wavAudio{
		audioType: audioType,
		audioBuf:  audioBuf,
		audioLen:  int(audioLen),
		spec:      spec,
		dataPos:   0,
		isPlaying: false,
		loop:      false,
		volume:    volume,
	}

	// 注册WAV播放器
	wav.id = registerAudio(wav)

	// 创建音频流
	callback := sdl.NewAudioStreamCallback(wavAudioCallback)
	wav.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(wav.id)),
	)

	if wav.stream == nil {
		sdl.Free(unsafe.Pointer(audioBuf))
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return wav, nil
}

// wav音频回调函数
func wavAudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	wav := getAudio(id).(*wavAudio)

	// 安全检查
	if wav == nil {
		return
	}

	wav.Lock()
	defer wav.Unlock()

	// fmt.Printf("wavAudioCallback id:%d, additionalAmount:%d, totalAmount:%d\n", id, additionalAmount, totalAmount)

	if wav.id != id || !wav.isPlaying || wav.audioLen == 0 {
		return
	}

	// 计算剩余数据量
	remaining := wav.audioLen - wav.dataPos
	if remaining <= 0 {
		if wav.loop {
			wav.dataPos = 0
			remaining = wav.audioLen
		} else {
			wav.isPlaying = false
			sdl.PauseAudioStreamDevice(stream)
			sdl.ClearAudioStream(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(wav.audioBuf)) + uintptr(wav.dataPos)))
		// wav.audioBuf+wav.dataPos
		sdl.PutAudioStreamData(stream, data, int32(dataToSend))
		wav.dataPos += dataToSend
	}

	// 再次检查循环
	if wav.loop && wav.dataPos >= wav.audioLen {
		wav.dataPos = 0
	}
}

// 播放控制方法（保持不变）
func (w *wavAudio) Play() bool {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil || w.id == 0 {
		return false
	}

	if w.isPlaying {
		// 这里肯定会进来，但是断不到点
		return false
	}

	w.isPlaying = true
	w.dataPos = 0
	sdl.ClearAudioStream(w.stream)
	sdl.SetAudioStreamGain(w.stream, w.volume)
	sdl.ResumeAudioStreamDevice(w.stream)

	return true
}

// 暂停
func (w *wavAudio) Pause() {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil || w.id == 0 {
		return
	}

	w.isPlaying = false
	sdl.PauseAudioStreamDevice(w.stream)
}

// 恢复
func (w *wavAudio) Resume() {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil || w.id == 0 {
		return
	}

	w.isPlaying = true
	sdl.ResumeAudioStreamDevice(w.stream)
}

// 停止
func (w *wavAudio) Stop() {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil || w.id == 0 {
		return
	}

	w.isPlaying = false
	w.dataPos = 0
	sdl.PauseAudioStreamDevice(w.stream)
	sdl.ClearAudioStream(w.stream)
}

// 设置循环播放
func (w *wavAudio) SetLoop(loop bool) {
	w.Lock()
	defer w.Unlock()

	w.loop = loop
}

// 关闭播放器，释放资源
func (w *wavAudio) Close() {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil || w.id == 0 {
		return
	}

	if w.stream != nil {
		w.isPlaying = false
		w.dataPos = 0
		sdl.PauseAudioStreamDevice(w.stream)
		sdl.ClearAudioStream(w.stream)
		sdl.DestroyAudioStream(w.stream)
		w.stream = nil
	}

	if w.audioLen > 0 {
		sdl.Free(unsafe.Pointer(w.audioBuf))
		w.audioBuf = nil
		w.audioLen = 0
	}

	unregisterAudio(w.id)
}

// 获取音频类型
func (w *wavAudio) GetAudioType() AudioType {
	w.Lock()
	defer w.Unlock()

	return w.audioType
}

// 设置音量
func (w *wavAudio) SetVolume(volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	w.Lock()
	defer w.Unlock()

	w.volume = volume
	sdl.SetAudioStreamGain(w.stream, volume)
}

// 获取当前音量
func (w *wavAudio) GetVolume() float32 {
	w.Lock()
	defer w.Unlock()

	return w.volume
}

// 是否正在播放
func (w *wavAudio) IsPlaying() bool {
	w.Lock()
	defer w.Unlock()

	return w.isPlaying
}

// mp3格式音频
type mp3Audio struct {
	// 锁
	sync.Mutex
	// 音频类型
	audioType AudioType
	// SDL音频流
	stream *sdl.AudioStream
	// MP3音频数据
	audioData []byte
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// id
	id uint32
	// 音频规格
	sampleRate int32
	channels   int32
	// 音量
	volume float32
}

var _ IAudio = (*mp3Audio)(nil)

func newMp3Audio(audioFilePath string, audioType AudioType) (*mp3Audio, error) {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file, %v, %v", audioFilePath, err)
	}
	defer file.Close()

	d, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create oggvorbis reader, %v, %v", audioFilePath, err)
	}

	// 动态读取完整数据
	var pcmData []byte
	buf := make([]byte, 4096)
	for {
		n, err := d.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return nil, fmt.Errorf("failed to read mp3 data, %v, %v", audioFilePath, err)
		}
		if n > 0 {
			pcmData = append(pcmData, buf[:n]...)
		}
		if n == 0 {
			break
		}
	}

	spec := &sdl.AudioSpec{
		Freq:     int32(d.SampleRate()),
		Channels: 2,
		Format:   sdl.AudioS16,
	}

	volume := GSoundVolume
	if audioType == AudioTypeMusic {
		volume = GMusicVolume
	}

	callback := sdl.NewAudioStreamCallback(mp3AudioCallback)
	mp3 := &mp3Audio{
		audioType:  audioType,
		audioData:  pcmData,
		dataPos:    0,
		isPlaying:  false,
		loop:       false,
		sampleRate: int32(d.SampleRate()),
		channels:   2,
		volume:     volume,
	}

	mp3.id = registerAudio(mp3)

	mp3.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(mp3.id)),
	)

	if mp3.stream == nil {
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return mp3, nil
}

// 音频回调函数
func mp3AudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	mp3 := getAudio(id).(*mp3Audio)

	// 安全检查
	if mp3 == nil {
		return
	}

	mp3.Lock()
	defer mp3.Unlock()

	if mp3.id != id || !mp3.isPlaying || len(mp3.audioData) == 0 {
		return
	}

	// 计算剩余数据量
	remaining := len(mp3.audioData) - mp3.dataPos
	if remaining <= 0 {
		if mp3.loop {
			mp3.dataPos = 0
			remaining = len(mp3.audioData)
		} else {
			mp3.isPlaying = false
			sdl.PauseAudioStreamDevice(stream)
			sdl.ClearAudioStream(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := mp3.audioData[mp3.dataPos : mp3.dataPos+dataToSend]
		sdl.PutAudioStreamData(stream, (*uint8)(unsafe.Pointer(&data[0])), int32(dataToSend))
		mp3.dataPos += dataToSend
	}

	// 再次检查循环（如果刚好发送完所有数据）
	if mp3.loop && mp3.dataPos >= len(mp3.audioData) {
		mp3.dataPos = 0
	}
}

// 播放
func (o *mp3Audio) Play() bool {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return false
	}

	if o.isPlaying {
		// fmt.Printf("mp3Sound %d is already playing\n", o.id)
		// 这里肯定会进来，但是断不到点
		return false
	}

	// fmt.Printf("mp3Sound %d will playing\n", o.id)

	o.isPlaying = true
	o.dataPos = 0
	sdl.ClearAudioStream(o.stream)
	sdl.SetAudioStreamGain(o.stream, o.volume)
	sdl.ResumeAudioStreamDevice(o.stream)

	return true
}

// 暂停
func (o *mp3Audio) Pause() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = false
	sdl.PauseAudioStreamDevice(o.stream)
}

// 恢复
func (o *mp3Audio) Resume() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = true
	sdl.ResumeAudioStreamDevice(o.stream)
}

// 停止
func (o *mp3Audio) Stop() {
	o.Lock()
	defer o.Unlock()

	if o.stream == nil || o.id == 0 {
		return
	}

	o.isPlaying = false
	o.dataPos = 0
	sdl.PauseAudioStreamDevice(o.stream)
	sdl.ClearAudioStream(o.stream)
}

// 设置循环播放
func (o *mp3Audio) SetLoop(loop bool) {
	o.Lock()
	defer o.Unlock()

	o.loop = loop
}

// 关闭播放器，释放资源
func (o *mp3Audio) Close() {
	o.Lock()
	defer o.Unlock()

	// 安全检查
	if o.stream == nil || o.id == 0 {
		return
	}

	if o.stream != nil {
		o.isPlaying = false
		o.dataPos = 0
		sdl.PauseAudioStreamDevice(o.stream)
		sdl.ClearAudioStream(o.stream)
		sdl.DestroyAudioStream(o.stream)
		o.stream = nil
	}

	unregisterAudio(o.id)
}

// 获取音频类型
func (o *mp3Audio) GetAudioType() AudioType {
	o.Lock()
	defer o.Unlock()

	return o.audioType
}

// 设置音量
func (o *mp3Audio) SetVolume(volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	o.Lock()
	defer o.Unlock()

	o.volume = volume
	sdl.SetAudioStreamGain(o.stream, volume)
}

// 获取当前音量
func (o *mp3Audio) GetVolume() float32 {
	o.Lock()
	defer o.Unlock()

	return o.volume
}

// 是否正在播放
func (o *mp3Audio) IsPlaying() bool {
	o.Lock()
	defer o.Unlock()

	return o.isPlaying
}
