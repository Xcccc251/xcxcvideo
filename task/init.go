package task

import (
	"github.com/robfig/cron/v3"
	"log"
)

func Init() {
	cr := cron.New()

	cr.AddFunc("0/5 * * * * ? ", func() {
		log.Println("执行定时任务")
	})
	cr.Start()
}
