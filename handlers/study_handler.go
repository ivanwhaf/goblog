package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/core"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"time"
)

// StudyHandler /study GET
func StudyHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "study.html", gin.H{})
}

// ApiStudyDictionaryRetrievalHandler /api/study/dictionary/retrieval GET
func ApiStudyDictionaryRetrievalHandler(c *gin.Context) {
	key := c.DefaultPostForm("key", "")
	limitStr := c.DefaultPostForm("limit", "5")
	var limit int64
	limit = util.StringToInt64(limitStr)
	dictionaries, _ := stores.DictionaryStore.GetRelativeWordsByKey(key, limit)
	c.JSON(200, gin.H{
		"words": dictionaries,
	})
}

// ApiStudyDictionaryHandler /api/study/dictionary POST
func ApiStudyDictionaryHandler(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("username") == nil {
		c.String(200, "-1")
		return
	}
	key := c.DefaultPostForm("key", "")
	value := c.DefaultPostForm("value", "")
	keyType := c.DefaultPostForm("key_type", "")
	valueType := c.DefaultPostForm("value_type", "")
	_ = stores.DictionaryStore.AddDictionary(&core.Dictionary{
		Key_:       key,
		Value_:     value,
		KeyType:    keyType,
		ValueType:  valueType,
		CreateDate: time.Now(),
	})
	c.String(200, "1")
}
