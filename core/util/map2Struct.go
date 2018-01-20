package util

import (
	"github.com/mitchellh/mapstructure"
)

func Map2Struct(data map[string]interface{}, obj interface{}) error {
	return mapstructure.Decode(data, obj)
}
