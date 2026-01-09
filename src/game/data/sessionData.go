package data

import (
	"log/slog"
	"os"

	emath "sunny_land/src/engine/utils/math"

	"github.com/bitly/go-simplejson"
)

/**
 * @brief 管理不同游戏场景之间的游戏状态
 *
 * 存储玩家生命值、分数、当前关卡等信息，
 * 使这些数据在场景切换时能够保持。
 */
type SessionData struct {
	// 当前生命值
	currentHealth int
	// 最大生命值
	maxHealth int
	// 当前得分
	currentScore int
	// 最高得分
	highScore int
	// 是否胜利
	isWin bool

	// 进入关卡时的生命值，读/存档用
	levelHealth int
	// 进入关卡时的得分，读/存档用
	levelScore int
	// 进入关卡时的地图路径，读/存档用
	mapPath string
}

// 创建新的会话数据
func NewSessionData() *SessionData {
	return &SessionData{
		currentHealth: 3,
		maxHealth:     3,
		currentScore:  0,
		highScore:     0,
		levelHealth:   3,
		levelScore:    0,
		mapPath:       "assets/maps/level1.tmj",
	}
}

// 获取当前生命值
func (sd *SessionData) GetCurHealth() int {
	return sd.currentHealth
}

// 获取最大生命值
func (sd *SessionData) GetMaxHealth() int {
	return sd.maxHealth
}

// 获取当前得分
func (sd *SessionData) GetCurrentScore() int {
	return sd.currentScore
}

// 获取最高得分
func (sd *SessionData) GetHighScore() int {
	return sd.highScore
}

// 获取进入关卡时的生命值
func (sd *SessionData) GetLevelHealth() int {
	return sd.levelHealth
}

// 获取进入关卡时的得分
func (sd *SessionData) GetLevelScore() int {
	return sd.levelScore
}

// 获取进入关卡时的地图路径
func (sd *SessionData) GetMapPath() string {
	return sd.mapPath
}

// 设置当前生命值
func (sd *SessionData) SetCurrentHealth(health int) {
	sd.currentHealth = emath.Clamp(health, 0, sd.maxHealth)
}

// 设置最大生命值
func (sd *SessionData) SetMaxHealth(health int) {
	if health <= 0 {
		sd.maxHealth = 1
	} else {
		sd.maxHealth = health
	}
	sd.SetCurrentHealth(sd.currentHealth)
}

// 设置最高得分
func (sd *SessionData) SetHighScore(score int) {
	sd.highScore = score
}

// 设置进入关卡时的生命值
func (sd *SessionData) SetLevelHealth(health int) {
	sd.levelHealth = health
}

// 设置进入关卡时的得分
func (sd *SessionData) SetLevelScore(score int) {
	sd.levelScore = score
}

// 设置进入关卡时的地图路径
func (sd *SessionData) SetMapPath(path string) {
	sd.mapPath = path
}

// 增加分数
func (sd *SessionData) AddScore(score int) {
	sd.currentScore += score
	// 如果当前分数超过最高分，则更新最高分
	sd.SetHighScore(max(sd.highScore, sd.currentScore))
}

// 重置游戏数据以准备开始新游戏，保留最高分
func (sd *SessionData) Reset() {
	sd.currentHealth = sd.maxHealth
	sd.currentScore = 0
	sd.levelHealth = 3
	sd.levelScore = 0
	sd.mapPath = "assets/maps/level1.tmj"
	slog.Info("session data reset",
		slog.Int("currentHealth", sd.currentHealth),
		slog.Int("maxHealth", sd.maxHealth),
		slog.Int("currentScore", sd.currentScore),
		slog.Int("highScore", sd.highScore),
		slog.Int("levelHealth", sd.levelHealth),
		slog.Int("levelScore", sd.levelScore),
		slog.String("mapPath", sd.mapPath),
	)
}

// 设置下一个关卡的信息，包括地图路径、关卡开始时的得分生命
func (sd *SessionData) SetNextLevel(mapPath string) {
	sd.mapPath = mapPath
	sd.levelHealth = sd.currentHealth
	sd.levelScore = sd.currentScore
}

// 将当前游戏数据保存到JSON文件，存档
func (sd *SessionData) SaveToFile(filename string) bool {
	j := simplejson.New()
	// 将成员变量序列化到 JSON 对象中
	j.Set("level_score", sd.levelScore)
	j.Set("level_health", sd.levelHealth)
	j.Set("max_health", sd.maxHealth)
	j.Set("high_score", sd.highScore)
	j.Set("map_path", sd.mapPath)

	// 打开文件进行写入
	ofs, err := os.Create(filename)
	if err != nil {
		slog.Error("create file error", slog.String("filename", filename))
		return false
	}
	defer ofs.Close()

	// 将JSON对象写入文件
	jsonBytes, err := j.EncodePretty()
	if err != nil {
		slog.Error("encode json error", slog.String("filename", filename))
		return false
	}

	_, err = ofs.Write(jsonBytes)
	if err != nil {
		slog.Error("write file error", slog.String("filename", filename))
		return false
	}

	slog.Info("session data save to file", slog.String("filename", filename))
	return true
}

// 从JSON文件加载游戏数据，读档
func (sd *SessionData) LoadFromFile(filename string) bool {
	// 打开文件进行读取
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		slog.Warn("read file error", slog.String("filename", filename))
		// 如果存档文件不存在，这不一定是错误
		return false
	}

	// 从文件解析JSON数据
	j, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		slog.Error("decode json error", slog.String("filename", filename))
		return false
	}

	sd.currentScore = j.Get("level_score").MustInt(0)
	sd.currentHealth = j.Get("level_health").MustInt(3)
	sd.maxHealth = j.Get("max_health").MustInt(3)
	// 文件中的最高分，与当前最高分取最大值
	sd.highScore = max(sd.highScore, j.Get("high_score").MustInt(0))
	sd.mapPath = j.Get("map_path").MustString("assets/maps/level1.tmj")

	slog.Info("session data load from file", slog.String("filename", filename))
	return true
}

// 设置是否胜利
func (sd *SessionData) SetIsWin(isWin bool) {
	sd.isWin = isWin
}

// 获取是否胜利
func (sd *SessionData) GetIsWin() bool {
	return sd.isWin
}

// 同步最高分(文件与当前分数取最大值)
func (sd *SessionData) SyncHighScore(fileName string) bool {
	jsonBytes, err := os.ReadFile(fileName)
	if err != nil {
		slog.Warn("read file error", slog.String("filename", fileName))
		// 文件找不到无法同步
		return false
	}

	// 从文件解析JSON数据
	j, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		slog.Error("decode json error", slog.String("filename", fileName))
		return false
	}

	// 同步最高分
	highScoreInFile := j.Get("high_score").MustInt(0)
	// 根据文件中的最高分和当前最高分来决定处理方式
	if highScoreInFile < sd.highScore {
		// 文件中的最高分 低于 当前最高分
		j.Set("high_score", sd.highScore)
		jsonBytes, err = j.EncodePretty()
		if err != nil {
			slog.Error("encode json error", slog.String("filename", fileName))
			return false
		}
		os.WriteFile(fileName, jsonBytes, 0644)
		slog.Info("session data sync high score to file", slog.String("filename", fileName))
	} else {
		sd.highScore = highScoreInFile
		// 文件中的最高分 不低于 当前最高分
		slog.Info("session data high score not lower than current score", slog.String("filename", fileName))
	}

	return true
}
