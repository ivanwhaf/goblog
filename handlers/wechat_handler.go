package handlers

import (
	"github.com/gin-gonic/gin"
	"goblog/services"
)

// WechatHandler / POST
func WechatHandler(c *gin.Context) {
	xmlMsg, body, err := services.ParsePostMsg(c)
	if err != nil {
		c.String(200, "本喵不明白你在说什么🐱")
		return
	}
	msgType := xmlMsg.MsgType
	if msgType == "text" {
		// text
		services.HandleTextMsg(c, body)
	} else if msgType == "image" {
		// image
		services.HandleImageMsg(c, body)
	} else if msgType == "voice" {
		// voice
		services.HandleVoiceMsg(c, body)
	} else if msgType == "event" {
		// event
		services.HandleEventMsg(c, body)
	} else {
		// other
		services.HandleDefaultMsg(c)
	}
}
