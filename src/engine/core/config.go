package core

import (
	"encoding/json"
	"log/slog"
	"os"
)

// Config对应JSON的根结构
type configJson struct {
	Window        windowConfig        `json:"window"`
	Graphics      graphicsConfig      `json:"graphics"`
	Performance   performanceConfig   `json:"performance"`
	Audio         audioConfig         `json:"audio"`
	InputMappings map[string][]string `json:"input_mappings"`
}

// WindowConfig对应"window"字段
type windowConfig struct {
	Title     string `json:"title"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Resizable bool   `json:"resizable"`
}

// GraphicsConfig对应"graphics"字段
type graphicsConfig struct {
	Vsync bool `json:"vsync"`
}

// PerformanceConfig对应"performance"字段
type performanceConfig struct {
	TargetFPS int `json:"target_fps"`
}

// AudioConfig对应"audio"字段
type audioConfig struct {
	MusicVolume float32 `json:"music_volume"`
	SoundVolume float32 `json:"sound_volume"`
}

// 管理应用程序配置
type Config struct {
	// 窗口标题
	WindowTitle string
	// 窗口宽度
	WindowWidth int
	// 窗口高度
	WindowHeight int
	// 是否可调整窗口大小
	WindowResizable bool
	// 是否开启垂直同步
	VsyncEnabled bool
	// 目标帧率
	TargetFPS int
	// 音效大小
	SoundVolume float32
	// 音乐大小
	MusicVolume float32
	// 按键映射
	InputMappings map[string][]string
}

// 创建配置
func NewConfig(filePath string) *Config {
	config := &Config{}
	config.Init()
	config.LoadFromFile(filePath)
	return config
}

// 初始化
func (c *Config) Init() {
	c.WindowTitle = "SunnyLand"
	c.WindowWidth = 1280
	c.WindowHeight = 720
	c.WindowResizable = true
	c.VsyncEnabled = true
	c.TargetFPS = 144
	c.SoundVolume = 0.5
	c.MusicVolume = 0.5
	c.InputMappings = make(map[string][]string)

	// 一些默认按键映射
	c.InputMappings["move_up"] = []string{"W", "Up"}
	c.InputMappings["move_down"] = []string{"S", "Down"}
	c.InputMappings["move_left"] = []string{"A", "Left"}
	c.InputMappings["move_right"] = []string{"D", "Right"}
	c.InputMappings["jump"] = []string{"J", "Space"}
	c.InputMappings["pause"] = []string{"P", "Escape"}
	c.InputMappings["attack"] = []string{"K", "MouseLeft"}
}

// 从文件中加载配置
func (c *Config) LoadFromFile(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Warn("load config file failed, use default config create file", slog.String("filePath", filePath), slog.String("err", err.Error()))
		if !c.SaveToFile(filePath) {
			slog.Error("default configuration file cannot be created", slog.String("filePath", filePath))
			return false
		}
		return false
	}

	var config configJson
	if err := json.Unmarshal(data, &config); err != nil {
		slog.Error("unmarshal config file failed", slog.String("filePath", filePath), slog.String("err", err.Error()))
		return false
	}

	// 应用解析后的配置
	c.WindowTitle = config.Window.Title
	c.WindowWidth = config.Window.Width
	c.WindowHeight = config.Window.Height
	c.WindowResizable = config.Window.Resizable
	c.VsyncEnabled = config.Graphics.Vsync
	c.TargetFPS = config.Performance.TargetFPS
	c.SoundVolume = config.Audio.SoundVolume
	c.MusicVolume = config.Audio.MusicVolume
	c.InputMappings = config.InputMappings

	slog.Info("load config file success", slog.String("filePath", filePath))
	return true
}

// 保存配置到文件
func (c *Config) SaveToFile(filePath string) bool {
	configJson := configJson{
		Window: windowConfig{
			Title:     c.WindowTitle,
			Width:     c.WindowWidth,
			Height:    c.WindowHeight,
			Resizable: c.WindowResizable,
		},
		Graphics: graphicsConfig{
			Vsync: c.VsyncEnabled,
		},
		Performance: performanceConfig{
			TargetFPS: c.TargetFPS,
		},
		Audio: audioConfig{
			MusicVolume: c.MusicVolume,
			SoundVolume: c.SoundVolume,
		},
		InputMappings: c.InputMappings,
	}
	data, err := json.MarshalIndent(configJson, "", "  ")
	if err != nil {
		slog.Error("marshal config file failed", slog.String("filePath", filePath), slog.String("err", err.Error()))
		return false
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		slog.Error("write config file failed", slog.String("filePath", filePath), slog.String("err", err.Error()))
		return false
	}

	slog.Info("save config file success", slog.String("filePath", filePath))
	return true
}
