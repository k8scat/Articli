package juejin

import (
	"strconv"
	"time"
)

func FormatTime(s, layout string) string {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ""
	}
	return time.Unix(i, 0).Format(layout)
}
