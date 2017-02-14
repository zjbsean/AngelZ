package util

import "time"

//时间格式所设定的时间必须是：2006-01-02 15:04:05
const (
	//DateFormat : 年-月-日
	DateFormat = "2006-01-02"
	//DateTimeHFormat : 年-月-日 时
	DateTimeHFormat = "2006-01-02 15"
	//DateTimeHMFormat : 年-月-日 时:分
	DateTimeHMFormat = "2006-01-02 15:04"
	//DateTimeHMSFormat : 年-月-日 时:分:秒
	DateTimeHMSFormat = "2006-01-02 15:04:05"
	//DateTimeHMSMFormat : 年-月-日 时:分:秒.毫秒
	DateTimeHMSMFormat = "2006-01-02 15:04:05.000"
)

//IsSameYear : 是否同一年
func IsSameYear(t1 *time.Time, t2 *time.Time) bool {
	if t1.Year() == t2.Year() {
		return true
	}
	return false
}

//IsSameYearUnix : 是否同一年
func IsSameYearUnix(ut1 int64, ut2 int64) bool {
	t1 := time.Unix(ut1, 0)
	t2 := time.Unix(ut2, 0)
	return IsSameYear(&t1, &t2)
}

//IsSameMonth : 是否同一月
func IsSameMonth(t1 *time.Time, t2 *time.Time) bool {
	if t1.Month() == t2.Month() && IsSameYear(t1, t2) {
		return true
	}
	return false
}

//IsSameMonthUnix : 是否同一月
func IsSameMonthUnix(ut1 int64, ut2 int64) bool {
	t1 := time.Unix(ut1, 0)
	t2 := time.Unix(ut2, 0)
	return IsSameMonth(&t1, &t2)
}

//IsSameDay : 是否同一天(0点)
func IsSameDay(t1 *time.Time, t2 *time.Time) bool {
	if t1.Day() == t2.Day() && IsSameMonth(t1, t2) {
		return true
	}
	return false
}

//IsSameDayUnix : 是否同一天(0点)
func IsSameDayUnix(ut1 int64, ut2 int64) bool {
	t1 := time.Unix(ut1, 0)
	t2 := time.Unix(ut2, 0)

	return IsSameDay(&t1, &t2)
}

//IsSameDayWithOffset : 是否同一天（0点+offset秒）
func IsSameDayWithOffset(t1 *time.Time, t2 *time.Time, offset int32) bool {
	d := time.Duration(offset * -1)
	nt1 := t1.Add(d)
	nt2 := t2.Add(d)
	return IsSameDay(&nt1, &nt2)
}

//IsSameDayWithOffsetUnix : 是否同一天（0点+offset秒）
func IsSameDayWithOffsetUnix(ut1 int64, ut2 int64, offset int32) bool {
	t1 := time.Unix(ut1, 0)
	t2 := time.Unix(ut2, 0)
	return IsSameDayWithOffset(&t1, &t2, offset)
}

//CurDayStart : 当前开始时刻（0点）
func CurDayStart() *time.Time {
	curTime := time.Now()
	t := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, curTime.Location())
	return &t
}

//CurDayStartUnix : 当前开始时刻（0点）的Unix值
func CurDayStartUnix() int64 {
	t := CurDayStart()
	return t.Unix()
}

//CurDayStartWithOffset : 当前开始时刻（0点+偏移秒数）
func CurDayStartWithOffset(offset uint32) *time.Time {
	curTime := CurDayStart()
	d := time.Duration(int64(offset) * int64(time.Second))
	t := curTime.Add(d)
	return &t
}

//CurDayStartWithOffsetUnix : 当天开始时刻（0点+偏移秒数）的Unix值
func CurDayStartWithOffsetUnix(offset uint32) int64 {
	t := CurDayStartWithOffset(offset)
	tUnix := t.Unix() + int64(offset)
	return tUnix
}

//CurTimeHourStart : 当前小时开始时刻
func CurTimeHourStart() *time.Time {
	curTime := time.Now()
	t := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), curTime.Hour(), 0, 0, 0, curTime.Location())
	return &t
}

//CurTimeHourStartUnix : 当前小时开始时刻的Unix值
func CurTimeHourStartUnix() int64 {
	t := CurTimeHourStart()
	return t.Unix()
}

//GetTimeDayStartTime : 获取t时刻那天的开始时间（0点）
func GetTimeDayStartTime(t *time.Time) *time.Time {
	tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &tt
}

//GetTimeDayStartTimeUnix : 获取t时刻那天的开始时间的Unix值（0点）
func GetTimeDayStartTimeUnix(t *time.Time) int64 {
	tt := GetTimeDayStartTime(t)
	return tt.Unix()
}

//GetTimeDayStartTimeWithOffset : 获取t时刻那天的开始时间（0点+offset）
func GetTimeDayStartTimeWithOffset(t *time.Time, offset int32) *time.Time {
	d := time.Duration(int64(offset) * int64(time.Second) * -1)
	tt := t.Add(d)
	ptt := GetTimeDayStartTime(&tt)
	d = time.Duration(int64(offset) * int64(time.Second))
	rt := ptt.Add(d)
	return &rt
}

//GetTimeDayStartTimeWithOffsetUnix : 获取t时刻那天的开始时间的Unix值（0点+offset）
func GetTimeDayStartTimeWithOffsetUnix(t *time.Time, offset int32) int64 {
	pt := GetTimeDayStartTimeWithOffset(t, offset)
	return pt.Unix()
}

//GetTimeHourStartTime : 获取t时刻小时开始时刻
func GetTimeHourStartTime(t *time.Time) *time.Time {
	tt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	return &tt
}

//GetTimeHourStartTimeUnix :  获取t时刻小时开始时刻的Unix值
func GetTimeHourStartTimeUnix(t *time.Time) int64 {
	pt := GetTimeHourStartTime(t)
	return pt.Unix()
}
