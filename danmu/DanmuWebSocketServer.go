package danmu

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"XcxcVideo/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DanmuMessage struct {
	Token string `json:"token"`
	Data  Data   `json:"data"`
}
type Data struct {
	Content   string  `json:"content"`
	Fontsize  int     `json:"fontsize"`
	Mode      int     `json:"mode"`
	Color     string  `json:"color"`
	TimePoint float64 `json:"timePoint"`
}
type ImResponse struct {
	Type string        `json:"type"`
	Time models.MyTime `json:"time"`
	Data interface{}   `json:"data"`
}

// 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（开发阶段）
	},
}

// 视频连接管理
var videoConnectionMap = struct {
	sync.RWMutex
	connections map[int]map[*websocket.Conn]bool
}{
	connections: make(map[int]map[*websocket.Conn]bool),
}

// 发送消息给指定视频的所有连接
func sendMessage(vid int, text string) {
	videoConnectionMap.RLock()
	// 获取指定视频的连接列表
	connections, exists := videoConnectionMap.connections[vid]
	defer videoConnectionMap.RUnlock()

	if !exists {
		return
	}

	// 并行发送消息
	var wg sync.WaitGroup
	for conn := range connections {
		wg.Add(1)
		go func(c *websocket.Conn) {
			defer wg.Done()
			// 向客户端发送消息
			if err := c.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
				videoConnectionMap.Lock()
				delete(videoConnectionMap.connections[vid], c) // 移除出错的连接
				videoConnectionMap.Unlock()
				c.Close() // 关闭出错的连接
			}
		}(conn)
	}
	wg.Wait() // 等待所有 Goroutine 完成
}

// 处理 WebSocket 连接
func DanmuWebSocketHandler(c *gin.Context) {
	vid, _ := strconv.Atoi(c.Param("vid"))
	// 升级 HTTP 请求为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	// 连接关闭
	defer func() {
		videoConnectionMap.Lock()
		if _, exists := videoConnectionMap.connections[vid]; exists {
			delete(videoConnectionMap.connections[vid], conn)
			if len(videoConnectionMap.connections[vid]) == 0 {
				delete(videoConnectionMap.connections, vid)
			}
		}
		videoConnectionMap.Unlock()
		sendMessage(vid, "当前观看人数"+strconv.Itoa(len(videoConnectionMap.connections[vid])))

		conn.Close()
	}()
	// 添加连接到视频连接映射
	videoConnectionMap.Lock()
	if videoConnectionMap.connections[vid] == nil {
		videoConnectionMap.connections[vid] = make(map[*websocket.Conn]bool)
	}
	videoConnectionMap.connections[vid][conn] = true
	videoConnectionMap.Unlock()

	// 发送当前观看人数
	sendMessage(vid, "当前观看人数"+strconv.Itoa(len(videoConnectionMap.connections[vid])))

	// 后续循环处理消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // 连接关闭时退出循环
		}
		sendDanmu(msg, conn, vid)
		log.Printf("Received from vid %d: %s", vid, msg)
	}
}

func sendDanmu(msg []byte, c *websocket.Conn, vid int) {
	var danmuMessage DanmuMessage
	json.Unmarshal(msg, &danmuMessage)

	token := strings.TrimPrefix(danmuMessage.Token, "Bearer ")
	userClaim, err := helper.AnalysisToken(token)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("登录已过期"))
		return
	}
	userId := userClaim.UserId
	result, err := models.RDb.Exists(context.Background(), define.USER_PREFIX+strconv.Itoa(userId)).Result()
	if result == 0 {
		c.WriteMessage(websocket.TextMessage, []byte("登录已过期"))
		return
	}

	var danmu models.Danmu
	danmu.Uid = userId
	danmu.Vid = vid
	danmu.Color = danmuMessage.Data.Color
	danmu.Content = danmuMessage.Data.Content
	danmu.Fontsize = danmuMessage.Data.Fontsize
	danmu.Mode = danmuMessage.Data.Mode
	danmu.TimePoint = danmuMessage.Data.TimePoint
	danmu.CreateDate = models.MyTime(time.Now())
	//todo 消息队列
	models.Db.Model(new(models.Danmu)).Create(&danmu)
	service.UpdateVideoStats(vid, "danmu", true, 1)
	models.RDb.SAdd(context.Background(), define.DANMU_IDSET+strconv.Itoa(vid), danmu.Id)
	danmuJson, _ := json.Marshal(&danmu)
	sendMessage(vid, string(danmuJson))

}
