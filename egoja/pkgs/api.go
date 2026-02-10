package pkgs

var cacheImport = make(map[string]map[string]any)

func GetCache(name string) (map[string]any, bool) {
	m, ok := cacheImport[name]
	return m, ok
}

func RegisterImport(name string, value map[string]any) {
	m, ok := cacheImport[name]
	if !ok {
		cacheImport[name] = value
		return
	}
	for k, v := range value {
		m[k] = v
	}
	cacheImport[name] = m

}

func DeleteImport(name string) {
	delete(cacheImport, name)

}
