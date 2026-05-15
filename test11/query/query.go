package query

import (
	"fmt"
	"scheduler/test11/db"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

type Jobs struct {
	url      string
	name     string
	sample   int
	cron     string
	email    string
	last_run string
	next_run string
}

func SelectQuery() {
	db := db.GetDB()
	query := `SELECT name, url, sample, cron, email, last_run, next_run FROM url_list;`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Unable to read DB rows", err)
	}
	for rows.Next() {
		var job Jobs
		err := rows.Scan(&job.name, &job.url, &job.sample, &job.cron, &job.email, &job.last_run, &job.next_run)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(job.cron)
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(job.cron)
		if err != nil {
			fmt.Println("Unable to parse the time. ", err)
		}
		fmt.Print(schedule.Next(time.Now()))
		schedule.
	}
}
