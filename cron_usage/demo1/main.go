package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	nextTime = expr.Next(now)
	fmt.Println(now, nextTime)

	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("is scheduler")
	})
	time.Sleep(time.Second * 30)
}
