package error

import (
	"encoding/json"
)

type GError struct {
	Code    string `json:"code"`    //错误码
	Message string `json:"message"` //错误信息
}

func (e GError) Error() string {
	jsonError, _ := json.Marshal(e)
	return string(jsonError)
}

func OutError(code string, messsage string) error {
	return GError{
		code,
		messsage,
	}
}
