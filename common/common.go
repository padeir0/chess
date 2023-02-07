package common

func Abs32(a int32) int32 {
	y := a >> 31
	return (a ^ y) - y
}

func Abs(a int) int {
	y := int32(a) >> 31
	return int((int32(a) ^ y) - y)
}
