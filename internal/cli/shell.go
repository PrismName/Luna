package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/seaung/Luna/internal/plugin"
)

// Command 表示一个CLI命令
type Command struct {
	Name        string
	Description string
	Usage       string
	Action      func(args []string) error
}

// CommandContext 保存命令执行的上下文信息
type CommandContext struct {
	Target     string
	PluginName string
	Options    map[string]string
}

// Shell 表示交互式命令行界面
type Shell struct {
	Commands       map[string]Command
	PluginMgr      *plugin.PluginManager
	Context        CommandContext
	Prompt         string
	History        []string
	HistoryMaxSize int
}

// NewShell 创建一个新的Shell实例
func NewShell() *Shell {
	return &Shell{
		Commands:       make(map[string]Command),
		PluginMgr:      plugin.NewPluginManager(),
		Prompt:         "luna > ",
		History:        make([]string, 0),
		HistoryMaxSize: 100,
		Context: CommandContext{
			Options: make(map[string]string),
		},
	}
}

// RegisterCommand 注册一个命令
func (s *Shell) RegisterCommand(cmd Command) {
	s.Commands[cmd.Name] = cmd
}

// AddToHistory 添加命令到历史记录
func (s *Shell) AddToHistory(cmdLine string) {
	if len(s.History) >= s.HistoryMaxSize {
		// 移除最旧的历史记录
		s.History = s.History[1:]
	}
	s.History = append(s.History, cmdLine)
}

// contains 检查字符串是否存在于切片中
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// setupCommands 设置所有可用命令
func (s *Shell) setupCommands() {
	s.RegisterCommand(Command{
		Name:        "help",
		Description: "显示帮助信息",
		Usage:       "help [command]",
		Action:      s.cmdHelp,
	})

	s.RegisterCommand(Command{
		Name:        "exit",
		Description: "退出程序",
		Usage:       "exit",
		Action:      s.cmdExit,
	})

	s.RegisterCommand(Command{
		Name:        "load",
		Description: "加载插件",
		Usage:       "load <plugin_path>",
		Action:      s.cmdLoadPlugin,
	})

	s.RegisterCommand(Command{
		Name:        "list",
		Description: "列出所有已加载的插件",
		Usage:       "list",
		Action:      s.cmdListPlugins,
	})

	s.RegisterCommand(Command{
		Name:        "search",
		Description: "搜索插件",
		Usage:       "search <keyword>",
		Action:      s.cmdSearchPlugins,
	})

	s.RegisterCommand(Command{
		Name:        "exec",
		Description: "执行指定名称的插件",
		Usage:       "exec <plugin_name> [target]",
		Action:      s.cmdExecPlugin,
	})

	s.RegisterCommand(Command{
		Name:        "unload",
		Description: "卸载指定名称的插件",
		Usage:       "unload <plugin_name>",
		Action:      s.cmdUnloadPlugin,
	})

	s.RegisterCommand(Command{
		Name:        "set",
		Description: "设置参数值",
		Usage:       "set <option> <value>",
		Action:      s.cmdSetOption,
	})

	s.RegisterCommand(Command{
		Name:        "unset",
		Description: "清除参数值",
		Usage:       "unset <option>",
		Action:      s.cmdUnsetOption,
	})

	s.RegisterCommand(Command{
		Name:        "use",
		Description: "选择要使用的插件",
		Usage:       "use <plugin_name>",
		Action:      s.cmdUsePlugin,
	})

	s.RegisterCommand(Command{
		Name:        "run",
		Description: "运行当前选择的插件",
		Usage:       "run",
		Action:      s.cmdRunPlugin,
	})

	s.RegisterCommand(Command{
		Name:        "show",
		Description: "显示信息",
		Usage:       "show [options|plugins]",
		Action:      s.cmdShow,
	})

	s.RegisterCommand(Command{
		Name:        "history",
		Description: "显示命令历史",
		Usage:       "history",
		Action:      s.cmdHistory,
	})
}

// cmdHelp 显示帮助信息
func (s *Shell) cmdHelp(args []string) error {
	if len(args) > 0 {
		// 显示特定命令的帮助
		cmdName := args[0]
		cmd, exists := s.Commands[cmdName]
		if !exists {
			return fmt.Errorf("未知命令: %s", cmdName)
		}

		fmt.Printf("命令: %s\n", cmd.Name)
		fmt.Printf("描述: %s\n", cmd.Description)
		fmt.Printf("用法: %s\n", cmd.Usage)
		return nil
	}

	// 显示所有命令
	fmt.Println("可用命令:")
	fmt.Println("==========")

	for name, cmd := range s.Commands {
		fmt.Printf("%-10s - %s\n", name, cmd.Description)
	}
	fmt.Println("\n使用 'help <command>' 获取特定命令的详细信息")
	return nil
}

