package main

import (
	"XcxcVideo/router"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 定义升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，实际项目中请根据需求设置
	},
}

// 处理 WebSocket 连接
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close() // 确保关闭连接

	log.Println("WebSocket connection established")

	// 循环读取消息并回显
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s\n", message)

		// 回显消息
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
func main() {
	r := router.Router()
	go func() {
		// 配置 WebSocket 路由
		http.HandleFunc("/im", handleWebSocket)
		// 单独启动 WebSocket 服务器
		port := ":7071"
		log.Printf("WebSocket server started on port %s\n", port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatalf("WebSocket server failed to start: %v\n", err)
		}
	}()
	r.Run(":7070")

}
