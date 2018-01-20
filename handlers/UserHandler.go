package handlers

import (
	"GoH/core/mysql"
	"GoH/core/util"
	"net/http"
	"GoH/models"
	"GoH/core/response"
	"GoH/core/constant"
	"GoH/core/session"
	"fmt"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	db := mysql.OpenMysql()
	defer db.Close()
	mysql.Orm.SetTable("z_user")
	data := mysql.Orm.FindOne(db)
	user := &models.User{}
	userData, _ := data[0]
	util.Map2Struct(userData, user)
	response.WriteResponse(w, constant.Success, "操作成功", user)
}

func GetSessionUser(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.GlobalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	user := sess.Get(constant.GohSessionName).(string)
	userMap, _ := util.Json2map(user)
	response.WriteResponse(w, constant.Success, "操作成功", userMap)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := mysql.OpenMysql()
	defer db.Close()
	mysql.Orm.SetTable("z_user")
	data := mysql.Orm.FindAll(db)
	for k, v := range data {
		mysql.Orm.SetTable("z_user_auths")
		auth := mysql.Orm.Where(fmt.Sprintf("IDENTIFIER='%s'", v["USER_CODE"])).FindOne(db)
		if len(auth) == 0 {
			fmt.Printf("%s\n", v["USER_CODE"])
			param := make(map[string]interface{})
			param["ID"] = fmt.Sprintf("%d", int64(k+3000000))
			param["USER_ID"] = fmt.Sprintf("%d", v["ID"])
			param["IDENTITY_TYPE"] = "3"
			param["IDENTIFIER"] = v["USER_CODE"]
			param["CREDENTIAL"] = util.Md5(fmt.Sprintf("%s%s", "123456", "123456"))
			param["FACTOR"] = "123456"
			mysql.Orm.Insert(db, param)
		}
	}
}

func GetUserList(w http.ResponseWriter, r *http.Request) {
	db := mysql.OpenMysql()
	defer db.Close()
	mysql.Orm.SetTable("z_user")
	data := mysql.Orm.FindAll(db)
	list := make([]models.User, len(data))
	for k, v := range data {
		user := &models.User{}
		util.Map2Struct(v, user)
		list[k] = *user
	}
	response.WriteResponse(w, constant.Success, "操作成功", list)
}
