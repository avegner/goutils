package stats

import (
	"fmt"
	"time"
)

func defaultReportCallback(msg string) {
	fmt.Println(msg)
}

func Took(operation string, report ...func(msg string)) func() {
	start := time.Now()
	r := defaultReportCallback
	if len(report) > 0 {
		r = report[0]
	}
	return func() {
		d := time.Now().Sub(start)
		msg := fmt.Sprintf("%s took %v", operation, d)
		r(msg)
	}
}
