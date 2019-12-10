package gobilibili

import (
	"fmt"
	"strings"
)

//DefaultHandler print cmd msg log
func DefaultHandler(c *Context) bool {
	cmd, err := c.Msg.Get("cmd").String()
	if err != nil {
		return true
	}
	if cmd == "LIVE" {
		fmt.Println("直播开始。。。")
		return false
	}
	if cmd == "PREPARING" {
		fmt.Println("房主准备中。。。")
		return false
	}
	// 弹幕升级了，弹幕cmd获得的值不是DANMU_MSG, 而是DANMU_MSG: + 版本, 例如: DANMU_MSG:4:0:2:2:2:0
	if cmd == "DANMU_MSG" || strings.HasPrefix(cmd, "DANMU_MSG") {
		commentText, err := c.Msg.Get("info").GetIndex(1).String()
		if err != nil {
			fmt.Println("Json decode error failed: ", err)
			return false
		}

		commentUser, err := c.Msg.Get("info").GetIndex(2).GetIndex(1).String()
		if err != nil {
			fmt.Println("Json decode error failed: ", err)
			return false
		}
		fmt.Println(commentUser, " say: ", commentText)
		return false
	}
	return false
}

//DebugHandler debug msg info
func DebugHandler(c *Context) bool {
	jbytes, _ := c.Msg.EncodePretty()
	fmt.Println(string(jbytes))
	return false
}
