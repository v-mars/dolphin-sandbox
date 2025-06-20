package dependencies

import (
	"dolphin-sandbox/biz/pack/runner/types"
	"sync"
)

var preloadScriptMap = map[string]string{}
var preloadScriptMapLock = &sync.RWMutex{}

func SetupDependency(packageName string, version string) {
	preloadScriptMapLock.Lock()
	defer preloadScriptMapLock.Unlock()
	preloadScriptMap[packageName] = version
}

func GetDependency(packageName string, version string) string {
	preloadScriptMapLock.RLock()
	defer preloadScriptMapLock.RUnlock()
	return preloadScriptMap[packageName]
}

func ListDependencies() []types.Dependency {
	var dependencies []types.Dependency
	preloadScriptMapLock.RLock()
	defer preloadScriptMapLock.RUnlock()
	for packageName, version := range preloadScriptMap {
		dependencies = append(dependencies, types.Dependency{
			Name:    packageName,
			Version: version,
		})
	}

	return dependencies
}
