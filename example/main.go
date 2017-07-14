package main

import (
	"github.com/lyyyuna/gobilibili"
)

func main() {
	bili := gobilibili.NewBiliBiliClient()
	bili.ConnectServer(1016)
}
