package hechodoc

import "strings"

func camelCaseFromStringList(words []string) string{
	for i,v:=range words{
		if i==0{
			continue
		}
		words[i]=strings.Title(v)
	}
	return strings.Join(words,"")
}
