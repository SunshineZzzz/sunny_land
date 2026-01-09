package ui

import (
	"log/slog"
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/render"
	"sunny_land/src/engine/ui/state"

	"github.com/go-gl/mathgl/mgl32"
)

/**
 * @brief 可交互UI元素的基类,继承自UIElement
 *
 * 定义了可交互UI元素的通用属性和行为。
 * 管理UI状态的切换和交互逻辑。
 * 提供事件处理、更新和渲染的虚方法。
 */
type UIInteractive struct {
	// 继承UI元素基础实现
	UIElement
	// 可交互元素很可能需要其他引擎组件
	context *econtext.Context
	// 当前状态
	state state.IUIState
	// 精灵图名称<->精灵图映射
	spriteMap map[string]*render.Sprite
	// 音效名称<->音效文件路径
	soundMap map[string]string
	// 当前显示的精灵图对象
	currentSprite *render.Sprite
	// 是否可以交互
	interactive bool
}

// 确保UIInteractive实现了IUIInteractive接口
var _ state.IUIInteractive = (*UIInteractive)(nil)

// 确保UIInteractive实现了IUIElement接口
var _ state.IUIElement = (*UIInteractive)(nil)

// 构建可交互UI元素
func BuildUIInteractive(uii *UIInteractive, context *econtext.Context, position mgl32.Vec2, size mgl32.Vec2) *UIInteractive {
	uii.context = context
	uii.spriteMap = make(map[string]*render.Sprite)
	uii.soundMap = make(map[string]string)
	uii.interactive = true
	BuildUIElement(&uii.UIElement, position, size)
	return uii
}

// 添加精灵图对象
func (uii *UIInteractive) AddSprite(name string, sprite *render.Sprite) {
	// 可交互UI元素必须有一个size用于交互检测，因此如果参数列表中没有指定，则用图片大小作为size
	if uii.size.X() <= 0.0 && uii.size.Y() <= 0.0 {
		uii.size = uii.context.GetResourceManager().GetTextureSize(sprite.GetTextureId())
	}
	uii.spriteMap[name] = sprite
}

// 设置当前显示的精灵
func (uii *UIInteractive) SetSprite(name string) {
	if sprite, ok := uii.spriteMap[name]; ok {
		uii.currentSprite = sprite
		return
	}

	slog.Warn("sprite not in map", slog.String("name", name))
}

// 添加音效
func (uii *UIInteractive) AddSound(name string, soundPath string) {
	uii.soundMap[name] = soundPath
}

// 播放音效
func (uii *UIInteractive) PlaySound(name string) {
	if soundPath, ok := uii.soundMap[name]; ok {
		uii.context.GetAudioPlayer().PlaySound(soundPath)
		return
	}

	slog.Warn("sound not in map", slog.String("name", name))
}

// 设置当前状态

func (uii *UIInteractive) SetState(state state.IUIState) {
	if uii.state == nil {
		slog.Warn("state is nil")
	}

	uii.state = state
	uii.state.Enter()
}

// 获取当前状态
func (uii *UIInteractive) GetState() state.IUIState {
	return uii.state
}

// 设置是否可交互
func (uii *UIInteractive) SetInteractive(interactive bool) {
	uii.interactive = interactive
}

// 获取是否可交互
func (uii *UIInteractive) GetInteractive() bool {
	return uii.interactive
}

// 处理输入事件
func (uii *UIInteractive) HandleInput(ctx *econtext.Context) bool {
	if uii.UIElement.HandleInput(ctx) {
		return true
	}

	// 先更新子节点，再更新自己（状态）
	if uii.state != nil && uii.interactive {
		if nextState := uii.state.HandleInput(ctx); nextState != nil {
			uii.SetState(nextState)
			return true
		}
	}
	return false
}

// 渲染可交互UI元素
func (uii *UIInteractive) Render(ctx *econtext.Context) {
	if !uii.visible {
		return
	}

	// 先渲染自身
	ctx.GetRenderer().DrawUISprite(uii.currentSprite, uii.GetScreenPosition(), &uii.size)

	// 再渲染子元素（调用基类方法）
	uii.UIElement.Render(ctx)
}

// 如果有点击事件，则重写该方法
func (uii *UIInteractive) Clicked() {
}
