package utils

import (
	"encoding/json"
	"fmt"
)

func MapToStruct(dataMap interface{}, response interface{}) error {
	b, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Errorf("cant json encode data map: %s", err.Error())
	}
	err = json.Unmarshal(b, response)
	if err != nil {
		return fmt.Errorf("cant encode dataMap: %s", err.Error())
	}
	return nil
}
