package utils

import (
	"abhi/search/search"
	"fmt"

	"github.com/robfig/cron"
)

func StartCronJob(){

	c := cron.New()
	c.AddFunc("0 * * * *", search.RunEngine) // run every hour
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs \n", cronCount)
}