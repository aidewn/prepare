package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	var (
		schedule map[string]*CronJob
		cronJob  *CronJob
		expr     *cronexpr.Expression
		now      time.Time
		err      error
	)

	schedule = make(map[string]*CronJob)
	now = time.Now()
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	schedule["job1"] = cronJob

	if expr, err = cronexpr.Parse("*/6 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	schedule["job2"] = cronJob

	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)
		for {
			now = time.Now()
			for jobName, cronJob = range schedule {
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					go func(jobName string) {
						fmt.Println(jobName, " is run")
					}(jobName)
				}

				cronJob.nextTime = cronJob.expr.Next(now)
				fmt.Println(jobName, "next run ", cronJob.nextTime)
			}

			select {
			case <-time.NewTimer(100 * time.Millisecond).C:
			}
		}

	}()

	time.Sleep(100 * time.Second)
}
