# Luna 插件开发指南

## 插件系统概述

Luna 的插件系统基于 [Yaegi](https://github.com/traefik/yaegi) 实现，允许开发者创建自定义插件来扩展 Luna 的功能。插件系统支持动态加载、执行、搜索和卸载插件。

## 插件开发

### 插件结构

每个 Luna 插件必须实现 `VulnPlugin` 接口，该接口定义了以下方法：

```go
type VulnPlugin interface {
	Meta() PluginMeta
	Run(target string) (bool, error)
}
```

其中 `PluginMeta` 结构体包含插件的基本信息：

```go
type PluginMeta struct {
	Name        string // 插件名称
	Version     string // 插件版本
	Description string // 插件描述
}
```

### 创建新插件

1. 复制 `templates/plugin_template.go` 作为起点
2. 修改插件元数据（名称、版本、描述）
3. 在 `Run` 方法中实现插件的主要功能
4. 确保导出一个名为 `Plugin` 的变量，Luna 将通过此变量加载插件

### 示例插件

```go
package main

import (
	"fmt"
)

type PluginMeta struct {
	Name        string
	Version     string
	Description string
}

type VulnPlugin interface {
	Meta() PluginMeta
	Run(target string) (bool, error)
}

type MyPlugin struct {
	meta PluginMeta
}

var Plugin = &MyPlugin{
	meta: PluginMeta{
		Name:        "my_plugin",
		Version:     "1.0.0",
		Description: "这是一个示例插件",
	},
}

func (p *MyPlugin) Meta() PluginMeta {
	return p.meta
}

func (p *MyPlugin) Run(target string) (bool, error) {
	fmt.Printf("插件 %s 正在运行，目标: %s\n", p.meta.Name, target)
	
	if target == "" {
		return false, fmt.Errorf("目标不能为空")
	}
	
	// 在这里实现你的插件逻辑
	
	return true, nil
}
```

## 编译和使用插件

### 编译插件

```bash
go build -o my_plugin.so -buildmode=plugin my_plugin.go
```

### 在 Luna 中使用插件

1. 加载插件：`load /path/to/my_plugin.so`
2. 列出已加载的插件：`list`
3. 搜索插件：`search <关键字>`
4. 使用插件：
   - 方法1：先选择插件 `use <插件名>` 然后运行 `run`
   - 方法2：直接执行 `exec <插件名> [目标]`
5. 卸载插件：`unload <插件名>`

## 插件命令参考

| 命令 | 描述 | 用法 |
|------|------|------|
| `load` | 加载插件 | `load <plugin_path>` |
| `list` | 列出所有已加载的插件 | `list` |
| `search` | 搜索插件 | `search <keyword>` |
| `use` | 选择要使用的插件 | `use <plugin_name>` |
| `run` | 运行当前选择的插件 | `run` |
| `exec` | 执行指定名称的插件 | `exec <plugin_name> [target]` |
| `unload` | 卸载指定名称的插件 | `unload <plugin_name>` |
| `set` | 设置参数值 | `set <option> <value>` |
| `unset` | 清除参数值 | `unset <option>` |

## 最佳实践

1. 为插件提供详细的描述，便于用户理解插件功能
2. 实现适当的错误处理和日志输出
3. 在插件中添加帮助信息或使用示例
4. 遵循 Go 语言的代码规范和最佳实践
5. 为复杂插件添加配置选项

## 贡献

欢迎贡献新的插件或改进现有插件。请确保您的插件遵循上述指南，并提供充分的文档。