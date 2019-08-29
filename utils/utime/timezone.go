package utime

import (
	"time"
)

var TimezoneIndia = time.FixedZone("IST", 5*3600+1800) // 印度时间UTC+5:30
var TimezoneCST = time.FixedZone("CST", 8*3600)
