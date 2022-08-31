package gutil

func UniqueStrings(values []string) []string {
	if values == nil {
		return nil
	}

	unique := make([]string, 0)
	for _, v := range values {
		if !Contains(unique, v) {
			unique = append(unique, v)
		}
	}

	return unique
}

func RemoveFromStrings(values []string, removal ...string) []string {
	if values == nil {
		return nil
	}

	clean := make([]string, 0)
	for _, v := range values {
		if Contains(removal, v) {
			continue
		}
		clean = append(clean, v)
	}
	return clean
}
