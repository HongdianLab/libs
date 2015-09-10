package task

import (
	"encoding/json"
	"fmt"
)

const (
	HEARTBEAT_DURATION = 3
)

func RegisterOrUpdate(name string, task *Task) error {
	key := fmt.Sprintf("/tasks/%s/%s", name, task.Name)
	taskJson, err := json.Marshal(&task)
	if err != nil {
		return err
	}
	value := string(taskJson)
	_, err = Client().Set(key, value, 3*HEARTBEAT_DURATION)
	return err
}
