# Luna

A PoC framework named Luna，基于 [Yaegi](https://github.com/traefik/yaegi) 实现的动态插件系统，支持插件的注册、加载、执行、搜索和管理。

## 功能特点

- 动态加载 Go 语言编写的插件
- 支持按名称执行插件
- 支持关键字搜索插件
- 提供插件模板，方便开发者创建自己的插件
- 完整的命令行界面，易于使用

## 使用方法

### 安装

```bash
go get github.com/seaung/Luna
```

### 启动 Luna

```bash
go run cmd/lua/luna.go
```

### 插件管理命令

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

## 示例

### 编译并加载示例插件

```bash
# 进入示例目录
cd examples

# 编译示例插件
./build_plugin.sh

# 在 Luna 中加载插件
load /path/to/sample_plugin.so

# 列出已加载的插件
list

# 执行插件
exec sample_plugin test_target
```

## 开发自己的插件

请参考 `templates/README.md` 和 `templates/plugin_template.go` 了解如何开发自己的插件。

### 插件开发快速入门

1. 复制 `templates/plugin_template.go` 作为起点
2. 修改插件元数据（名称、版本、描述）
3. 在 `Run` 方法中实现插件的主要功能
4. 编译插件：`go build -o my_plugin.so -buildmode=plugin my_plugin.go`
5. 在 Luna 中加载并测试插件

## 贡献

欢迎贡献新的插件或改进现有功能。请确保您的代码遵循 Go 语言的代码规范和最佳实践。
