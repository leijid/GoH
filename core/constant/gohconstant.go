package constant

const Debug = true

//session key name
const GohSessionName = "goh_session_account" //session键值
const GohUserRedisDb = 1                     //redis用户库

//response code
const Success = "0000"             //成功
const SystemForbidden = "403"      //禁止访问
const SystemError = "500"          //系统错误
const LoginPasswordError = "10001" //账号或密码错误
const LoginAccountLock = "10002"   //账号锁定

//lock
const LockTimes = 5                 //输入错误次数锁定账号
const LockExpireTime = 24 * 60 * 60 //锁定时间
const LockAccountPrefix = "lock_"   //锁定账号的redis键值前缀

//sms code
const SmsCodeExpireTime = 5 * 60 //短信验证码超时时间，单位:s
const SmsCodePrefix = "sms_"     //短信验证码的redis键值前缀

//http request
const HttpTimeOut = 30 //http请求超时时间