// cmdExit 退出程序
func (s *Shell) cmdExit(args []string) error {
	fmt.Println("再见!")
	os.Exit(0)
	return nil
}

// cmdLoadPlugin 加载插件
func (s *Shell) cmdLoadPlugin(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["load"].Usage)
	}

	path := args[0]
	err := s.PluginMgr.LoadPlugin(path)
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}

	fmt.Printf("插件 '%s' 加载成功\n", path)
	return nil
}

// cmdListPlugins 列出所有已加载的插件
func (s *Shell) cmdListPlugins(args []string) error {
	plugins := s.PluginMgr.ListPlugins()

	if len(plugins) == 0 {
		fmt.Println("没有加载任何插件")
		return nil
	}

	fmt.Println("已加载的插件:")
	fmt.Println("=============")

	for _, p := range plugins {
		meta := p.Meta()
		fmt.Printf("%-20s - %s (v%s)\n", meta.Name, meta.Description, meta.Version)
	}

	return nil
}

// cmdSearchPlugins 搜索插件
func (s *Shell) cmdSearchPlugins(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["search"].Usage)
	}

	keyword := args[0]
	plugins := s.PluginMgr.SearchPlugins(keyword)

	if len(plugins) == 0 {
		fmt.Printf("没有找到包含关键字 '%s' 的插件\n", keyword)
		return nil
	}

	fmt.Printf("搜索结果 ('%s'):\n", keyword)
	fmt.Println("====================")

	for _, p := range plugins {
		meta := p.Meta()
		fmt.Printf("%-20s - %s (v%s)\n", meta.Name, meta.Description, meta.Version)
	}

	return nil
}

// cmdExecPlugin 执行指定名称的插件
func (s *Shell) cmdExecPlugin(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["exec"].Usage)
	}

	pluginName := args[0]
	target := ""

	if len(args) > 1 {
		target = args[1]
	} else if s.Context.Target != "" {
		target = s.Context.Target
	}

	if target == "" {
		fmt.Println("警告: 未指定目标，可以使用 'set target <value>' 设置默认目标")
	}

	fmt.Printf("执行插件 '%s'...\n", pluginName)
	success, err := s.PluginMgr.ExecutePlugin(pluginName, target)

	if err != nil {
		return fmt.Errorf("执行插件失败: %v", err)
	}

	if success {
		fmt.Println("插件执行成功")
	} else {
		fmt.Println("插件执行失败")
	}

	return nil
}

// cmdUnloadPlugin 卸载指定名称的插件
func (s *Shell) cmdUnloadPlugin(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["unload"].Usage)
	}

	pluginName := args[0]
	err := s.PluginMgr.UnloadPlugin(pluginName)

	if err != nil {
		return fmt.Errorf("卸载插件失败: %v", err)
	}

	fmt.Printf("插件 '%s' 已卸载\n", pluginName)
	return nil
}

// cmdSetOption 设置选项值
func (s *Shell) cmdSetOption(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: %s", s.Commands["set"].Usage)
	}

	option := args[0]
	value := args[1]

	// 特殊处理target选项
	if option == "target" {
		s.Context.Target = value
	}

	// 保存到选项映射
	s.Context.Options[option] = value
	fmt.Printf("%s => %s\n", option, value)
	return nil
}

// cmdUnsetOption 清除选项值
func (s *Shell) cmdUnsetOption(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["unset"].Usage)
	}

	option := args[0]

	// 特殊处理target选项
	if option == "target" {
		s.Context.Target = ""
	}

	// 从选项映射中删除
	delete(s.Context.Options, option)
	fmt.Printf("%s 已清除\n", option)
	return nil
}

// cmdUsePlugin 选择要使用的插件
func (s *Shell) cmdUsePlugin(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["use"].Usage)
	}

	pluginName := args[0]
	plugins := s.PluginMgr.ListPlugins()

	// 验证插件是否存在
	var found bool
	for _, p := range plugins {
		if p.Meta().Name == pluginName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("找不到插件: %s", pluginName)
	}

	s.Context.PluginName = pluginName
	// 更新提示符以显示当前插件
	s.Prompt = fmt.Sprintf("luna (%s) > ", pluginName)
	fmt.Printf("使用插件: %s\n", pluginName)
	return nil
}

