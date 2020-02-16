package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02 15:04:15"
)

func GetNow() time.Time {
	zone, _ := time.LoadLocation("Asia/Taipei")
	return time.Now().In(zone)
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}
