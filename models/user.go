package models

import "time"

//omitempty:数据为空时，该字段无需转换成JSON
//`json:"-"`:表示转换JSON时该字段忽略
type User struct {
	ID            int64     `json:"id"`
	USER_CODE     string    `json:"userCode,omitempty"`
	USER_NAME     string    `json:"userName,omitempty"`
	NICK_NAME     string    `json:"nickName,omitempty"`
	REAL_NAME     string    `json:"realName,omitempty"`
	USER_TYPE     string    `json:"userType,omitempty"`
	MOBILE_PHONE  string    `json:"mobilePhone,omitempty"`
	FIX_PHONE     string    `json:"fixPhone,omitempty"`
	EMAIL         string    `json:"email,omitempty"`
	ID_TYPE       string    `json:"idType,omitempty"`
	ID_CARD       string    `json:"idCard,omitempty"`
	BIRTHDAY      time.Time `json:"birthday"`
	SEX           string    `json:"sex,omitempty"`
	AGE           int32     `json:"age"`
	HEAD_IMG      string    `json:"headImg,omitempty"`
	QQ            string    `json:"qq,omitempty"`
	WECHAT        string    `json:"wechat,omitempty"`
	WEIBO         string    `json:"weibo,omitempty"`
	ALIPAY        string    `json:"alipay,omitempty"`
	IS_NEED_EMAIL int8      `json:"isNeedEmail"`
	APP_CODE      string    `json:"appCode,omitempty"`
	REGISTER_IP   string    `json:"registerIp,omitempty"`
	REGISTER_TIME time.Time `json:"registerTime"`
	REGISTER_TYPE string    `json:"registerType,omitempty"`
	REFEREE       string    `json:"referee,omitempty"`
	ORIGIN_TYPE   string    `json:"originType,omitempty"`
	PWD_STRENGTH  string    `json:"pwdStrength,omitempty"`
	SOURCE        string    `json:"source,omitempty"`
	STATUS        int32     `json:"status"`
	ONLINE_STATUS string    `json:"onlineStatus,omitempty"`
}
