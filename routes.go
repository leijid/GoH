package main

import (
	"GoH/handlers"
	"net/http"
)

type Route struct {
	Name        string           //路由名称
	Method      string           //请求方式
	Pattern     string           //请求表达式
	HandlerFunc http.HandlerFunc //请求处理方法
	Auth        bool             //是否鉴权
	Role        []string         //用户角色
}

type Routes []Route

var routes Routes = Routes{
	Route{Name: "Index", Method: "GET", Pattern: "/", HandlerFunc: handlers.Index, Auth: true, Role: []string{"ADMIN", "USER"}},
	Route{Name: "Login", Method: "GET", Pattern: "/login", HandlerFunc: handlers.Login, Auth: false},
	Route{Name: "GetUser", Method: "GET", Pattern: "/get", HandlerFunc: handlers.GetUser, Auth: false},
	Route{Name: "Logout", Method: "GET", Pattern: "/logout", HandlerFunc: handlers.Logout, Auth: false},
	Route{Name: "GetUserList", Method: "GET", Pattern: "/get_list", HandlerFunc: handlers.GetUserList, Auth: false},
	Route{Name: "GetSmsCode", Method: "GET", Pattern: "/get_sms_code", HandlerFunc: handlers.GetSmsCode, Auth: false},
	Route{Name: "Update", Method: "GET", Pattern: "/update", HandlerFunc: handlers.UpdateUser, Auth: false},
	Route{Name: "GetSessionUser", Method: "GET", Pattern: "/get_sess_user", HandlerFunc: handlers.GetSessionUser, Auth: false},
}
