package response

import (
	"GoH/core/util"
	"net/http"
)

type Response struct {
	Code    string      `json:"code"`           //返回码
	Message string      `json:"message"`        //返回信息
	Data    interface{} `json:"data,omitempty"` //返回数据
}

func WriteResponse(w http.ResponseWriter, code string, message string, data interface{}) {
	res := &Response{code, message, data}
	w.Write([]byte(util.ToJson(res)))
}
