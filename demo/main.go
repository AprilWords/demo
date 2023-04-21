package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
 1. 显示抽奖开始消息：

{"action": "start"}

2. 关闭抽奖消息：
{"action": "stop"}

3. 展示抽奖数字
{"action": "show-number", "key": 0, "val": 12}
*/

var ms = make(chan []byte, 1000)
var wss []*websocket.Conn
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewConnection(c *gin.Context) {

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务端错误",
		})
		return
	}
	wss = append(wss, ws)
	// 开始通信

}

// http处理
func dealRequest(c *gin.Context) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	ms <- buf.Bytes()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
	})

}
func main() {
	r := gin.Default()

	r.GET("/wsTest", NewConnection)
	r.POST("/api/order", dealRequest)
	go func() {
		for message := range ms {
			for _, wsClient := range wss {
				wsClient.WriteMessage(websocket.TextMessage, message)

			}
		}
	}()
	r.Run("0.0.0.0:9966")
}
