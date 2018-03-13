package gobilibili

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
)

//Context 消息上下文环境,提供快捷提取消息数据的功能
type Context struct {
	Msg    *simplejson.Json
	RoomID int
}

//DanmuInfo 弹幕信息
type DanmuInfo struct {
	UID         int    `json:"uid"`          //用户ID
	Uname       string `json:"uname"`        //用户名称
	Rank        int    `json:"rank"`         //用户排名
	Level       int    `json:"level"`        //用户等级
	Text        string `json:"text"`         //说的话
	MedalLevel  int    `json:"medal_level"`  //勋章等级
	MedalName   string `json:"medal_name"`   //勋章名称
	MedalAnchor string `json:"medal_anchor"` //勋章所属主播
}

//GetDanmuInfo 在Handler中调用，从simplejson.Json中提取弹幕信息
func (p *Context) GetDanmuInfo() (dInfo DanmuInfo) {
	dInfo.Text, _ = p.Msg.Get("info").GetIndex(1).String()
	dInfo.Uname, _ = p.Msg.Get("info").GetIndex(2).GetIndex(1).String()
	dInfo.UID, _ = p.Msg.Get("info").GetIndex(2).GetIndex(0).Int()
	dInfo.MedalLevel, _ = p.Msg.Get("info").GetIndex(3).GetIndex(0).Int()
	dInfo.MedalName, _ = p.Msg.Get("info").GetIndex(3).GetIndex(1).String()
	dInfo.MedalAnchor, _ = p.Msg.Get("info").GetIndex(3).GetIndex(2).String()
	dInfo.Level, _ = p.Msg.Get("info").GetIndex(4).GetIndex(0).Int()
	dInfo.Rank, _ = p.Msg.Get("info").GetIndex(4).GetIndex(2).Int()
	return
}

//GetOnlineNumber 在Handler中调用，从simplejson.Json中提取房间在线人数
func (p *Context) GetOnlineNumber() int {
	return p.Msg.Get("online").MustInt()
}

//WelcomeGuardInfo 管理进房信息
type WelcomeGuardInfo struct {
	GuardLevel string `json:"guard_level"`
	UID        int    `json:"uid"`
	Username   string `json:"username"`
}

//GetWelcomeGuardInfo 在Handler中调用，从一个simplejson.Json中提取管理进房信息
func (p *Context) GetWelcomeGuardInfo() (wInfo WelcomeGuardInfo) {
	wInfo.GuardLevel = p.Msg.Get("data").Get("guard_level").MustString()
	wInfo.UID = p.Msg.Get("data").Get("uid").MustInt()
	wInfo.Username = p.Msg.Get("data").Get("username").MustString()
	return
}

//WelcomeInfo 普通人员进房信息
type WelcomeInfo struct {
	IsAdmin bool   `json:"is_admin"`
	UID     int    `json:"uid"`
	Uname   string `json:"uname"`
	Vip     int    `json:"vip"`
	Svip    int    `json:"svip"`
}

//GetWelcomeInfo 在Handler中调用，从一个simplejson.Json中提取普通人员进房信息
func (p *Context) GetWelcomeInfo() (wInfo WelcomeInfo) {
	wInfo.IsAdmin = p.Msg.Get("data").Get("is_admin").MustBool() || p.Msg.Get("data").Get("isadmin").MustBool()
	wInfo.UID = p.Msg.Get("data").Get("uid").MustInt()
	wInfo.Uname = p.Msg.Get("data").Get("uname").MustString()
	wInfo.Vip = p.Msg.Get("data").Get("vip").MustInt()
	wInfo.Svip = p.Msg.Get("data").Get("svip").MustInt()
	return
}

