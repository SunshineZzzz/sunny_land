package object

import (
	"log/slog"
	"reflect"

	econtext "sunny_land/src/engine/context"
)

// 组件接口
type IComponent interface {
	// 初始化组件
	Init()
	// 更新组件
	Update(float64, *econtext.Context)
	// 处理输入
	HandleInput(*econtext.Context)
	// 渲染
	Render(*econtext.Context)
	// 清理
	Clean()
	// 设置组件所属的游戏对象
	SetOwner(*GameObject)
	// 获取组件所属的游戏对象
	GetOwner() *GameObject
}

// 游戏对象，负责管理游戏对象的组件
type GameObject struct {
	// 名称
	name string
	// 标签
	tag string
	// 组件列表
	components map[reflect.Type]IComponent
	// 延迟删除标识，将来由场景类负责删除
	needRemove bool
}

// 创建游戏对象
func NewGameObject(name string, tag string) *GameObject {
	slog.Debug("new game object", slog.String("name", name), slog.String("tag", tag))
	return &GameObject{
		name:       name,
		tag:        tag,
		components: make(map[reflect.Type]IComponent),
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
	gt.components = make(map[reflect.Type]IComponent)
}

// 添加组件
func (gt *GameObject) AddComponent(component IComponent) {
	componentType := reflect.TypeOf(component)
	if _, exists := gt.components[componentType]; exists {
		slog.Error("component already exists in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
			slog.String("componentType", componentType.String()))
		return
	}
	gt.components[componentType] = component
	component.SetOwner(gt)
	component.Init()
	slog.Debug("add component to game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.String("componentType", componentType.String()))
}

// 获取组件
func (gt *GameObject) GetComponent(component IComponent) IComponent {
	componentType := reflect.TypeOf(component)
	if component, exists := gt.components[componentType]; exists {
		return component
	}
	slog.Error("component not found in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.String("componentType", componentType.String()))
	return nil
}

// 移除组件
func (gt *GameObject) RemoveComponent(component IComponent) {
	componentType := reflect.TypeOf(component)
	if _, exists := gt.components[componentType]; !exists {
		slog.Error("component not found in game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
			slog.String("componentType", componentType.String()))
		return
	}
	component.Clean()
	delete(gt.components, componentType)
	slog.Debug("remove component from game object", slog.String("gameObject.Name", gt.name), slog.String("gameObject.Tag", gt.tag),
		slog.String("componentType", componentType.String()))
}

// 检查组件是否存在
func (gt *GameObject) HasComponent(component IComponent) bool {
	componentType := reflect.TypeOf(component)
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
