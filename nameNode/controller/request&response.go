package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "UserID"

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
		http.StatusBadRequest,
		"文件不存在",
	}
	CodeFileExist ResCode = ResCode{
		http.StatusBadRequest,
		"文件已存在",
	}
	CodeUserInputErr ResCode = ResCode{
		http.StatusForbidden,
		"用户名或密码错误",
	}
	CodeUserExist ResCode = ResCode{
		http.StatusBadRequest,
		"用户名或邮箱已经存在",
	}
	CodeNodeNotFound ResCode = ResCode{
		http.StatusNotFound,
		"文件节点不存在",
	}
	CodeDatabaseForbidden ResCode = ResCode{
		http.StatusBadRequest,
		"数据库配置禁止",
	}
	CodeDiskIsFull ResCode = ResCode{
		http.StatusInternalServerError,
		"文件存储容量已满,请扩容",
	}
)

func (c ResCode) Msg() string {
	return c.ResMsg
}

//错误返回
func ResponseErr(c *gin.Context, code ResCode) {
	c.JSON(code.Code, gin.H{
		"msg": code.Msg(),
	})
}

//正确返回
func ResponseSuc(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
