package main

import (
	"GoH/core/log"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"GoH/core/session"
	"GoH/core/util"
	"GoH/core/response"
	"GoH/core/constant"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name, route.Auth, route.Role)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}

func Logger(inner http.Handler, name string, auth bool, role []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if auth {
			sess, _ := session.GlobalSessions.SessionStart(w, r)
			defer sess.SessionRelease(w)
			sessionStr := sess.Get(constant.GohSessionName)
			if sessionStr == nil {
				forbidden(w, "用户未登录")
				return
			} else {
				sessionMap, err := util.Json2map(sessionStr.(string))
				if err != nil {
					forbidden(w, "用户未登录")
					return
				}
				userRole := sessionMap["userRole"]
				isAccess := false
				for _, ur := range userRole.([]interface{}) {
					for _, r := range role {
						if ur.(string) == r {
							isAccess = true
						}
					}
				}
				if !isAccess {
					forbidden(w, "权限不足")
					return
				}
			}
		}
		inner.ServeHTTP(w, r)
		log.Infof(
			"%s\t%s\t%s\tcost:%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func forbidden(w http.ResponseWriter, message string) {
	log.Warningf("403 Forbidden")
	response.WriteResponse(w, constant.SystemForbidden, "未登录或权限不足", nil)
}
