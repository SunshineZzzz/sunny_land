package object

import (
	"log/slog"

	econtext "sunny_land/src/engine/context"
	"sunny_land/src/engine/physics"
	"sunny_land/src/engine/utils/def"
)

// 游戏对象，负责管理游戏对象的组件
type GameObject struct {
	// 名称
	name string
	// 标签
	tag string
	// 组件列表
	components map[def.ComponentType]physics.IComponent
	// 延迟删除标识，将来由场景类负责删除
	needRemove bool
}

// 确保GameObject实现了IGameObject接口
var _ physics.IGameObject = (*GameObject)(nil)

// 创建游戏对象
func NewGameObject(name string, tag string) *GameObject {
	slog.Debug("new game object", slog.String("name", name), slog.String("tag", tag))
	return &GameObject{
		name:       name,
		tag:        tag,
		components: make(map[def.ComponentType]physics.IComponent),
		needRemove: false,
	}
}

// 更新
func (gt *GameObject) Update(deltaTime float64, context *econtext.Context) {
	for _, component := range gt.components {
		component.Update(deltaTime, context)
	}
}

// 渲染
func (gt *GameObject) Render(context *econtext.Context) {
	for _, component := range gt.components {
		component.Render(context)
	}
}

// 处理输入
func (gt *GameObject) HandleInput(context *econtext.Context) {
	for _, component := range gt.components {
		component.HandleInput(context)
	}
}

// 清理
func (gt *GameObject) Clean() {
	for _, component := range gt.components {
		component.Clean()
	}
	gt.components = make(map[def.ComponentType]physics.IComponent)
}

// 添加组件
func (gt *GameObject) AddComponent(component physics.IComponent) physics.IComponent {
	if _, exists := gt.components[component.GetType()]; exists {
		slog.Error("component already exists in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
			slog.Any("componentType", component.GetType()))
		return nil
	}
	gt.components[component.GetType()] = component
	component.SetOwner(gt)
	component.Init()
	slog.Debug("add component to game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.Any("componentType", component.GetType()))
	return component
}

// 获取组件
func (gt *GameObject) GetComponent(componentType def.ComponentType) physics.IComponent {
	if component, exists := gt.components[componentType]; exists {
		return component
	}
	slog.Error("component not found in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.Any("componentType", componentType))
	return nil
}

// 移除组件
func (gt *GameObject) RemoveComponent(componentType def.ComponentType) {
	component, exists := gt.components[componentType]
	if !exists {
		slog.Error("component not found in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
			slog.Any("componentType", componentType))
		return
	}
	component.Clean()
	delete(gt.components, componentType)
	slog.Debug("remove component from game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.Any("componentType", componentType))
}

// 检查组件是否存在
func (gt *GameObject) HasComponent(componentType def.ComponentType) bool {
	_, exists := gt.components[componentType]
	return exists
}

// 设置名称
func (gt *GameObject) SetName(name string) {
	gt.name = name
}

// 获取名称
func (gt *GameObject) GetName() string {
	return gt.name
}

// 设置标签
func (gt *GameObject) SetTag(tag string) {
	gt.tag = tag
}

// 获取标签
func (gt *GameObject) GetTag() string {
	return gt.tag
}

// 设置是否需要删除
func (gt *GameObject) SetNeedRemove(needRemove bool) {
	gt.needRemove = needRemove
}

// 检查是否需要删除
func (gt *GameObject) NeedRemove() bool {
	return gt.needRemove
}
