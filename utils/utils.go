package utils

import (
	"encoding/json"
	"os"
)

func ReadJSONFile(path string, v interface{}) error {
	filePtr, err := os.Open(path)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}
