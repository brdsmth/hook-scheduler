package services

import (
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

func RunCron() {
	log.Print("Initiated cron")
	c := cron.New()
	c.AddFunc("* * * * *", func() { // Runs every minute
		currentTime := time.Now().Format(time.RFC3339)
		log.Println()
		log.Println("---	Cron running	---")
		log.Printf("Current time:\t%s", currentTime)
		jobs, err := GetDueJobs() // Implement this function
		if err != nil {
			log.Println("Error fetching jobs:", err)
			return
		}
		Queue(jobs)
		log.Println("---	Cron stopped	---")
		log.Println()
	})
	c.Start()
}
