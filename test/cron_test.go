package test

import (
	"github.com/robfig/cron/v3"
	"log"
	"testing"
)

func TestCron(t *testing.T) {
	// 创建一个 Cron 实例
	c := cron.New()

	// 添加定时任务
	_, err := c.AddFunc("@every 10s", func() {
		log.Println("每隔 10 秒执行一次任务")
	})
	if err != nil {
		log.Fatalf("添加任务失败: %v", err)
	}

	// 启动 Cron
	c.Start()
	defer c.Stop() // 确保程序退出时停止 Cron

	// 主程序阻塞
	log.Println("定时任务已启动，按 Ctrl+C 退出程序")
	select {}
}
