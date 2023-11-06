package helpers

func MapKeys(m map[string]string) []string {
	var res = make([]string, len(m))
	var i int
	for k, _ := range m {
		res[i] = k
		i++
	}
	return res
}
