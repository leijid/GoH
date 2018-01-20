package handlers

import (
	"GoH/core/session"
	"net/http"
	"GoH/services"
	"GoH/core/response"
	"GoH/core/constant"
	"GoH/core/util"
)

func GetSmsCode(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	services.GetSmsCode(query["phone"][0])
}

func Login(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ticket, err := services.Login(query["u"][0], query["p"][0], query["c"][0], w, r)
	if err != nil {
		errMap, err := util.Json2map(util.ToJson(err))
		if err != nil {
			response.WriteResponse(w, constant.SystemError, "系统异常", nil)
		} else {
			response.WriteResponse(w, errMap["code"].(string), errMap["message"].(string), nil)
		}
	} else {
		response.WriteResponse(w, constant.Success, "操作成功", ticket)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session.GlobalSessions.SessionDestroy(w, r)
}
