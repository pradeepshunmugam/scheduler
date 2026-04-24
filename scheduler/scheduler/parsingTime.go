package scheduler

import (
	"fmt"
	"scheduler/scheduler/scheduler/dbcon"
	"time"

	"github.com/robfig/cron/v3"
)

// cronschedule, name, url string, sample int
func calculateNextRun(job jobs) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(job.cron)
	if err != nil {
		fmt.Println("Unable to parse the time. ", err)
	}
	next := schedule.Next(time.Now())
	// fmt.Printf("cron : %v   current time : %v    next schedule time : %v\n", cronschedule, time.Now(), next)
	//Inser next run in db.
	query := (`UPDATE url_list SET next_run = $1 where name = $2`)
	dbcon.DBActivity(query, next, job.name)
	dbcon.CheckNextRun(job)

}
