# gobilibili
![](https://ws1.sinaimg.cn/large/521c75dcly1fo040yjlaij20g90a6myb.jpg)

B 站直播弹幕 Go 版。
在[原项目](https://github.com/lyyyuna/gobilibili) 基础上作了以下修改:

* 自动获取直播间真实ID,兼容短ID
* 增加赠送礼物/进入房间在线人数变动消息处理
* 增加方便的消息事件订阅机制，用于将自己的处理器注册到消息处理链
* 原处理逻辑现在被包装为一个处理器，名为DefaultHandler
* 增加了一些函数，用于事件触发时快速获取消息结构(徽章，等级，礼物信息，类型，数量，标志等等)
* 为大部分函数添加了错误处理

## 安装
    go get github.com/bigemon/gobilibili
## 示例

### 实时打印弹幕

```
package main

import "github.com/bigemon/gobilibili"

func main() {
	bili := gobilibili.NewBiliBiliClient()
	bili.RegHandleFunc(gobilibili.CmdAll, gobilibili.DefaultHandler)
}
```
#### 事件订阅
如果你希望订阅不同的事件，请尝试`gobilibili.Cmd*`开头的一系列常量。
以下是一些示例,你也可以随时在example目录下查看.


*订阅弹幕事件，并输出弹幕信息*

```
bili := gobilibili.NewBiliBiliClient()
bili.RegHandleFunc(gobilibili.CmdDanmuMsg, func(c *gobilibili.Context) bool {
	dinfo := c.GetDanmuInfo()
	log.Printf("[%d]%d 说: %s\r\n", c.RoomID, dinfo.UID, dinfo.Text)
	return false
})
```

*进入房间*

```
bili.RegHandleFunc(gobilibili.CmdWelcome, func(c *gobilibili.Context) bool {
	winfo := c.GetWelcomeInfo()
	if winfo.Uname != "" {
		log.Printf("[%d]%s 进入了房间\r\n", c.RoomID, winfo.Uname)
	} else {
		log.Printf("[%d]%d 进入了房间\r\n", c.RoomID, winfo.UID)
	}
	return false
})
```
*投喂礼物*

```
bili.RegHandleFunc(gobilibili.CmdSendGift, func(c *gobilibili.Context) bool {
	gInfo := c.GetGiftInfo()
	log.Printf("[%d]%s %s 了 %s x %d (价值%.3f)\r\n", c.RoomID, gInfo.Uname, gInfo.Action, gInfo.GiftName, gInfo.Num, float32(gInfo.Price*gInfo.Num)/1000)
	return false
})
```
*在线人数变动*

```
bili.RegHandleFunc(gobilibili.CmdOnlineChange, func(c *gobilibili.Context) bool {
	online := c.GetOnlineNumber()
	log.Printf("[%d]房间里当前在线：%d\r\n", c.RoomID, online)
	return false
})
```
*状态切换为直播开始*

```
bili.RegHandleFunc(gobilibili.CmdLive, func(c *gobilibili.Context) bool {
	online := c.GetOnlineNumber()
	log.Println("主播诈尸啦!")
	return false
})
```
*状态切换为准备中*

```
bili.RegHandleFunc(gobilibili.CmdPreparing, func(c *gobilibili.Context) bool {
	online := c.GetOnlineNumber()
	log.Println("主播正在躺尸")
	return false
})
```
*返回值*
Handler和HandleFunc的返回值用于控制调用链是否继续向下执行。 
如果你希望其它调用链能够继续响应这个事件，请返回false。

## 消息调试
通过注册gobilibili.DebugHandler,可以在收到直播消息时查看原始消息。

```
package main

import "github.com/bigemon/gobilibili"

func main() {
	bili := gobilibili.NewBiliBiliClient()
	bili.RegHandleFunc(gobilibili.CmdAll, gobilibili.DebugHandler)
	bili.ConnectServer(102)
}
```
运行后,当直播间发生事件时,将会输出类似格式的JSON输出:

```
{
  "cmd": "DANMU_MSG",
  "info": [
    [
      0,
      1,
      25,
      16777215,
      1517402685,
      -136720455,
      0,
      "c42d0814",
      0
    ],
    "干嘛不播啦",
    [
      30731115,
      "Ed在",
      0,
      0,
      0,
      10000,
      1,
      ""
    ],
    [],
    [
      1,
      0,
      9868950,
      "\u003e50000"
    ],
    [],
    0,
    0,
    {
      "uname_color": ""
    }
  ]
}
```
以上示例的是一个弹幕消息.
其中`"cmd": "DANMU_MSG"`中的`"DANMU_MSG"`,就是调用`bili.RegHandleFunc`时需要传入的`cmd`参数。
你可以通过`gobilibili.CmdType("嘿,我是CmdType")`,将string转换为CmdType.
在这之后，你可以使用 bili.RegHandleFunc 或 bili.RegHandler 注册这个CmdType.

## 扩展
通过读取gobilibili.Context传入的Msg,可以处理尚未进行支持的事件.
请搭配上一节的消息调试进行食用。
以下是DefaultHandler的实现。

```
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
	if cmd == "DANMU_MSG" {
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
```






