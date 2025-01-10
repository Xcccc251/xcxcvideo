package websocketServer

import (
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ImResponse struct {
	Type string        `json:"type"`
	Time models.MyTime `json:"time"`
	Data interface{}   `json:"data"`
}

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
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
	// 设置客户端的 Pong 消息处理器，用于检测客户端响应
	conn.SetPongHandler(func(appData string) error {
		// 当接收到 Pong 消息时，更新读取期限
		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		return nil
	})

	// 设置读取超时
	go func() {
		for {
			// 每隔3秒发送 Ping 消息
			err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Printf("Ping error for userId %d: %v", userId, err)
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	var messageChan = make(chan models.ImMessage, 10)
	go func() {
		for {
			select {
			case imMessage := <-messageChan:
				switch imMessage.Code {
				case 101:
					sendWhisper(imMessage)
				}
			}
		}
	}()
	// 后续循环读取消息

	for {
		_, message, err := conn.ReadMessage()
		var messageMap map[string]interface{}
		json.Unmarshal(message, &messageMap)
		if code, ok := messageMap["code"].(float64); ok {
			fmt.Println("code", code)
			imMessage := models.ImMessage{}
			imMessage.Code = int(code)
			imMessage.Message = messageMap
			imMessage.UserId = userId
			messageChan <- imMessage
		}
		if err != nil {
			log.Printf("Read error for userID %d: %v", userId, err)
			break
		}
		log.Printf("Received from userID %d: %s", userId, message)
	}

}

func SendMessage(userId int, typeOfMsg string, data interface{}) error {
	message, _ := json.Marshal(ImResponse{
		Type: typeOfMsg,
		Time: models.MyTime(time.Now()),
		Data: data,
	})
	cm := connManager
	cm.mu.Lock()
	defer cm.mu.Unlock()
	conn, ok := cm.connections[userId]
	if !ok {
		return fmt.Errorf("user %s not connected", userId)
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}
