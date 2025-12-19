package component

import (
	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/object"
)

// 基础组件结构体
type Component struct {
	// 组件所属的游戏对象
	owner *object.GameObject
}

// 确保Component实现了IComponent接口
var _ object.IComponent = (*Component)(nil)

// 初始化组件
func (c *Component) Init() {}

// 更新组件
func (c *Component) Update(float64, *econtext.Context) {}

// 处理输入
func (c *Component) HandleEvents(context *econtext.Context) {}

// 渲染
func (c *Component) Render(context *econtext.Context) {}

// 清理组件
func (c *Component) Clean() {}

// 设置组件所属的游戏对象
func (c *Component) SetOwner(owner *object.GameObject) {
	c.owner = owner
}

// 获取组件所属的游戏对象
func (c *Component) GetOwner() *object.GameObject {
	return c.owner
}

// 非接口实现方法
