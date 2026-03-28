package sdk

import "sync"

var (
	registryMu sync.Mutex
	registry   []Plugin
)

// Register adds a plugin to the global registry.
// Call this from your plugin's init() function.
func Register(p Plugin) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, p)
}

// All returns all registered plugins.
func All() []Plugin {
	registryMu.Lock()
	defer registryMu.Unlock()
	result := make([]Plugin, len(registry))
	copy(result, registry)
	return result
}
