package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "UserID"

var ErrorUserNotLogin = errors.New("用户未登录")
var ErrorMiddlewareParams = errors.New("参数解析错误")

type ResCode struct {
	Code   int
	ResMsg string
}

var (
	CodeSuccess ResCode = ResCode{
		http.StatusOK,
		"success",
	}
	CodeInvalidParam ResCode = ResCode{
		http.StatusBadRequest,
		"请求参数错误",
	}
	CodeServerBusy ResCode = ResCode{
		http.StatusInternalServerError,
		"服务繁忙",
	}
	CodeInvalidToken ResCode = ResCode{
		http.StatusForbidden,
		"无效的token",
	}
	CodeNotFoundFile ResCode = ResCode{
		http.StatusNotFound,
		"文件不存在",
	}
	CodeFileExist ResCode = ResCode{
		http.StatusBadRequest,
		"文件已存在",
	}
	CodeUploadIDNotFound ResCode = ResCode{
		http.StatusBadRequest,
		"UploadID 不存在",
	}
	CodeHashNotMatch ResCode = ResCode{
		http.StatusBadRequest,
		"文件数据hash校验失败",
	}
	CodeFileNameInvalid ResCode = ResCode{
		http.StatusBadRequest,
		"文件名非法",
	}
	CodeFileError ResCode = ResCode{
		http.StatusInternalServerError,
		"文件已损坏",
	}
)

func (c ResCode) Msg() string {
	return c.ResMsg
}

// GetCurrentUser 获取当前用户ID
func GetCurrentUser(c *gin.Context) (UserID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}

	UserID, ok = uid.(int64)
	if !ok {
		return 0, ErrorMiddlewareParams
	}
	return UserID, nil
}

//cos错误返回
func ResponseErr(c *gin.Context, code ResCode) {
	c.JSON(code.Code, gin.H{
		"msg": code.Msg(),
	})
}

//cos正确返回
func ResponseSuc(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  CodeSuccess.Msg(),
		"data": data,
	})
}
