package main

import (
	"XcxcVideo/danmu"
	websocketServer "XcxcVideo/handler"
	"XcxcVideo/router"
	"XcxcVideo/task"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	go func() {
		// 创建一个 Cron 实例
		c := cron.New()
		// 添加定时任务
		_, err := c.AddFunc("@every 1800s", task.UpdateHotSearchWord)
		if err != nil {
			log.Fatalf("添加任务失败: %v", err)
		}
		// 启动 Cron
		c.Start()
		defer c.Stop() // 确保程序退出时停止 Cron
		select {}
	}()
	go func() {
		// 配置 WebSocket 路由
		http.HandleFunc("/im", websocketServer.HandleWebSocket)
		// 单独启动 WebSocket 服务器
		port := ":7071"
		log.Printf("WebSocket server started on port %s\n", port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatalf("WebSocket server failed to start: %v\n", err)
		}
	}()

	r.GET("/ws/danmu/:vid", danmu.DanmuWebSocketHandler)
	r.Run(":7070")

}
