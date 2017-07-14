package main

import (
	"github.com/lyyyuna/gobilibili"
)

func main() {
	bili := gobilibili.NewBiliBiliClient()
	go bili.ConnectServer(115)
	bili.HeartbeatLoop()
}
