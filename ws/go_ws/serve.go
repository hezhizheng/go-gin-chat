package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-gin-chat/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 客户端连接详情
type wsClients struct {
	Conn *websocket.Conn `json:"conn"`

	RemoteAddr string `json:"remote_addr"`

	Uid float64 `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	AvatarId string `json:"avatar_id"`

	ToUser interface{} `json:"to_user"`
}

// client & serve 的消息体
type msg struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// 变量定义初始化
var (
	wsUpgrader = websocket.Upgrader{}

	clientMsg = msg{}

	mutex = sync.Mutex{}

	rooms = [roomCount + 1][]wsClients{}

	rooms2 = make(chan wsClients)

	privateChat = []wsClients{}

	privateChat2 = make(chan wsClients)

	clientMsg2 = make(chan []byte)

	clientMsg3 = make(chan msg)
)

// 定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTypeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数

func Run(gin *gin.Context) {

	// @see https://github.com/gorilla/websocket/issues/523
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, _ := wsUpgrader.Upgrade(gin.Writer, gin.Request, nil)

	defer c.Close()

	go read(c)
	go write(c)

	select {}

}

func read(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		log.Println("message", string(message))
		if err != nil { // 离线通知
			log.Println("ReadMessage error1", err)
			break
		}

		serveMsgStr := message

		// 处理心跳响应 , heartbeat为与客户端约定的值
		//log.Println(string(serveMsgStr))
		if string(serveMsgStr) == `heartbeat` {

			clientMsg3 <- msg{
				Status: 0,
				Data:   "heartbeat ok",
			}
		}

		json.Unmarshal(message, &clientMsg)
		// log.Println("来自客户端的消息", clientMsg,c.RemoteAddr())
		if clientMsg.Data != nil {
			if clientMsg.Status == msgTypeOnline { // 进入房间，建立连接
				roomId, _ := getRoomId()

				rooms2 <- wsClients{
					Conn:       c,
					RemoteAddr: c.RemoteAddr().String(),
					Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
					Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
					RoomId:     roomId,
					AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
				}
			}

			_, serveMsg := formatServeMsgStr(clientMsg.Status)
			clientMsg3 <- serveMsg
		}
	}
}

func write(c *websocket.Conn) {
	//r := make(chan wsClients)
	for {
		select {
		case r := <-rooms2:
			log.Println("room2", r, c.RemoteAddr())
		case cl := <-clientMsg3:

			serveMsgStr, _ := json.Marshal(cl)

			switch cl.Status {

			case 0:
				c.WriteMessage(websocket.TextMessage, serveMsgStr)

			}

			log.Println("cl", cl, c.RemoteAddr())
		}
	}
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	_, roomIdInt := getRoomId()
	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	for index, con := range rooms[roomIdInt] {
		if con.RemoteAddr == conn.RemoteAddr().String() {
			data := map[string]interface{}{
				"username": con.Username,
				"uid":      con.Uid,
				"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
			}

			jsonStrServeMsg := msg{
				Status: msgTypeOffline,
				Data:   data,
			}
			serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

			disMsg := string(serveMsgStr)

			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:index], rooms[roomIdInt][index+1:]...)
			mutex.Unlock()
			notify(conn, disMsg)
		}
	}
}

// 格式化传送给客户端的消息数据
func formatServeMsgStr(status int) ([]byte, msg) {

	roomId, roomIdInt := getRoomId()

	data := map[string]interface{}{
		"username": clientMsg.Data.(map[string]interface{})["username"].(string),
		"uid":      clientMsg.Data.(map[string]interface{})["uid"].(float64),
		"room_id":  roomId,
		"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
	}

	if status == msgTypeSend || status == msgTypePrivateChat {
		data["avatar_id"] = clientMsg.Data.(map[string]interface{})["avatar_id"].(string)
		data["content"] = clientMsg.Data.(map[string]interface{})["content"].(string)

		toUidStr := clientMsg.Data.(map[string]interface{})["to_uid"].(string)
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := strconv.FormatFloat(data["uid"].(float64), 'f', -1, 64)
		intUid, _ := strconv.Atoi(stringUid)

		if _, ok := clientMsg.Data.(map[string]interface{})["image_url"]; ok {
			// 存在图片
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data["content"],
				"room_id":    data["room_id"],
				"image_url":  clientMsg.Data.(map[string]interface{})["image_url"].(string),
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
			})
		}

	}

	if status == msgTypeGetOnlineUser {
		data["count"] = GetOnlineRoomUserCount(roomIdInt)
		data["list"] = onLineUserList(roomIdInt)
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg
}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.(map[string]interface{})["room_id"].(string)

	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

// 获取在线用户列表
func onLineUserList(roomId int) []wsClients {
	return rooms[roomId]
}

// =======================对外方法=====================================

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
