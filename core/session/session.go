package session

import (
	"GoH/core/session/sess"
	_ "GoH/core/session/sess/redis"
	"github.com/dlintw/goconf"
)

var GlobalSessions *session.Manager

/**
 * 使用第三方beego框架的session管理器
 * session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
 * session.NewManager("file", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./tmp"}`)
 * session.NewManager("redis", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"db"}`)
 * redis存储session时，ProviderConfig配置格式为：db
 */
func InitSession(conf *goconf.ConfigFile) {
	sessionStore, _ := conf.GetString("session", "session_store")
	sessionName, _ := conf.GetString("session", "session_name")
	sessionLifeTime, _ := conf.GetInt("session", "seesion_life_time")
	sessionProviderConfig, _ := conf.GetString("session", "session_provider_config")
	sessionIdLength, _ := conf.GetInt("session", "seesion_id_length")
	sessionDisableHTTPOnly, _ := conf.GetBool("session", "session_disable_httponly")
	sessionEnableSetCookie, _ := conf.GetBool("session", "session_enable_set_cookie")
	sessionConfig := &session.ManagerConfig{
		CookieName:      sessionName,
		Gclifetime:      int64(sessionLifeTime),
		ProviderConfig:  sessionProviderConfig,
		SessionIDLength: int64(sessionIdLength),
		DisableHTTPOnly: sessionDisableHTTPOnly,
		EnableSetCookie: sessionEnableSetCookie,
	}
	var err error
	GlobalSessions, err = session.NewManager(sessionStore, sessionConfig)
	if err != nil {
		panic(err)
	}
	go GlobalSessions.GC()
}
