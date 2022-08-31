package hexafaktory

import (
	"fmt"

	"github.com/kamva/tracer"
)

func bytesMapToInterfaceMap(m map[string][]byte) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		res[k] = v
	}
	return res
}

func interfaceMapToBytesMap(m map[string]interface{}) (map[string][]byte, error) {
	res := make(map[string][]byte)
	for k, v := range m {
		b, ok := v.([]byte)
		if !ok {
			return nil, tracer.Trace(fmt.Errorf("key %s values's type is invalid", k))
		}
		res[k] = b
	}
	return res, nil
}
