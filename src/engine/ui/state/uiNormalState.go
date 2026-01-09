package state

import (
	econtext "sunny_land/src/engine/context"
)

/**
 * @brief 正常状态
 *
 * 正常状态是UI元素的默认状态。
 */
type UINormalState struct {
	// 继承基础UI状态
	UIState
}

// 确保UINormalState实现UIState接口
var _ IUIState = (*UINormalState)(nil)

// 创建正常状态实例
func NewUINormalState(owner IUIInteractive) *UINormalState {
	return &UINormalState{
		UIState: UIState{
			owner: owner,
		},
	}
}

// 进入状态
func (n *UINormalState) Enter() {
	n.owner.SetSprite("normal")
}

// 处理输入
func (n *UINormalState) HandleInput(ctx *econtext.Context) IUIState {
	inputManager := ctx.GetInputManager()
	mousePos := inputManager.GetLogicalMousePosition()
	// 如果鼠标在UI元素内，则切换到悬停状态
	if n.owner.IsPointInside(mousePos) {
		n.owner.PlaySound("hover")
		return NewUIHoverState(n.owner)
	}
	return nil
}