// cmdRunPlugin 运行当前选择的插件
func (s *Shell) cmdRunPlugin(args []string) error {
	if s.Context.PluginName == "" {
		return fmt.Errorf("请先使用 'use <plugin_name>' 选择一个插件")
	}

	if s.Context.Target == "" {
		return fmt.Errorf("请先使用 'set target <target_value>' 设置目标")
	}

	plugins := s.PluginMgr.ListPlugins()
	var targetPlugin plugin.VulnPlugin

	for _, p := range plugins {
		if p.Meta().Name == s.Context.PluginName {
			targetPlugin = p
			break
		}
	}

	if targetPlugin == nil {
		return fmt.Errorf("找不到插件: %s", s.Context.PluginName)
	}

	fmt.Printf("正在运行插件 '%s' 检测目标 '%s'...\n", s.Context.PluginName, s.Context.Target)

	vuln, err := targetPlugin.Run(s.Context.Target)
	if err != nil {
		return fmt.Errorf("插件运行失败: %v", err)
	}

	if vuln {
		fmt.Printf("[!] 目标 '%s' 存在漏洞!\n", s.Context.Target)
	} else {
		fmt.Printf("[+] 目标 '%s' 安全\n", s.Context.Target)
	}

	return nil
}

// cmdShow 显示信息
func (s *Shell) cmdShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: %s", s.Commands["show"].Usage)
	}

	switch args[0] {
	case "options":
		fmt.Println("当前设置:")
		fmt.Println("=========")
		fmt.Printf("当前插件: %s\n", s.Context.PluginName)
		fmt.Printf("目标: %s\n", s.Context.Target)

		// 显示其他选项
		for k, v := range s.Context.Options {
			if k != "target" { // 已经单独显示了
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	case "plugins":
		return s.cmdListPlugins(nil)
	default:
		return fmt.Errorf("未知的show子命令: %s", args[0])
	}

	return nil
}

// cmdHistory 显示命令历史
func (s *Shell) cmdHistory(args []string) error {
	if len(s.History) == 0 {
		fmt.Println("没有命令历史")
		return nil
	}

	fmt.Println("命令历史:")
	fmt.Println("=========")

	for i, cmd := range s.History {
		fmt.Printf("%3d  %s\n", i+1, cmd)
	}

	return nil
}

// Run 启动交互式shell
func (s *Shell) Run() {
	s.setupCommands()

	fmt.Println("Luna 漏洞扫描框架 - 交互式控制台")
	fmt.Println("输入 'help' 获取可用命令列表")

	for {
		// 准备自定义输入
		fmt.Print(s.Prompt)

		// 创建自定义输入选项
		suggestions := make([]string, 0, len(s.Commands))

		// 添加历史记录到建议列表
		for _, cmd := range s.History {
			suggestions = append(suggestions, cmd)
		}

		// 添加命令到建议列表
		for name := range s.Commands {
			// 避免重复添加已在历史记录中的命令
			if !contains(suggestions, name) {
				suggestions = append(suggestions, name)
			}
		}

		// 添加插件名称到建议列表中
		for _, p := range s.PluginMgr.ListPlugins() {
			cmdWithPlugin := "use " + p.Meta().Name
			if !contains(suggestions, cmdWithPlugin) {
				suggestions = append(suggestions, cmdWithPlugin)
			}
		}

		// 如果没有建议，添加一个默认选项
		if len(suggestions) == 0 {
			suggestions = append(suggestions, "help")
		}

		// 使用promptui.Select创建交互式选择
		select_prompt := promptui.Select{
			Label: "命令",
			Items: suggestions,
			Size:  10, // 显示10个选项
			Searcher: func(input string, index int) bool {
				suggestion := suggestions[index]
				return strings.Contains(strings.ToLower(suggestion), strings.ToLower(input))
			},
			StartInSearchMode: true, // 直接进入搜索模式
		}

		_, cmdLine, err := select_prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
				s.cmdExit(nil)
				return
			}
			fmt.Printf("提示错误: %v\n", err)
			continue
		}

		// 跳过空命令
		cmdLine = strings.TrimSpace(cmdLine)
		if cmdLine == "" {
			continue
		}

		// 添加到历史记录
		s.AddToHistory(cmdLine)

		// 解析命令和参数
		parts := strings.Fields(cmdLine)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		// 执行命令
		cmd, exists := s.Commands[cmdName]
		if !exists {
			fmt.Printf("未知命令: %s\n", cmdName)
			fmt.Println("输入 'help' 获取可用命令列表")
			continue
		}

		err = cmd.Action(args)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
		}
	}
}
