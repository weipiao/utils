package utils

import "time"

type Clock struct {
	timeFormate string
	loc *time.Location
}

//解析时间日期为对象
func (cl *Clock) ParseStrTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(cl.timeFormate, timeStr)
	return t, err
}

//格式化时间
func (cl *Clock)ToRFC3339(now time.Time) string {
	return now.Format(cl.timeFormate)
}

//今天日期
func (cl *Clock) NowDate() string {
	return time.Now().In(cl.loc).Format(cl.timeFormate)
}

func (cl *Clock) ConvertRFC3339TimeFormat(t string) (string, error) {
	tmp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t, nil
	}
	return tmp.In(cl.loc).Format(cl.timeFormate), nil
}

//计算两个日期的时间差
func  (cl *Clock) GetSecondDiffer(startTime, endTime string) int64 {
	var second int64
	t1, err := time.ParseInLocation(cl.timeFormate, startTime, cl.loc)
	t2, err := time.ParseInLocation(cl.timeFormate, endTime, cl.loc)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		second = diff
		return second
	}

	return second

}

//时间戳转日期
func (cl *Clock)  GetTimeUnixToTimeStr(unix int) string {
	tm := time.Unix(int64(unix), 0).In(cl.loc)
	return tm.Format(cl.timeFormate)
}

//时间字符串转时间戳(string —> int64)
func (cl *Clock)  StrToUnix(date string) int64 { 	//重要：获取时区
	theTime, _ := time.ParseInLocation(cl.timeFormate, date, cl.loc) //使用模板在对应时区转化为time.time类型
	unix := theTime.Unix()                                   //转化为时间戳，类型是int64
	return unix
}

//根据起止时间获取之间的连续时间
func  (cl *Clock)GetDateFromRange(startTime, endTime string) []string {
	startUnix := cl.StrToUnix(startTime)
	endUnix := cl.StrToUnix(endTime)

	//计算日期段内有多少天
	days := (endUnix-startUnix)/86400 + 1;

	//保存每天日期
	var date []string
	for i := 0; i < int(days); i++ {
		date = append(date, time.Unix(startUnix+int64(86400*i), 0).Format("2006-01-02"))
	}

	return date
}

func NewClock() *Clock{
   loc,_:=time.LoadLocation("Local")
   return &Clock{
	  timeFormate:"2006-01-02 15:04:05",
	  loc: loc,
  }
}
