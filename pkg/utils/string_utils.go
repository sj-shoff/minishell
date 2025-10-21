package utils

import "strings"

// TrimSpace обрезает пробелы в строке
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// SplitFields разбивает строку на поля
func SplitFields(s string, sep string) []string {
	return strings.Split(s, sep)
}

// SplitEnvVar разбивает переменную окружения на ключ и значение
func SplitEnvVar(env string) []string {
	return strings.SplitN(env, "=", 2)
}
