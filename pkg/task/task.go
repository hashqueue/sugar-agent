package task

import (
	"encoding/json"
	"errors"

	"sugar-agent/internal"
	"sugar-agent/pkg/utils"
)

func StartTask(msg []byte) (*internal.PerfData, error) {
	mqMessage := make(map[string]interface{})
	err := json.Unmarshal(msg, &mqMessage)
	utils.FailOnError(err, "Failed to unmarshal message")
	taskType := mqMessage["task_type"].(float64)
	taskConfig := mqMessage["metadata"].(map[string]interface{})["task_config"].(map[string]interface{})

	if taskType == 0 {
		intervals := taskConfig["intervals"].(float64)
		count := taskConfig["count"].(float64)
		perfData, err := internal.StartGetPerfDataTask(uint64(intervals), uint64(count))
		if err != nil {
			return nil, errors.New("get perf data task failed")
		}
		//b, err := json.MarshalIndent(perfData, "", "  ")
		//utils.FailOnError(err, "json marshal failed")
		//err = os.WriteFile("perfData.json", b, 0644)
		//utils.FailOnError(err, "write file failed")
		return perfData, nil
	} else {
		//fmt.Println("TaskType != 0", mqMessage.TaskType)
	}
	return nil, errors.New("task type not supported")
}
