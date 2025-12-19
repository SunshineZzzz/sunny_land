package input

import (
	"log/slog"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 定义输入动作的状态
type ActionState int

const (
	// 动作未激活
	ActionStateInActive ActionState = iota
	// 动作在本帧刚刚被按下
	ActionStatePressedThisFrame
	// 动作被持续按下
	ActionStateHeldDown
	// 动作在本帧刚刚被释放
	ActionStateReleasedThisFrame
)

// 输入管理器
type InputManager struct {
	// SDL渲染器，用于获取逻辑坐标
	sdlRenderer *sdl.Renderer
	// 存储动作名称到按键名称列表的映射，比如:
	// "attack" = ["K","MouseLeft"]
	// "MouseLeftClick" = ["MouseLeft"]
	actionsToKeynameMap *map[string][]string
	// 从输入(SDL_Scancode或SDL_MouseButton)到关联的动作名称列表，比如:
	// 1 = ["attack","MouseLeftClick"]
	// 14 = ["attack"]
	inputToActionsMap map[uint32][]string
	// 存储每个动作的当前状态，比如:
	// "attack" = ActionStateHeldDown
	// "MouseLeftClick" = ActionStateInActive
	actionStates map[string]ActionState
	// 退出标志
	shouldQuit bool
	// 鼠标位置(屏幕坐标)
	mousePosition mgl32.Vec2
}

// 创建输入管理器
func NewInputManager(sdlRenderer *sdl.Renderer, inputMappings *map[string][]string) *InputManager {
	if sdlRenderer == nil {
		panic("sdlRenderer is nil")
	}
	im := &InputManager{
		sdlRenderer:         sdlRenderer,
		actionsToKeynameMap: nil,
		inputToActionsMap:   make(map[uint32][]string),
		actionStates:        make(map[string]ActionState),
		shouldQuit:          false,
	}
	im.initializeMappings(inputMappings)
	// 获取初始鼠标位置
	sdl.GetMouseState(&im.mousePosition[0], &im.mousePosition[1])
	slog.Debug("create input manager", slog.Any("mousePosition", im.mousePosition))
	return im
}

// 初始化按键映射
func (im *InputManager) initializeMappings(inputMappings *map[string][]string) {
	if inputMappings == nil {
		panic("inputMappings is nil")
	}
	im.actionsToKeynameMap = inputMappings
	im.inputToActionsMap = make(map[uint32][]string)
	im.actionStates = make(map[string]ActionState)

	// 如果配置中没有定义鼠标按钮动作(通常不需要配置)，则添加默认映射, 用于UI
	if _, exists := (*im.actionsToKeynameMap)["MouseLeftClick"]; !exists {
		slog.Debug("MouseLeftClick not found in input mappings, add default mapping")
		(*im.actionsToKeynameMap)["MouseLeftClick"] = []string{"MouseLeft"}
	}
	if _, exists := (*im.actionsToKeynameMap)["MouseRightClick"]; !exists {
		slog.Debug("MouseRightClick not found in input mappings, add default mapping")
		(*im.actionsToKeynameMap)["MouseRightClick"] = []string{"MouseRight"}
	}

	// 遍历动作<->按键名称映射，构建按键<->动作映射
	for action, keyNames := range *im.actionsToKeynameMap {
		im.actionStates[action] = ActionStateInActive
		for _, keyName := range keyNames {
			scancode := sdl.GetScancodeFromName(keyName)
			mouseButton := im.mouseButtonFromName(keyName)
			if scancode != sdl.ScancodeUnknown {
				im.inputToActionsMap[uint32(scancode)] = append(im.inputToActionsMap[uint32(scancode)], action)
				slog.Debug("map action to keyboard scancode", slog.String("action", action), slog.String("keyName", keyName),
					slog.Any("scancode", scancode))
			} else if mouseButton != 0 {
				im.inputToActionsMap[uint32(mouseButton)] = append(im.inputToActionsMap[uint32(mouseButton)], action)
				slog.Debug("map action to mouse button", slog.String("action", action), slog.String("keyName", keyName),
					slog.Any("mouseButton", mouseButton))
			} else {
				slog.Warn("unknown key name in input mappings, ignore", slog.String("keyName", keyName), slog.String("action", action))
			}
		}
	}

	slog.Debug("input mappings initialized")
}

// 根据按键名称获取鼠标按钮值
func (im *InputManager) mouseButtonFromName(keyName string) sdl.MouseButtonFlags {
	if keyName == "MouseLeft" {
		return sdl.ButtonLeft
	}
	if keyName == "MouseRight" {
		return sdl.ButtonRight
	}
	if keyName == "MouseMiddle" {
		return sdl.ButtonMiddle
	}
	if keyName == "MouseX1" {
		return sdl.ButtonX1
	}
	if keyName == "MouseX2" {
		return sdl.ButtonX2
	}
	return 0
}

// 更新
func (im *InputManager) Update() {
	// 根据上一帧的值更新默认动作状态
	for action, state := range im.actionStates {
		switch state {
		case ActionStatePressedThisFrame:
			im.actionStates[action] = ActionStateHeldDown
		case ActionStateReleasedThisFrame:
			im.actionStates[action] = ActionStateInActive
		}
	}

	// 处理所有等待处理的SDL事件
	var event sdl.Event
	for sdl.PollEvent(&event) {
		im.processEvent(event)
	}
}

// 处理单个SDL事件
func (im *InputManager) processEvent(event sdl.Event) {
	switch event.Type() {
	case sdl.EventKeyDown, sdl.EventKeyUp:
		scancode := event.Key().Scancode
		isDown := event.Key().Down
		isRepeat := event.Key().Repeat

		if actions, exists := im.inputToActionsMap[uint32(scancode)]; exists {
			// 更新action状态
			for _, action := range actions {
				im.updateActionState(action, isDown, isRepeat)
			}
		}
	case sdl.EventMouseButtonDown, sdl.EventMouseButtonUp:
		button := event.Button().Button
		isDown := event.Button().Down

		if actions, exists := im.inputToActionsMap[uint32(button)]; exists {
			// 更新action状态
			for _, action := range actions {
				// 鼠标按钮事件不支持重复触发
				im.updateActionState(action, isDown, false)
			}
		}
	case sdl.EventMouseMotion:
		// 更新鼠标位置
		im.mousePosition[0] = event.Motion().X
		im.mousePosition[1] = event.Motion().Y
	case sdl.EventQuit:
		im.shouldQuit = true
	}
}

// 更新动作状态
func (im *InputManager) updateActionState(action string, isDown bool, isRepeat bool) {
	_, exists := im.actionStates[action]
	if !exists {
		slog.Warn("updateActionState: action not found", slog.String("action", action))
		return
	}

	if isDown {
		if isRepeat {
			// 按键重复按下
			im.actionStates[action] = ActionStateHeldDown
		} else {
			// 按键初次按下
			im.actionStates[action] = ActionStatePressedThisFrame
		}
		return
	}

	// 按键释放
	im.actionStates[action] = ActionStateReleasedThisFrame
}

// 检查动作是否按下
func (im *InputManager) IsActionDown(action string) bool {
	if state, exists := im.actionStates[action]; exists {
		return state == ActionStateHeldDown || state == ActionStatePressedThisFrame
	}
	return false
}

// 检查动作是否在本帧刚刚被按下
func (im *InputManager) IsActionPressed(action string) bool {
	if state, exists := im.actionStates[action]; exists {
		return state == ActionStatePressedThisFrame
	}
	return false
}

// 检查动作是否在本帧刚刚被释放
func (im *InputManager) IsActionReleased(action string) bool {
	if state, exists := im.actionStates[action]; exists {
		return state == ActionStateReleasedThisFrame
	}
	return false
}

// 是否退出
func (im *InputManager) ShouldQuit() bool {
	return im.shouldQuit
}

// 设置退出
func (im *InputManager) SetShouldQuit(shouldQuit bool) {
	im.shouldQuit = shouldQuit
}

// 获取鼠标位置(屏幕坐标)
func (im *InputManager) GetMousePosition() mgl32.Vec2 {
	return im.mousePosition
}

// 获取鼠标位置(逻辑坐标)
func (im *InputManager) GetLogicalMousePosition() mgl32.Vec2 {
	var logicalPos mgl32.Vec2
	sdl.RenderCoordinatesFromWindow(im.sdlRenderer, im.mousePosition[0], im.mousePosition[1], &logicalPos[0], &logicalPos[1])
	return logicalPos
}