//GiftInfo 礼物信息
type GiftInfo struct {
	Action    string `json:"action"`
	AddFollow int    `json:"addFollow"`
	BeatID    string `json:"beatId"`
	BizSource string `json:"biz_source"`
	Capsule   struct {
		Colorful struct {
			Change   int `json:"change"`
			Coin     int `json:"coin"`
			Progress struct {
				Max int `json:"max"`
				Now int `json:"now"`
			} `json:"progress"`
		} `json:"colorful"`
		Normal struct {
			Change   int `json:"change"`
			Coin     int `json:"coin"`
			Progress struct {
				Max int `json:"max"`
				Now int `json:"now"`
			} `json:"progress"`
		} `json:"normal"`
	} `json:"capsule"`
	EventNum   int    `json:"eventNum"`
	EventScore int    `json:"eventScore"`
	GiftID     int    `json:"giftId"`
	GiftName   string `json:"giftName"`
	GiftType   int    `json:"giftType"`
	Gold       int    `json:"gold"`
	// Medal       interface{} `json:"medal"`
	Metadata string `json:"metadata"`
	NewMedal int    `json:"newMedal"`
	NewTitle int    `json:"newTitle"`
	// NoticeMsg   interface{} `json:"notice_msg"`
	Num    int    `json:"num"`
	Price  int    `json:"price"`
	Rcost  int    `json:"rcost"`
	Remain int    `json:"remain"`
	Rnd    string `json:"rnd"`
	Silver int    `json:"silver"`
	// SmalltvMsg  interface{} `json:"smalltv_msg"`
	// SpecialGift interface{} `json:"specialGift"`
	Super     int    `json:"super"`
	Timestamp int    `json:"timestamp"`
	Title     string `json:"title"`
	TopList   *[]struct {
		Face       string `json:"face"`
		GuardLevel int    `json:"guard_level"`
		IsSelf     int    `json:"isSelf"`
		Rank       int    `json:"rank"`
		Score      int    `json:"score"`
		UID        int    `json:"uid"`
		Uname      string `json:"uname"`
	} `json:"top_list"`
	UID   int    `json:"uid"`
	Uname string `json:"uname"`
}

//GetGiftInfo 获取礼物信息
func (p *Context) GetGiftInfo() (gInfo GiftInfo) {
	jbytes, _ := p.Msg.Get("data").Encode()
	jbytes = bytes.Replace(jbytes, []byte(`"beatId":0,`), []byte(`"beatId":"0",`), -1)
	jbytes = bytes.Replace(jbytes, []byte(`"rnd":0,`), []byte(`"rnd":"0",`), -1)
	if err := json.Unmarshal(jbytes, &gInfo); err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(jbytes))
		gInfo.Action = p.Msg.Get("data").Get("action").MustString()
		gInfo.AddFollow = p.Msg.Get("data").Get("addFollow").MustInt()
		gInfo.BeatID = p.Msg.Get("data").Get("beatId").MustString()
		gInfo.BizSource = p.Msg.Get("data").Get("biz_source").MustString()
		gInfo.EventNum = p.Msg.Get("data").Get("eventNum").MustInt()
		gInfo.EventScore = p.Msg.Get("data").Get("eventScore").MustInt()
		gInfo.GiftID = p.Msg.Get("data").Get("giftId").MustInt()
		gInfo.GiftName = p.Msg.Get("data").Get("giftName").MustString()
		gInfo.GiftType = p.Msg.Get("data").Get("giftType").MustInt()
		gInfo.Gold = p.Msg.Get("data").Get("gold").MustInt()
		// gInfo.Medal = p.Msg.Get("data").Get("medal")
		gInfo.Metadata = p.Msg.Get("data").Get("metadata").MustString()
		gInfo.NewMedal = p.Msg.Get("data").Get("newMedal").MustInt()
		gInfo.NewTitle = p.Msg.Get("data").Get("newTitle").MustInt()
		// gInfo.NoticeMsg = p.Msg.Get("data").Get("")
		gInfo.Num = p.Msg.Get("data").Get("num").MustInt()
		gInfo.Price = p.Msg.Get("data").Get("price").MustInt()
		gInfo.Rcost = p.Msg.Get("data").Get("rcost").MustInt()
		gInfo.Remain = p.Msg.Get("data").Get("remain").MustInt()
		gInfo.Rnd = p.Msg.Get("data").Get("rnd").MustString()
		gInfo.Silver = p.Msg.Get("data").Get("silver").MustInt()
		// gInfo.SmalltvMsg = p.Msg.Get("data").Get("")
		// gInfo.SpecialGift = p.Msg.Get("data").Get("")
		gInfo.Super = p.Msg.Get("data").Get("super").MustInt()
		gInfo.Timestamp = p.Msg.Get("data").Get("timestamp").MustInt()
		gInfo.Title = p.Msg.Get("data").Get("title").MustString()
	}
	return
}
