package api

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code           int    `json:"code"`
	Hint           string `json:"hint"`
	HttpStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e Error) CompleteError() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return e.Message
	}
	return string(jsonData)
}
