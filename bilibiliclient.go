package gobilibili

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"

	"time"

	"github.com/bitly/go-simplejson"
)

const (
	min = 100000000000000
	max = 200000000000000
)

// BiliBiliClient define
type BiliBiliClient struct {
	CIDInfoUrl      string
	roomId          int
	ChatPort        int
	protocolversion uint16
	ChatHost        string
	serverConn      net.Conn
	uid             int
}

func NewBiliBiliClient() *BiliBiliClient {
	bili := new(BiliBiliClient)
	bili.ChatHost = "livecmt-1.bilibili.com"
	bili.ChatPort = 788
	bili.protocolversion = 1
	return bili
}

// ConnectServer define
func (bili *BiliBiliClient) ConnectServer(roomId int) {
	fmt.Println("Entering room ....")

	dstAddr := fmt.Sprintf("%s:%d", bili.ChatHost, bili.ChatPort)
	dstConn, err := net.Dial("tcp", dstAddr)
	if err != nil {
		fmt.Println("Failed to connect bilibili server.")
		return
	}
	bili.serverConn = dstConn
	fmt.Println("弹幕链接中。。。")
	bili.SendJoinChannel(bili.roomId)
	go bili.HeartbeatLoop()
	go bili.ReceiveMessageLoop()
	for {
		time.Sleep(time.Second * 10)
	}
}

// HeartbeatLoop define
func (bili *BiliBiliClient) HeartbeatLoop() {
	for {
		bili.SendSocketData(0, 16, bili.protocolversion, 2, 1, "")
		time.Sleep(time.Second * 30)
	}
}

// SendJoinChannel define
func (bili *BiliBiliClient) SendJoinChannel(channelId int) {
	bili.uid = rand.Intn(max) + min
	body := fmt.Sprintf("{'roomid':%d,'uid':%d}", channelId, bili.uid)
	bili.SendSocketData(0, 16, bili.protocolversion, 7, 1, body)
}

// SendSocketData define
func (bili *BiliBiliClient) SendSocketData(packetlength uint32, magic uint16, ver uint16, action uint32, param uint32, body string) {
	bodyBytes := []byte(body)
	if packetlength == 0 {
		packetlength = uint32(len(bodyBytes) + 16)
	}
	headerBytes := new(bytes.Buffer)
	var data = []interface{}{
		packetlength,
		magic,
		ver,
		action,
		param,
	}
	for _, v := range data {
		err := binary.Write(headerBytes, binary.BigEndian, v)
		if err != nil {
			fmt.Println("binary write failed: ", err)
		}
	}
	socketData := append(headerBytes.Bytes(), bodyBytes...)
	bili.serverConn.Write(socketData)
}

func (bili *BiliBiliClient) ReceiveMessageLoop() {
	for {
		buf := make([]byte, 4)
		if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
			fmt.Println("binary read failed: ", err)
		}
		expr := binary.BigEndian.Uint32(buf)
		buf = make([]byte, 4)
		if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
			fmt.Println("binary read failed: ", err)
		}
		buf = make([]byte, 4)
		if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
			fmt.Println("binary read failed: ", err)
		}
		num := binary.BigEndian.Uint32(buf)
		buf = make([]byte, 4)
		if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
			fmt.Println("binary read failed: ", err)
		}

		num2 := expr - 16
		if num2 != 0 {
			num = num - 1
			if num == 0 || num == 1 || num == 2 {
				buf = make([]byte, 4)
				if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
					fmt.Println("binary read failed: ", err)
				}
				num3 := binary.BigEndian.Uint32(buf)
				fmt.Println("房间人数为：", num3)
				continue
			} else if num == 3 || num == 4 {
				buf = make([]byte, num2)
				if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
					fmt.Println("binary read failed: ", err)
				}
				messages := string(buf)
				bili.parseDanMu(messages)
				continue
			} else if num == 5 || num == 6 || num == 7 {
				buf = make([]byte, num2)
				if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
					fmt.Println("binary read failed: ", err)
				}
				continue
			} else {
				if num != 16 {
					buf = make([]byte, num2)
					if _, err := io.ReadFull(bili.serverConn, buf); err != nil {
						fmt.Println("binary read failed: ", err)
					}
					continue
				} else {
					continue
				}
			}
		}
	}
}

func (bili *BiliBiliClient) parseDanMu(message string) {
	dic, err := simplejson.NewJson([]byte(message))
	if err != nil {
		fmt.Println("Json decode error failed: ", err)
		return
	}

	cmd, err := dic.Get("cmd").String()
	if err != nil {
		fmt.Println("Json decode error failed: ", err)
		return
	}

	if cmd == "LIVE" {
		fmt.Println("直播开始。。。")
		return
	}
	if cmd == "PREPARING" {
		fmt.Println("房主准备中。。。")
		return
	}
	if cmd == "DANMU_MSG" {
		commentText, err := dic.Get("info").GetIndex(1).String()
		if err != nil {
			fmt.Println("Json decode error failed: ", err)
			return
		}

		commentUser, err := dic.Get("info").GetIndex(2).GetIndex(1).String()
		if err != nil {
			fmt.Println("Json decode error failed: ", err)
			return
		}

		fmt.Println(commentUser, " say: ", commentText)
		return
	}

}
