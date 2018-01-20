package util

import (
	"encoding/json"
)

func ToJson(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func ToStruct(jsonStr string, obj *interface{}) {
	json.Unmarshal([]byte(jsonStr), &obj)
}

func Json2map(jsonStr string) (mapObj map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}
