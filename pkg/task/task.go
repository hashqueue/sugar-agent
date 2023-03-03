package task

import (
	"encoding/json"
	"os"

	"sugar-agent/internal"
	"sugar-agent/pkg/utils"
)

type Metadata struct {
	DurationTime uint64 `json:"durationTime"`
	Count        uint64 `json:"count"`
}

type MQMessage struct {
	TaskType int      `json:"taskType"`
	Metadata Metadata `json:"metadata"`
}

func StartTask(msg []byte) {
	var mqMessage MQMessage
	err := json.Unmarshal(msg, &mqMessage)
	utils.LogOnError(err, "Failed to unmarshal message")
	if mqMessage.TaskType == 0 {
		//fmt.Println("TaskType == 0", mqMessage)
		perfData := internal.StartGetPerfDataTask(mqMessage.Metadata.DurationTime, mqMessage.Metadata.Count)
		b, err := json.MarshalIndent(perfData, "", "  ")
		utils.LogOnError(err, "json marshal failed")
		err = os.WriteFile("perfData.json", b, 0644)
		utils.LogOnError(err, "write file failed")
	} else {
		//fmt.Println("TaskType != 0", mqMessage.TaskType)
	}
}
