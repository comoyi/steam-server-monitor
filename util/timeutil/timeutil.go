package timeutil

import "fmt"

func FormatDuration(second int64) string {
	var d int64
	var h int64
	var m int64
	var s int64
	var str string
	var flag = false

	d = second / 86400
	second -= d * 86400
	h = second / 3600
	second -= h * 3600
	m = second / 60
	second -= m * 60
	s = second

	if d > 0 {
		flag = true
		str = fmt.Sprintf("%s%d天", str, d)
	}
	if flag || h > 0 {
		flag = true
		str = fmt.Sprintf("%s%d时", str, h)
	}
	if flag || m > 0 {
		flag = true
		str = fmt.Sprintf("%s%d分", str, m)
	}
	str = fmt.Sprintf("%s%d秒", str, s)
	return str
}
