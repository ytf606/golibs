package timex

import (
	"database/sql/driver"
	"time"
)

type Time time.Time

const (
	DefaultFormat    = "2006-01-02 15:04:05"
	DefaultYMDFormat = "2006-01-02"
	DefaultYMFormat  = "200601"
)

var (
	NewTicker = time.NewTicker
)

type (
	Duration = time.Duration
	Ticker   = time.Ticker
)

func SetTimeZone() {
	time.LoadLocation("Asia/Shanghai")
}

func Now() time.Time {
	return time.Now()
}

func TableSuffix() string {
	return Now().Format(DefaultYMFormat)
}

// BeforeMonthToString 获取当前时间往前N个月的格式化时间
func BeforeMonthToString(mon int) string {
	return BeforeMonth(mon).Format(DefaultFormat)
}

// BeforeMonth 获取当前时间往前N个月的time
func BeforeMonth(mon int) time.Time {
	return time.Now().AddDate(0, -mon, 0)
}

func AddDateToString(year, month, day int, startTime string) string {
	return AddDate(year, month, day, startTime).Format(DefaultFormat)
}

func AddDate(year, month, day int, startTime string) time.Time {
	t := time.Now()
	if startTime != "" {
		t, _ = DateStringToTime(startTime)
	}
	return t.AddDate(year, month, day)
}

// GetNowMonthName 获取当前时间的月份名称
func GetNowMonthName() string {
	var monMap = map[time.Month]string{
		1: "一月", 2: "二月", 3: "三月", 4: "四月", 5: "五月", 6: "六月",
		7: "七月", 8: "八月", 9: "九月", 10: "十月", 11: "十一月", 12: "十二月",
	}

	nowMon := time.Now().Month()

	return monMap[nowMon]
}

// BeginOfMonthToString 获取本月第一天0点的格式化时间
func BeginOfMonthToString() string {
	return BeginOfMonth().Format(DefaultYMDFormat)
}

// BeginOfMonth 获取本月第一天0点
func BeginOfMonth() time.Time {
	nowTime := time.Now()
	return nowTime.AddDate(0, 0, -nowTime.Day()+1) // 本月第一天0点
}

func Sleep(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

func SleepMilli(milli int) {
	time.Sleep(time.Duration(milli) * time.Millisecond)
}

func DurationSecond(second int) time.Duration {
	return time.Duration(second) * time.Second
}

func UnixSecond(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0)
}

func DateStringToTime(datetime string) (time.Time, error) {
	lay := DefaultFormat
	if len(datetime) == 10 {
		lay = DefaultYMDFormat
	}
	return time.ParseInLocation(lay, datetime, time.Local)
}

func DateUnix(datetime string) (int64, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}

	date, err := time.ParseInLocation(DefaultFormat, datetime, loc)
	if err != nil {
		return 0, err
	}

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return date.Unix(), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+DefaultFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DefaultFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, DefaultFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(DefaultFormat)
}

func (t Time) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Now()
	if time.Time(t).IsZero() != true {
		tTime = time.Time(t)
	}
	return tTime.Format(DefaultFormat), nil
}
