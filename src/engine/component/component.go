package component

import (
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
)

// 基础组件结构体
type Component struct {
	// 组件所属的游戏对象
	owner physics.IGameObject
	// 组件类型
	componentType def.ComponentType
}

// 确保Component实现了IComponent接口
var _ physics.IComponent = (*Component)(nil)

// 初始化组件
func (c *Component) Init() {}

// 更新组件
func (c *Component) Update(float64, physics.IContext) {}

// 处理输入
func (c *Component) HandleInput(physics.IContext) {}

// 渲染
func (c *Component) Render(physics.IContext) {}

// 清理组件
func (c *Component) Clean() {}

// 设置组件所属的游戏对象
func (c *Component) SetOwner(owner physics.IGameObject) {
	c.owner = owner
}

// 获取组件所属的游戏对象
func (c *Component) GetOwner() physics.IGameObject {
	return c.owner
}

// 获取组件类型
func (c *Component) GetType() def.ComponentType {
	return c.componentType
}
