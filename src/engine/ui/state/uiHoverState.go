package state

import econtext "sunny_land/src/engine/context"

/**
 * @brief 悬停状态
 *
 * 当鼠标悬停在UI元素上时，会切换到该状态。
 */
type UIHoverState struct {
	// 继承基础UI状态
	UIState
}

// 确保UIHoverState实现IUIState接口
var _ IUIState = (*UIHoverState)(nil)

// 创建悬停状态实例
func NewUIHoverState(owner IUIInteractive) *UIHoverState {
	return &UIHoverState{
		UIState: UIState{
			owner: owner,
		},
	}
}

// 进入状态
func (h *UIHoverState) Enter() {
	h.owner.SetSprite("hover")
}

// 处理输入
func (h *UIHoverState) HandleInput(ctx *econtext.Context) IUIState {
	inputManager := ctx.GetInputManager()
	mousePos := inputManager.GetLogicalMousePosition()
	// 如果鼠标不在UI元素内，则切换到正常状态
	if !h.owner.IsPointInside(mousePos) {
		return NewUINormalState(h.owner)
	}
	// 如果鼠标按下，则返回按下状态
	if inputManager.IsActionPressed("MouseLeftClick") {
		return NewUIPressedState(h.owner)
	}
	return nil
}
