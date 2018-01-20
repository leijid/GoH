package services

import (
	"GoH/core/mysql"
	"fmt"
	"GoH/core/util"
	herr "GoH/core/error"
	"GoH/core/session"
	"net/http"
	"GoH/core/constant"
	hredis "GoH/core/redis"
	"github.com/garyburd/redigo/redis"
	"database/sql"
)

type SessionAccount struct {
	UserId   int64    `json:"userId"`
	UserName string   `json:"userName"`
	NickName string   `json:"nickName"`
	UserRole []string `json:"userRole"`
}

func GetSmsCode(phone string) {
	c := hredis.RedisClient.Get()
	defer c.Close()
	c.Do("SELECT", constant.GohUserRedisDb)
	code := util.RandomCreateBytes(6, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
	c.Do("SETEX", fmt.Sprintf("%s%s", constant.SmsCodePrefix, phone), constant.SmsCodeExpireTime, code)
	//后续短信发送功能
}

func Login(loginName string, password string, code string, w http.ResponseWriter, r *http.Request) (string, error) {
	db := mysql.OpenMysql()
	defer db.Close()
	c := hredis.RedisClient.Get()
	defer c.Close()
	c.Do("SELECT", constant.GohUserRedisDb)
	mysql.Orm.SetTable("z_user_auths")
	userAuth := mysql.Orm.Where(fmt.Sprintf("IDENTIFIER='%s'", loginName)).FindOne(db)
	if len(userAuth) > 0 {
		inputPassEncrypt := util.Md5(fmt.Sprintf("%s%s", password, userAuth[0]["FACTOR"]))
		if userAuth[0]["CREDENTIAL"] == inputPassEncrypt {
			kvi, _ := redis.Int64(c.Do("GET", fmt.Sprintf("%s%s", constant.LockAccountPrefix, loginName)))
			if kvi >= constant.LockTimes {
				//判断账号是否超出限制，超出不允许登录
				return "", herr.OutError(constant.LoginAccountLock, "账号已被冻结，请24小时后再试！")
			} else {
				mysql.Orm.SetTable("z_user")
				user := mysql.Orm.Where(fmt.Sprintf("ID=%d", userAuth[0]["USER_ID"])).FindOne(db)
				if len(user) > 0 {
					setSessionAccount(user, db, w, r)
					//登录成功删除冻结标识
					c.Do("DEL", fmt.Sprintf("%s%s", constant.LockAccountPrefix, loginName))
					//登录成功返回ticket,并缓存ticket和用户信息之间的关系
					ticketCode := util.RandomCreateBytes(45)
					ticket := fmt.Sprintf("%s%s", "TICKET_", string(ticketCode))
					c.Do("SETEX", ticket, constant.LockExpireTime, util.ToJson(user[0]))
					return string(ticket), nil
				}
			}
		} else {
			//记录错误次数，超过限制锁定账号
			kvi, _ := redis.Int64(c.Do("GET", fmt.Sprintf("%s%s", constant.LockAccountPrefix, loginName)))
			c.Do("SETEX", fmt.Sprintf("%s%s", constant.LockAccountPrefix, loginName), constant.LockExpireTime, kvi+1)
			return "", herr.OutError(constant.LoginPasswordError, "登录名或密码错误！")
		}
	} else {
		return "", herr.OutError(constant.LoginPasswordError, "登录名或密码错误！")
	}
	return "", nil
}

func setSessionAccount(user map[int]map[string]interface{}, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	sess, _ := session.GlobalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	acct := &SessionAccount{}
	acct.UserId = user[0]["ID"].(int64)
	acct.UserName = user[0]["USER_NAME"].(string)
	acct.NickName = user[0]["NICK_NAME"].(string)
	mysql.Orm.SetTable("z_user_role as ur")
	userRole := mysql.Orm.LeftJoin("z_role as r", "ur.ROLE_ID=r.ID").Where(fmt.Sprintf("ur.USER_ID=%d", user[0]["ID"])).NoLimit().FindAll(db)
	role := make([]string, len(userRole))
	for k, v := range userRole {
		role[k] = v["ROLE_CODE"].(string)
	}
	acct.UserRole = role
	sess.Set(constant.GohSessionName, util.ToJson(*acct))
}
