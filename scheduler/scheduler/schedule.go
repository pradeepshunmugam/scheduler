package scheduler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/scheduler/db.go"
	"scheduler/scheduler/logger"
	"strconv"
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
	last_run time.Time
	next_run time.Time
}

func UpdateNextRun() {
	db := db.GetDB()
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
		query := (`UPDATE url_list SET next_run = $1 where name = $2`)
		_, err = db.Exec(query, next, job.name)
		if err != nil {
			logger.Log.Error(`Unable to insert the next_run.`, zap.Error(err))
		} else {
			logger.Log.Info("next run is updated. ", zap.Time("Next run :", next), zap.String("url name : ", job.name))
		}
	}
}

func Run() {
	db := db.GetDB()
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
		if !job.next_run.After(time.Now()) {
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
	//loop to check the url status conitnously
	if !job.next_run.After(time.Now()) { //Run all the job which passed the next run
		UpdateNextRun() //update the next run so that next run will be updated and the job will not conituously for which passed the next run.
		for i := range scheduleList {
			wg.Add(1) //waitgroup to close the function
			url := fmt.Sprint(scheduleList[i]["url"])
			name := fmt.Sprint(scheduleList[i]["name"])
			sampleStr := fmt.Sprint(scheduleList[i]["sample"])

			sample, err := strconv.Atoi(sampleStr)
			if err != nil {
				logger.Log.Error("Error in converting sample to integer", zap.Error(err))
			}
			//goroutine created to concurrently execute
			go func(u string, name string, s int) {
				defer wg.Done()
				checkStatus(u, name, sample, ch)
			}(url, name, sample)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for i := 0; i < len(scheduleList); i++ {
		logger.Log.Info("url status", zap.String("Result", <-ch)) //receiveing the output through channel
	}

}

func checkStatus(url, name string, sample int, ch chan string) {
	db := db.GetDB()
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Info("Recovered")
		}
	}()
	client := http.Client{Timeout: 10 * time.Second}
	urlStatus, err := client.Get(url)
	if err != nil {
		// logger.Log.Error("http get error", zap.Error(err))
		// sendNotification(name, err.Error()) //adding this on 11-Jun-2026
		// // panic("oops something happened !!! Panic")
		// return

		fmt.Println("http get error:", err)
		sendNotification(name, err.Error())
		return
	}
	res := fmt.Sprintf("Url status of %v is %v .", url, urlStatus.StatusCode)
	if urlStatus.StatusCode >= 200 {
		sendNotification(name, res) //adding this on 11-Jun-2026
	}
	insert := `INSERT INTO urlstatus(name, status,statuscode) VALUES($1, $2, $3)`
	_, err = db.Exec(insert, name, urlStatus.Status, urlStatus.StatusCode)
	if err != nil {
		fmt.Println("unable to insert result in DB.", err)
	}
	updateLastRun := `UPDATE url_list SET last_run = (SELECT event_timestamp FROM urlstatus WHERE urlstatus.name = url_list.name  ORDER BY event_timestamp DESC LIMIT 1 ) WHERE name IN (SELECT $1 FROM urlstatus);`
	_, err = db.Exec(updateLastRun, name)
	if err != nil {
		fmt.Println("Unable to insert last run", err)
	}
	ch <- res //sending the output through channel
}
