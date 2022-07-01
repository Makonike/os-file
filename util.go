package main

func IsBlank(str string) bool {
	return str == ""
}

func DivWithUp(div, divisor int) int {
	if div%divisor == 0 {
		return div / divisor
	} else {
		return div/divisor + 1
	}
}
