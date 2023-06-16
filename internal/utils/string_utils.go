package utils

import "strings"

func IsNilOrEmpty(str *string) bool {
	return str == nil || strings.Trim(*str, "") == ""
}
