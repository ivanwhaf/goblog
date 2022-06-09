package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type AccessToken struct {
	AppId       string
	Secret      string
	AccessToken string
	LifeTime    int64
}

var accessToken AccessToken

func (a *AccessToken) UpdateToken() {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + a.AppId + "&secret=" + a.Secret
	res, _ := http.Get(url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var j = make(map[string]interface{})
	err := json.Unmarshal(body, &j)
	if err != nil {
		fmt.Println("json unmarshal error")
	}
	a.AccessToken = j["access_token"].(string)
	a.LifeTime = j["expires_in"].(int64)
}

func TokenHandle(a *AccessToken) {
	for {
		if a.LifeTime > 10 {
			time.Sleep(time.Second * 5)
			a.LifeTime -= 5
		} else {
			a.UpdateToken()
			fmt.Println("update token")
		}
	}
}

func TokenRoutineStart() {
	accessToken = AccessToken{}
	go TokenHandle(&accessToken)
}

type RecMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	MsgId        string
}

type RecTextMsg struct {
	RecMsg
	Content string
}

type RecImageMsg struct {
	RecMsg
	PicUrl  string
	MediaId string
}

type RecVoiceMsg struct {
	RecMsg
	MediaId     string
	Format      string
	Recognition string
}

type ReplyMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
}

type ReplyTextMsg struct {
	ReplyMsg
	Content string
}

func (r *ReplyTextMsg) ToXml() string {
	resp := "<xml><ToUserName><![CDATA[" + r.ToUserName + "]]></ToUserName><FromUserName><![CDATA[" + r.FromUserName +
		"]]></FromUserName><CreateTime>" + r.CreateTime + "</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[" +
		r.Content + "]]></Content></xml>"
	return resp
}

type ReplyImageMsg struct {
	ReplyMsg
	MediaId string
}

func (r *ReplyImageMsg) ToXml() string {
	resp := "<xml><ToUserName><![CDATA[" + r.ToUserName + "]]></ToUserName><FromUserName><![CDATA[" + r.FromUserName +
		"]]></FromUserName><CreateTime>" + r.CreateTime +
		"</CreateTime><MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[" +
		r.MediaId + "]]></MediaId></Image></xml>"
	return resp
}

type RecEventMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Event        string
}

func GetRandomMediaId() (string, error) {
	fileName := config.GetConfig().File.TempFilePath + "media_id.txt"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(file), "\n")
	mediaId := lines[rand.Intn(len(lines))]
	return mediaId, nil
}

func GetRobotResponse(content string) string {
	m := map[string]string{
		"key":    config.GetConfig().Turing.TuringKey,
		"info":   content,
		"userid": "wechat-robot",
	}
	bytesData, _ := json.Marshal(m)
	data := bytes.NewBuffer(bytesData)
	res, _ := http.Post(config.GetConfig().Turing.TuringApiUrl, "application/json;charset=utf-8", data)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read resp.Body failed")
		return ""
	}

	var r = make(map[string]interface{})
	fmt.Println("Body:", string(body))
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return ""
	}

	code := r["code"].(float64)
	if code == 100000 {
		return r["text"].(string)
	} else if code == 200000 {
		return r["url"].(string)
	} else if code == 302000 {
		return r["list"].(string)
	} else if code == 308000 {
		return r["list"].(string)
	} else {
		return "本喵有点累了要休息了，明天再找我玩吧o(=•ェ•=)m😕！"
	}
}

func DownloadPicFromUrl(url string) string {
	sl := strings.Split(url, "//")[1]
	sl = strings.Split(sl, "/")[1]
	ext := strings.Split(sl, "_")[1]
	fmt.Println(ext)
	return ext
}

func SaveMediaId(mediaId string) {
	file, err := os.OpenFile(config.GetConfig().File.TempFilePath+"media_id.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open media_id.txt failed")
	}
	defer file.Close()
	_, err = io.WriteString(file, mediaId+"\n")
	if err != nil {
		fmt.Println("write media id failed")
	}
}

