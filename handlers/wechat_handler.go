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
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type AccessToken struct {
	AppId       string
	Secret      string
	AccessToken string
	LifeTime    float64
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
		return
	}
	if t, ok := j["access_token"]; ok {
		a.AccessToken = t.(string)
		fmt.Println("update token", a.AccessToken)
	}
	if e, ok := j["expires_in"]; ok {
		a.LifeTime = e.(float64)
	}
}

func TokenHandle(a *AccessToken) {
	for {
		if a.LifeTime > 10 {
			time.Sleep(time.Second * 5)
			a.LifeTime -= 5
		} else {
			a.UpdateToken()
			time.Sleep(time.Second * 1)
		}
	}
}

func TokenRoutineStart() {
	accessToken = AccessToken{
		AppId:       config.GetConfig().Wechat.AppId,
		Secret:      config.GetConfig().Wechat.AppSecret,
		AccessToken: "",
		LifeTime:    0,
	}
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
		"]]></FromUserName><CreateTime>" + r.CreateTime + "</CreateTime><MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[" +
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
		fmt.Println("read media id file error", err)
		return "", err
	}
	lines := strings.Split(string(file), "\n")
	// CRLF LF!!
	mediaId := lines[rand.Intn(len(lines))]
	return mediaId, nil
}

func GetTuringRobotResponse(content string) string {
	m := map[string]string{
		"key":    config.GetConfig().Turing.TuringKey,
		"info":   content,
		"userid": "wechat-robot",
	}
	bytesData, _ := json.Marshal(m)
	data := bytes.NewBuffer(bytesData)
	res, err := http.Post(config.GetConfig().Turing.TuringApiUrl, "application/json;charset=utf-8", data)
	if err != nil {
		fmt.Println("http post error", err)
		return "本喵听不懂你在说什么o(=•ェ•=)m😕"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read resp.Body failed", err)
		return "本喵听不懂你在说什么o(=•ェ•=)m😕"
	}

	var j = make(map[string]interface{})

	err = json.Unmarshal(body, &j)
	if err != nil {
		fmt.Println("json unmarshal failed", err)
		return "本喵听不懂你在说什么o(=•ェ•=)m😕"
	}

	code := j["code"].(float64)
	if code == 100000 {
		// text
		return j["text"].(string)
	} else if code == 200000 {
		// pic
		return j["url"].(string)
	} else if code == 302000 {
		// news
		return j["list"].(string)
	} else if code == 308000 {
		// dish
		return j["list"].(string)
	} else {
		return "本喵有点累了要休息了，明天再找我玩吧o(=•ェ•=)m😕！"
	}
}

func DownloadPicFromUrl(url string, savePath string) string {
	sl := strings.Split(url, "//")[1]
	sl = strings.Split(sl, "/")[1]
	ext := strings.Split(sl, "_")[1]
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error", err)
		return ""
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read body error", err)
		return ""
	}
	fileName := time.Now().Format("2006-01-02 15:03:04") + "." + ext
	err = ioutil.WriteFile(savePath+fileName, data, 0644)
	if err != nil {
		fmt.Println("save pic error", err)
		return ""
	}
	fmt.Println("download pic from url")
	return fileName
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
	url := "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=" + accessToken.AccessToken + "&type=image"
	fmt.Println("start upload permanent material image")
	//file, _ := ioutil.ReadFile(config.GetConfig().File.TempFilePath + filename)
	//fmt.Println("FILE:", file[:10], len(file))
	//
	//res, err := http.Post(url, "multipart/form-data", bytes.NewReader(file))
	//if err != nil {
	//	fmt.Println("post material image error", err)
	//	return ""
	//}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)

	fileWriter, err := writer.CreateFormFile("media", config.GetConfig().File.TempFilePath+filename)
	if err != nil {
		fmt.Println("file write to buffer error", err)
		return ""
	}
	fh, err := os.Open(config.GetConfig().File.TempFilePath + filename)
	if err != nil {
		fmt.Println("open file error", err)
		return ""
	}

	_, _ = io.Copy(fileWriter, fh)

	writer.Close()
	contentType := writer.FormDataContentType()
	res, err := http.Post(url, contentType, buff)
	if err != nil {
		fmt.Println("http post error", err)
		return ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read resp.Body failed", err)
		return ""
	}

	var r = make(map[string]interface{})
	//fmt.Println("Body:", string(body))
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshal failed", err)
		return ""
	}
	mediaId := r["media_id"].(string)
	fmt.Println("upload permanent material image:", mediaId)
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

		replyContent := GetTuringRobotResponse(msg.Content)
		fmt.Println("Turing robot reply:", replyContent)
		replyMsg := ReplyTextMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			Content: replyContent,
		}
		c.String(200, replyMsg.ToXml())
		return
	} else if msgType == "image" {
		// image
		msg := RecImageMsg{}
		_ = xml.Unmarshal(body, &msg)

		fmt.Println("[Image] From:", msg.FromUserName, msg.PicUrl)
		savePath := config.GetConfig().File.TempFilePath
		picName := DownloadPicFromUrl(msg.PicUrl, savePath)
		mediaId := UploadPermanentMaterialImage(picName)
		if mediaId != "" {
			SaveMediaId(mediaId)
		}
		randomMediaId, _ := GetRandomMediaId()
		fmt.Println("random media id:", randomMediaId)
		replyMsg := ReplyImageMsg{
			ReplyMsg: ReplyMsg{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				CreateTime:   time.Now().Format("2006-01-02 15:03:04"),
			},
			MediaId: randomMediaId,
		}
		//fmt.Println(replyMsg.ToXml())
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
		return
	} else {
		c.String(200, "本喵不明白你在说什么🐱")
		return
	}
}
