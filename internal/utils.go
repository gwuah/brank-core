package internal

import "strconv"

func ConvertToUint64(num string) uint64 {
	id64, _ := strconv.ParseUint(num, 10, 64)
	return id64
}

func ConvertToInt(num string) int {
	id64, _ := strconv.ParseInt(num, 10, 64)
	return int(id64)
}
func Bool(b bool) *bool {
	temp := b
	return &temp
}
