package digitx

import (
	"math"

	"github.com/spf13/cast"
)

// 获取整数百分比
func Percentage(member interface{}, denominator interface{}) int64 {
	denominatorNum := cast.ToInt64(denominator)
	if denominatorNum <= 0 {
		return 0
	}
	return cast.ToInt64(math.Floor(cast.ToFloat64(member)/cast.ToFloat64(denominator)*100 + 0.5))
}
