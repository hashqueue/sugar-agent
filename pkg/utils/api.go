package utils

import (
	"encoding/json"
	"errors"
)

func UserLogin(baseUrl string, reqData map[string]interface{}) (string, error) {
	client := &HTTPClient{
		BaseURL: baseUrl,
	}
	resp, err := client.Post("/api/v1/system/users/login/", map[string]string{
		"Content-Type": "application/json",
	}, reqData)
	FailOnError(err, "Failed to send request")
	data := make(map[string]interface{})
	err = json.Unmarshal(resp, &data)
	FailOnError(err, "Failed to unmarshal response body")
	if data["code"].(float64) == 20000 && data["message"] == "登录成功" {
		// convert token to string
		return data["data"].(map[string]interface{})["access"].(string), nil
	}
	return "", errors.New("failed to login")
}

func UpdateTaskStatus(baseUrl string, reqData map[string]interface{}, taskUUID string, token string) error {
	client := &HTTPClient{
		BaseURL: baseUrl,
	}
	resp, err := client.Patch("/api/v1/task-results/"+taskUUID+"/", map[string]string{
		"Content-Type":  "application/json",
		"Authorization": `Bearer ` + token,
	}, reqData)
	FailOnError(err, "Failed to send request")
	data := make(map[string]interface{})
	err = json.Unmarshal(resp, &data)
	FailOnError(err, "Failed to unmarshal response body")
	if data["code"].(float64) == 20000 && data["message"] == "success" {
		return nil
	}
	return errors.New("failed to update task status")
}
