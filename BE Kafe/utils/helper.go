package utils

import "strconv"

// Atoi converts string to int
func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
