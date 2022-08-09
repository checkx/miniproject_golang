package helpers

import "strconv"

// Create function convert string to uint and return uint
func ConvertStringToUint(str string) uint {
	var uintVal uint
	result, _ := strconv.ParseUint(str, 10, 64)
	uintVal = uint(result)
	return uintVal
}

// Create function conver uint to string and return string
func ConvertUintToString(uintVal uint) string {
	var str string
	str = strconv.FormatUint(uint64(uintVal), 10)
	return str
}
