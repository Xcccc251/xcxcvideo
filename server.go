package main

import (
	"XcxcVideo/common/helper"
	"XcxcVideo/router"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

// 定义升级器
// 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源连接（生产环境需加强验证）
	},
}

// 连接管理器
type ConnectionManager struct {
	connections map[int]*websocket.Conn // 用户ID到连接的映射
	mu          sync.Mutex              // 保护并发访问
}

// 创建一个全局的连接管理器
var connManager = ConnectionManager{
	connections: make(map[int]*websocket.Conn),
	mu:          sync.Mutex{},
}

// 添加连接
func (cm *ConnectionManager) Add(userId int, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[userId] = conn
}

// 移除连接
func (cm *ConnectionManager) Remove(userId int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, userId)
}

// 发送消息给指定用户
func (cm *ConnectionManager) Send(userId int, message []byte) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	conn, ok := cm.connections[userId]
	if !ok {
		return fmt.Errorf("user %s not connected", userId)
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}

type Message struct {
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 处理 WebSocket 连接
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()
	// 首次读取消息以获取 Token
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read initial message:", err)
		return
	}
	var msg Message
	err = json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("Failed to unmarshal initial message:", err)
		return
	}
	token := strings.TrimPrefix(msg.Content, "Bearer ")
	userClaim, _ := helper.AnalysisToken(token)
	userId := userClaim.UserId

	connManager.Add(userId, conn)
	defer connManager.Remove(userId)

	log.Printf("WebSocket connection established for userId: %d", userId)
	//go func() {
	//	newmsg := struct {
	//		Type string      `json:"Type"`
	//		Data interface{} `json:"Data"`
	//	}{
	//		Type: "reply",
	//		Data: "hello",
	//	}
	//	jsonMsg, _ := json.Marshal(&newmsg)
	//
	//	for {
	//		fmt.Println("send to", userId)
	//		connManager.Send(userId, jsonMsg)
	//		time.Sleep(3 * time.Second)
	//	}
	//}()

	// 后续循环读取消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for userID %d: %v", userId, err)
			break
		}
		log.Printf("Received from userID %d: %s", userId, message)
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
