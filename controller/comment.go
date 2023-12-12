package controller

import (
	"ezgin/model"
	"ezgin/services"
	"ezgin/utils"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func PostComment(c *gin.Context) {
	uid, ok := utils.GetUid(c)
	if !ok {
		RespFailed(c, CodeNeedLogin)
		return
	}
	//获取参数并校验
	ParamComment := new(model.ParamComment)
	if err := c.ShouldBind(ParamComment); err != nil {
		RespFailed(c, CodeInvalidParam)
		return
	}
	if ParamComment.FromUID != uid {
		RespFailed(c, CodeInvalidUser)
		return
	}
	if ParamComment.FromUID == 0 || ParamComment.ToUID == 0 || ParamComment.Message == "" {
		RespFailed(c, CodeInvalidParam)
		return
	}
	//根据错误类型返回响应
	mid, err := services.PostComment(ParamComment)
	if err != nil {
		RespFailed(c, CodeServiceBusy)
		log.Printf("%v", err)
		return
	}
	RespSuccess(c, &model.Comment{ID: mid})
}

func GetComment(c *gin.Context) {
	id := c.Param("uid")
	//string to int
	uid, err := strconv.Atoi(id)
	if err != nil {
		RespFailed(c, CodeInvalidParam)
		return
	}
	commentInfo, err := services.GetComment(uid)
	if err != nil {
		RespFailed(c, CodeServiceBusy)
		log.Printf("%v", err)
		return
	}
	RespSuccess(c, commentInfo)
}

func DeleteComment(c *gin.Context) {
	uid, ok := utils.GetUid(c)
	if !ok {
		RespFailed(c, CodeNeedLogin)
		return
	}

	id := c.Param("mid")
	//string to int
	mid, err := strconv.ParseInt(id, 10, 64)

	if uidFrom, err := services.GetCommentFrom(mid); err != nil || uidFrom != uid {
		RespFailed(c, CodeInvalidUser)
		return
	}

	if err != nil {
		RespFailed(c, CodeInvalidParam)
		return
	}
	err = services.DeleteComment(mid)
	if err != nil {
		RespFailed(c, CodeServiceBusy)
		log.Printf("%v", err)
		return
	}
	RespSuccess(c, nil)
}
