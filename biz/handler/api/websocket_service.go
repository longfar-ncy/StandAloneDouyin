package api

import (
	"context"
	"douyin/dal/db"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"strings"
)

type Client struct {
	ID       string
	ToUserID string
	Conn     *websocket.Conn
	Send     chan []byte
}

func (c *Client) readPump() {
	defer func() {
		MannaClient.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		//c.Conn.PongHandler()
		SendMsg := new(SendMsg)

		err := c.Conn.ReadJSON(&SendMsg)
		if err != nil {
			hlog.Error("api.websocket_service.readPump.ReadJSON err:", err.Error())
			MannaClient.Unregister <- c
			_ = c.Conn.Close()
			break
		}

		if SendMsg.Type == 1 { // 发送消息

			MannaClient.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(SendMsg.Content), // 发送过来的消息
			}
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send: // 对方在线逻辑
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			ReplyMsg := ReplyMsg{
				Code:    777777,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(ReplyMsg)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)

			uid, touid, err := ExtractNumbers(c.ToUserID)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//fmt.Println(msg)
			err = db.CreateMessage(uid, touid, string(msg[26:len(msg)-2])) // 将消息放到数据库
			if err != nil {
				hlog.Error("api.websocket_service.writePump.CreateMessage err:", err.Error())
			}
		case message, ok := <-c.Send: // 不在线逻辑
			if ok {
				ReplyMsg := ReplyMsg{
					Code:    777777,
					Content: fmt.Sprintf("%s", string(message)),
				}
				msg, _ := json.Marshal(ReplyMsg)
				_ = c.Conn.WriteMessage(websocket.TextMessage, msg)

				uid, touid, err := ExtractNumbers(c.ToUserID)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend {
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "不能给非好友发消息"
					hlog.Error("api.websocket_service.writePump.IsFriend err:", errNo.Error())
				}

				//fmt.Println(msg)
				err = db.CreateMessage(uid, touid, string(msg[26:len(msg)-2])) // 将消息放到数据库
				if err != nil {
					hlog.Error("api.websocket_service.writePump.CreateMessage err:", err.Error())
				}
			}
		}
	}
}

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

func CreateID(uid, toUid string) string {
	var builder strings.Builder
	builder.WriteString(uid)
	builder.WriteString("->")
	builder.WriteString(toUid)
	return builder.String()
}

func ServeWs(ctx context.Context, c *app.RequestContext) {
	fromUserID := c.GetUint64(global.Config.JWTConfig.IdentityKey)
	hlog.Info("biz.handler.api.websocket_service.ServeWs GetFromUserID:", fromUserID)
	toUid := c.Query("to_user_id")

	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		client := &Client{
			ID:       CreateID(strconv.FormatUint(fromUserID, 10), toUid),
			ToUserID: CreateID(toUid, strconv.FormatUint(fromUserID, 10)),
			Conn:     conn,
			Send:     make(chan []byte),
		}

		MannaClient.Register <- client

		go client.writePump()
		client.readPump()
	})
	if err != nil {
		hlog.Error("biz.handler.api.websocket_service.ServeWs err:", err.Error())
	}
}
