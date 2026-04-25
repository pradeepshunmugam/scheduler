package scheduler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/logger"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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

func UpdateNextRun(db *sql.DB) {
	rows, err := db.Query("SELECT name, url, sample, cron, email, last_run, next_run FROM url_list;")
	if err != nil {
		logger.Log.Error("Unable to connect to DB to read the urls!!, ", zap.Error(err))
	}
	for rows.Next() {
		var job Jobs
		err := rows.Scan(&job.name, &job.url, &job.sample, &job.cron, &job.email, &job.last_run, &job.next_run)
		defer rows.Close()
		if err != nil {
			logger.Log.Error("Unable to scan the urls!!, ", zap.Error(err))
		}
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(job.cron)
		if err != nil {
			fmt.Println("Unable to parse the time. ", err)
		}

		//Inser next run in db.

		next := schedule.Next(time.Now())
		currrentTime := time.Now()
		query := (`UPDATE url_list SET next_run = $1 where name = $2`)
		_, err = db.Exec(query, next, job.name)
		if err != nil {
			logger.Log.Error(`Unable to insert the next_run.`, zap.Error(err))
		} else {
			logger.Log.Info("next run is updated. ", zap.Time("current time :", currrentTime), zap.String("url name : ", job.name))
		}
	}
}

func Run(db *sql.DB) {
	i := 0
	var scheduleList []map[string](interface{})
	query := `SELECT name, url, next_run, sample FROM url_list;`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	var job Jobs
	for rows.Next() {
		err := rows.Scan(&job.name, &job.url, &job.next_run, &job.sample)
		if err != nil {
			fmt.Println(err)
		}
		parsedNextRun, err := time.Parse(time.RFC3339, job.next_run)
		currrentTime := time.Now()
		diff := parsedNextRun.Sub(currrentTime)
		minutesPart := int(diff.Minutes()) % 60
		if minutesPart < 5 {
			i++
			scheduleMap := map[string]interface{}{
				"name":   job.name,
				"url":    job.url,
				"sample": job.sample,
			}
			scheduleList = append(scheduleList, scheduleMap)
		}
	}
	data, err := json.Marshal(scheduleList)
	if err != nil {
		logger.Log.Error("failed to marshal scheduleList")
		return
	}
	logger.Log.Info("Going to schedule  ", zap.Int("total : ", i), zap.String("schedule list : ", string(data)))
	ch := make(chan string) //channel to receive url status
	var wg sync.WaitGroup
	//loop to check the url status cconitnously
	for i := range scheduleList {
		wg.Add(1) //waitgroup to close the function
		url := fmt.Sprint(scheduleList[i]["url"])
		//goroutine created to concurrently execute
		go func(u string) {
			defer wg.Done()
			checkStatus(u, ch)
		}(url)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for i := 0; i < len(scheduleList); i++ {
		logger.Log.Info("url status", zap.String("Result", <-ch)) //receiveing the output through channel
	}

}

func checkStatus(url string, ch chan string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Info("Recovered")
		}
	}()
	client := http.Client{Timeout: 10 * time.Second}
	urlStatus, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		panic("oops something happened !!! Panic")
	}
	res := fmt.Sprintf("Url status of %v is %v .", url, urlStatus.StatusCode)
	ch <- res //sending the output through channel

}
