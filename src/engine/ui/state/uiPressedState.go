package state

import econtext "sunny_land/src/engine/context"

/**
 * @brief 按下状态
 *
 * 当鼠标按下UI元素时，会切换到该状态。
 */
type UIPressedState struct {
	// 继承基础UI状态
	UIState
}

// 确保UIPressedState实现IUIState接口
var _ IUIState = (*UIPressedState)(nil)

// 创建按下状态实例
func NewUIPressedState(owner IUIInteractive) *UIPressedState {
	return &UIPressedState{
		UIState: UIState{
			owner: owner,
		},
	}
}

// 进入状态
func (p *UIPressedState) Enter() {
	p.owner.SetSprite("pressed")
	p.owner.PlaySound("pressed")
}

// 处理输入
func (p *UIPressedState) HandleInput(ctx *econtext.Context) IUIState {
	inputManager := ctx.GetInputManager()
	mousePos := inputManager.GetLogicalMousePosition()
	// 如果鼠标不在UI元素内，则切换到正常状态
	if !p.owner.IsPointInside(mousePos) {
		return NewUINormalState(p.owner)
	}
	// 如果鼠标松开，则返回悬停状态
	if inputManager.IsActionReleased("MouseLeftClick") {
		p.owner.Clicked()
		return NewUIHoverState(p.owner)
	}
	return nil
}
