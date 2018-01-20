package models

//omitempty:数据为空时，该字段无需转换成JSON
//`json:"-"`:表示转换JSON时该字段忽略
type UserAuths struct {
	ID            int64  `json:"id"`
	USER_ID       int64  `json:"userId"`
	IDENTITY_TYPE string `json:"identityType,omitempty"`
	IDENTIFIER    string `json:"identifier,omitempty"`
	CREDENTIAL    string `json:"credential,omitempty"`
	FACTOR        string `json:"factor,omitempty"`
}
