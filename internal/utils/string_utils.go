package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func IsNilOrEmpty(str *string) bool {
	return str == nil || strings.Trim(*str, "") == ""
}

func isCepValid(cep *string) bool {
	regex := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return regex.MatchString(*cep)
}

func CepMask(cep string) (string, error) {
	if isCepValid(&cep) {
		if !strings.Contains(cep, "-") {
			return fmt.Sprintf("%s-%s", cep[:5], cep[5:]), nil
		}
		return cep, nil
	}

	return "", errors.New("cep inválido")
}
