package enum

import "github.com/gin-gonic/gin"

type ResponseStatus struct {
	success bool
	code    int
	msg     string
	data    interface{}
}

var Success = ResponseStatus{true, 200, "", nil}
var ServerError = ResponseStatus{false, 500, "服务器内部错误", nil}
var ParamError = ResponseStatus{false, 400, "参数错误", nil}
var NullParam = ResponseStatus{false, 1001, "空参数", nil}
var UploadError = ResponseStatus{false, 1002, "文件上传错误", nil}
var TotalError = ResponseStatus{false, 1003, "错误", nil}

func (r *ResponseStatus) ToFailGinH() (maps gin.H) {
	maps = make(map[string]interface{})
	maps["success"] = r.success
	maps["code"] = r.code
	maps["msg"] = r.msg
	return
}

func (r *ResponseStatus) ToFailGinHWithMsg(msg string) (maps gin.H) {
	maps = make(map[string]interface{})
	maps["success"] = r.success
	maps["code"] = r.code
	maps["msg"] = msg
	return
}

func (r *ResponseStatus) ToSuccessGinH() (maps gin.H) {
	maps = make(map[string]interface{})
	maps["success"] = r.success
	maps["code"] = r.code
	maps["msg"] = r.msg
	return
}

func (r *ResponseStatus) ToSuccessGinHWithMsgAndData(msg string, data interface{}) (maps gin.H) {
	maps = make(map[string]interface{})
	maps["success"] = r.success
	maps["code"] = r.code
	maps["msg"] = msg
	maps["data"] = data
	return
}
