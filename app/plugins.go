package app

import (
	"github.com/RAiWorks/RapidGo/core/plugin"
	exampleplugin "github.com/RAiWorks/RapidGo/plugins/example"
)

// RegisterPlugins registers all application plugins with the manager.
func RegisterPlugins(m *plugin.PluginManager) {
	m.Add(exampleplugin.New())
}
