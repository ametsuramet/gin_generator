package utils

import (
	"strings"
)

type StringConv struct {
}

func (u *StringConv) ToCamel(input string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range input {
		if k == 0 {
			camelCase = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}
