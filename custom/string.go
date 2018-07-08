package custom

import (
	"regexp"
	"strings"
)

func String2key(str string) string {
	re1, _ := regexp.Compile("[^a-z0-9]")
	re2, _ := regexp.Compile("--*")
	re3, _ := regexp.Compile("-*$")

	result1 := string(re1.ReplaceAll([]byte(strings.ToLower(str)), []byte("-")))
	result2 := string(re2.ReplaceAll([]byte(result1), []byte("-")))
	return string(re3.ReplaceAll([]byte(result2), []byte("")))
}

func String2InCaseSensitivePattern(s string) string{
	var sl, su, slu, pattern string
	var patterns = []string{}
	
	for i := 0; i < len(s); i++ {
		sl = strings.ToLower(string(s[i]))
		su = strings.ToUpper(sl)
		slu = strings.Join([]string{"[", su, sl, "]"}, "")
		patterns = append(patterns, slu)
	}
	pattern = strings.Join(patterns, "")
	return pattern	
}