package stringx

import (
	"strings"

	"github.com/spf13/cast"
)

func SplitToInt64(str string, sep string) []int64 {
	tmp := strings.Split(str, sep)
	res := make([]int64, len(tmp))
	for k, x := range tmp {
		res[k] = cast.ToInt64(x)
	}
	return res
}

func SplitToInt(str string, sep string) []int {
	tmp := strings.Split(str, sep)
	res := make([]int, len(tmp))
	for k, x := range tmp {
		res[k] = cast.ToInt(x)
	}
	return res
}

func SplitToInt8(str string, sep string) []int8 {
	tmp := strings.Split(str, sep)
	res := make([]int8, len(tmp))
	for k, x := range tmp {
		res[k] = cast.ToInt8(x)
	}
	return res
}
