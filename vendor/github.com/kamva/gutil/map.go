package gutil

import (
	"errors"
	"fmt"
)

// KeyValuesToMap convert each pair of slice to a key/value of map.
func KeyValuesToMap(keyValues ...interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			m[fmt.Sprint(keyValues[i])] = keyValues[i+1]
		} else {
			return nil, errors.New(fmt.Sprintf("missing %v index value", i))
		}
	}

	return m, nil
}

// MapToKeyValue convert map to slice that contain map key and value pairs.
func MapToKeyValue(m map[string]interface{}) []interface{} {
	kv := make([]interface{}, len(m)*2)

	i := 0
	for k, v := range m {
		kv[i] = k
		kv[i+1] = v
		i += 2
	}

	return kv
}

// ExtendMap extend dest map by src map.
// overwrite argument set overwrite policy.
func ExtendMap(dest, src map[string]interface{}, overwrite bool) {
	if dest == nil {
		panic("destination map is nil, can not extend nil map")
	}

	for key, val := range src {
		// If key exists in dest and we can not overwrite it, so continue.
		if _, ok := dest[key]; ok && !overwrite {
			continue
		}
		dest[key] = val
	}
}

// ExtendMap extend dest map by src map.
// overwrite argument set overwrite policy.
func ExtendStrMap(dest, src map[string]string, overwrite bool) {
	if dest == nil {
		panic("destination map is nil, can not extend nil map")
	}

	for key, val := range src {
		// If key exists in dest and we can not overwrite it, so continue.
		if _, ok := dest[key]; ok && !overwrite {
			continue
		}
		dest[key] = val
	}
}

// MapPathExtractor extract the map paths recursively.
type MapPathExtractor struct {
	// Depth is depth of the keys. set to 1 to get just first level of keys.
	Depth int

	// Separator defines separator string between keys on each depth.
	// e.g if we set "." separator set to "." so we will have a.b for
	// {a:{b:"hello"}} map.
	Separator string
}

func (e MapPathExtractor) Extract(m map[string]interface{}) []string {
	return e.extractRecursive(m, "", 1)
}

func (e MapPathExtractor) extractRecursive(m map[string]interface{}, prefix string, currentDepth int) []string {
	if currentDepth > e.Depth {
		return nil
	}
	var paths []string
	for k, v := range m {
		path := e.generatePath(prefix, k)
		paths = append(paths, path)
		if newMap, ok := v.(map[string]interface{}); ok {
			paths = append(paths, e.extractRecursive(newMap, path, currentDepth+1)...)
		}
	}
	return paths
}

func (e MapPathExtractor) generatePath(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return fmt.Sprintf("%s%s%s", prefix, e.Separator, key)
}