// 示例插件
// 这是一个简单的示例插件，用于测试Luna的插件系统

package main

import (
	"fmt"
	"strings"
	"time"
)

// PluginMeta 定义插件的元数据
type PluginMeta struct {
	Name        string
	Version     string
	Description string
}

// VulnPlugin 接口定义了插件必须实现的方法
type VulnPlugin interface {
	Meta() PluginMeta
	Run(target string) (bool, error)
}

// SamplePlugin 是示例插件的具体实现
type SamplePlugin struct {
	meta PluginMeta
}

// 确保SamplePlugin实现了VulnPlugin接口
var _ VulnPlugin = (*SamplePlugin)(nil)

// 创建插件实例
// 注意：必须命名为Plugin，这是Luna加载插件时查找的符号
var Plugin = &SamplePlugin{
	meta: PluginMeta{
		Name:        "sample_plugin",
		Version:     "1.0.0",
		Description: "Luna示例插件 - 用于测试插件系统",
	},
}

// Meta 返回插件的元数据
func (p *SamplePlugin) Meta() PluginMeta {
	return p.meta
}

// Run 实现插件的主要功能
func (p *SamplePlugin) Run(target string) (bool, error) {
	fmt.Printf("[%s] 示例插件正在运行...\n", time.Now().Format("15:04:05"))

	if target == "" {
		return false, fmt.Errorf("目标不能为空")
	}

	fmt.Printf("分析目标: %s\n", target)

	// 简单的演示逻辑
	time.Sleep(1 * time.Second)
	fmt.Println("正在处理...")
	time.Sleep(1 * time.Second)

	// 检查目标是否包含特定字符串
	if strings.Contains(strings.ToLower(target), "test") {
		fmt.Println("发现测试目标！")
		return true, nil
	} else if strings.Contains(strings.ToLower(target), "fail") {
		fmt.Println("目标处理失败")
		return false, nil
	}

	fmt.Println("目标处理完成")
	return true, nil
}