func UploadPermanentMaterialImage(filename string) string {
	accessToken := ""
	url := "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=" + accessToken + "&type=image"
	fmt.Println("start upload permanent material image")
	file, _ := ioutil.ReadFile(config.GetConfig().File.TempFilePath + filename)
	res, err := http.Post(url, "multipart/form-data", bytes.NewReader(file))
	if err != nil {
		fmt.Println("post material image error")
		return ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read resp.Body failed")
		return ""
	}

	var r = make(map[string]string)
	fmt.Println("Body:", string(body))
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return ""
	}
	mediaId := r["media_id"]
	return mediaId
}

// WechatHandler / POST
func WechatHandler(c *gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)

	xmlMsg := RecMsg{}
	_ = xml.Unmarshal(body, &xmlMsg)

	if len(string(body)) == 0 {
		c.String(200, "本喵不明白你在说什么🐱")
		return
	}

	fmt.Println(xmlMsg.MsgType)

	msgType := xmlMsg.MsgType
	if msgType == "text" {
		// text
		msg := RecTextMsg{}
		_ = xml.Unmarshal(body, &msg)
		fmt.Println("[Text] From:", msg.FromUserName, msg.Content)

		if msg.Content == "【收到不支持的消息类型，暂无法显示】" {
			mediaId, _ := GetRandomMediaId()
			replyMsg := ReplyImageMsg{
				ReplyMsg: ReplyMsg{
					ToUserName:   msg.FromUserName,
					FromUserName: msg.ToUserName,
					CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
				},
				MediaId: mediaId,
			}
			fmt.Println(replyMsg.ToXml())
			c.String(200, replyMsg.ToXml())
			return
		}

		replyContent := GetRobotResponse(msg.Content)
		fmt.Println("Turing robot reply:", replyContent)
		replyMsg := ReplyTextMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			Content: replyContent,
		}
		fmt.Println(replyMsg.ToXml())
		c.String(200, replyMsg.ToXml())
		return
	} else if msgType == "image" {
		// image
		msg := RecImageMsg{}
		_ = xml.Unmarshal(body, &msg)

		fmt.Println("[Image] From:", msg.FromUserName, msg.PicUrl)
		picType := DownloadPicFromUrl(msg.PicUrl)
		mediaId := UploadPermanentMaterialImage(picType)
		if mediaId != "" {
			SaveMediaId(mediaId)
		}
		randomMediaId, _ := GetRandomMediaId()
		replyMsg := ReplyImageMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			MediaId: randomMediaId,
		}
		fmt.Println(replyMsg.ToXml())
		c.String(200, replyMsg.ToXml())
		return
	} else if msgType == "voice" {
		// voice
		msg := RecVoiceMsg{}
		_ = xml.Unmarshal(body, &msg)
		fmt.Println("[Voice] From:", msg.FromUserName, msg.MediaId)

		recognition := msg.Recognition
		replyMsg := ReplyTextMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			Content: recognition,
		}
		fmt.Println(replyMsg.ToXml())
		c.String(200, replyMsg.ToXml())
		return
	} else if msgType == "event" {
		// event
		msg := RecEventMsg{}
		_ = xml.Unmarshal(body, &msg)
		fmt.Println("[Event] From:", msg.FromUserName, msg.Event)
		content := "欢迎关注嗷大喵的公众号！本公众号提供各种实用功能！请准备好您的膝盖\\r\\n你可以试着发文字、图片和我聊天，我知道好多东西（天气，快递，笑话，翻译，算数，食谱，车票，语音识别）~~"
		replyMsg := ReplyTextMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			Content: content,
		}
		c.String(200, replyMsg.ToXml())
	} else {
		c.String(200, "本喵不明白你在说什么🐱")
		return
	}
}
