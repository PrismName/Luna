package plugin

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/traefik/yaegi/interp"
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

type PluginManager struct {
	plugins map[string]VulnPlugin
	mxt     sync.Mutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]VulnPlugin),
	}
}

func (pm *PluginManager) LoadPlugin(path string) error {
	pm.mxt.Lock()
	defer pm.mxt.Unlock()

	i := interp.New(interp.Options{
		// DisableCapabilites: []string{"syscall", "os/exec"},
	})

	i.Use(interp.Symbols)

	code, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = i.Eval(string(code))
	if err != nil {
		return err
	}

	v, err := i.Eval("Plugin")
	if err != nil {
		return fmt.Errorf("plugin Symbol not found")
	}

	plugin, ok := v.Interface().(VulnPlugin)
	if !ok {
		return fmt.Errorf("Invalid plugin type")
	}

	pm.plugins[plugin.Meta().Name] = plugin

	return nil
}

func (pm *PluginManager) ListPlugins() []VulnPlugin {
	pm.mxt.Lock()
	defer pm.mxt.Unlock()

	var list []VulnPlugin
	for _, p := range pm.plugins {
		list = append(list, p)
	}

	return list
}

// GetPlugin 根据名称获取插件
func (pm *PluginManager) GetPlugin(name string) (VulnPlugin, bool) {
	pm.mxt.Lock()
	defer pm.mxt.Unlock()

	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// ExecutePlugin 根据插件名执行插件
func (pm *PluginManager) ExecutePlugin(name string, target string) (bool, error) {
	plugin, exists := pm.GetPlugin(name)
	if !exists {
		return false, fmt.Errorf("插件 '%s' 不存在", name)
	}

	return plugin.Run(target)
}

// SearchPlugins 根据关键字搜索插件
func (pm *PluginManager) SearchPlugins(keyword string) []VulnPlugin {
	pm.mxt.Lock()
	defer pm.mxt.Unlock()

	var results []VulnPlugin
	keyword = strings.ToLower(keyword)

	for _, p := range pm.plugins {
		meta := p.Meta()
		if strings.Contains(strings.ToLower(meta.Name), keyword) ||
			strings.Contains(strings.ToLower(meta.Description), keyword) {
			results = append(results, p)
		}
	}

	return results
}

// UnloadPlugin 卸载指定名称的插件
func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mxt.Lock()
	defer pm.mxt.Unlock()

	if _, exists := pm.plugins[name]; !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	delete(pm.plugins, name)
	return nil
}
