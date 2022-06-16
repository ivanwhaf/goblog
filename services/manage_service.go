package services

import (
	"goblog/core"
	"goblog/stores"
	"strconv"
)

func GetCommentsSlcForManage(comments []*core.Comment) []map[string]string {
	var commentsSlc []map[string]string
	for _, comment := range comments {
		m := map[string]string{
			"Id":          strconv.FormatInt(comment.Id, 10),
			"ArticleId":   strconv.FormatInt(comment.ArticleId, 10),
			"ReplyName":   comment.ReplyName,
			"Content":     comment.Content,
			"CommentDate": comment.CommentDate.Format("2006-01-02 15:03:04"),
			"Ip":          comment.Ip,
			"Location":    comment.Location,
		}
		commentsSlc = append(commentsSlc, m)
	}
	return commentsSlc
}

func GetAdminsSlcForManage(admins []*core.Admin) []map[string]string {
	var adminsSlc []map[string]string
	for _, admin := range admins {
		m := map[string]string{
			"id":        strconv.FormatInt(admin.Id, 10),
			"username":  admin.Username,
			"password":  admin.Password,
			"nickname":  admin.Nickname,
			"sex":       admin.Sex,
			"authority": strconv.Itoa(int(admin.Authority)),
		}
		latestLogin, _ := stores.LoginStore.GetLatestLoginByUsername(admin.Username)
		if latestLogin.Ip != "" {
			m["latestLoginDate"] = latestLogin.LoginDate.Format("2006-01-02 15:03:04")
		}
		adminsSlc = append(adminsSlc, m)
	}
	return adminsSlc
}
