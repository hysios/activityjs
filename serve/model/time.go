package model

import "strconv"

type time struct {
	sec    int
	hour   int
	min    int
	wday   int
	year   int
	yday   int
	mon    int
	mday   int
	isdst  int
	offset int
}

func getTime(t int) time {
	var tt time
	t /= 1000
	day := t / (24 * 60 * 60)
	secs := t % (24 * 60 * 60)
	tt.sec = secs % 60
	mins := secs / 60
	tt.hour = mins / 60
	tt.min = mins % 60
	tt.wday = (day + 4) % 7
	year := (((day * 4) + 2) / 1461)
	tt.year = year + 70
	day -= ((year * 1461) + 1) / 4

	tt.yday = day
	if leap, ok := tt.isLeap(); day > 58+leap {
		if ok {
			day += 1
		} else {
			day += 2
		}
	}

	tt.mon = ((day * 12) + 6) / 367
	tt.mday = day + 1 - ((tt.mon*367)+5)/12
	tt.isdst = 0
	return tt
}

func (t *time) isLeap() (int, bool) {
	if !(t.year&3 > 0) {
		return 1, true
	}
	return 0, false
}

func (t *time) String() string {
	return strconv.Itoa(t.year+1900) +
		"-" + pad(t.mon+1) +
		"-" + pad(t.mday) + " " +
		pad(t.hour+t.offset) + ":" +
		pad(t.min) + ":" +
		pad(t.sec)
}

func (t *time) Unix() int {
	year := t.year - 70
	days := (year*1461 - 2) / 4
	secs := days * 24 * 60 * 60
	secs += (t.yday + 1) * 24 * 60 * 60
	secs += (t.hour + t.offset) * 60 * 60
	secs += t.min * 60
	secs += t.sec
	secs *= 1000
	return secs
}

func (t Time) String() string {
	tt := getTime(int(t))
	return strconv.Itoa(tt.year+1900) +
		"-" + pad(tt.mon+1) +
		"-" + pad(tt.mday) + " " +
		pad(tt.hour+tt.offset) + ":" +
		pad(tt.min) + ":" +
		pad(tt.sec)

}

func (t *time) Timezone(offst int) time {
	var tt = *t
	tt.offset = offst
	return tt
}

func (t Time) Timezone(offst int) Time {
	tt := getTime(int(t))
	tt = tt.Timezone(offst)

	return Time(tt.Unix())
}

func pad(n int) string {
	if n >= 10 {
		return strconv.Itoa(n)
	} else {
		return "0" + strconv.Itoa(n)
	}
}
