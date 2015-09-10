package task

import (
	"encoding/json"
	"fmt"
)

const (
	HEARTBEAT_DURATION = 3
)

//register or update
func Register(name string, task *Task, params ...uint64) error {
	var expire uint64 = 3 * HEARTBEAT_DURATION
	if len(params) > 0 {
		expire = params[0]
	}
	err := register(name, task, expire)
	return err
}

func register(name string, task *Task, expire uint64) error {
	key := fmt.Sprintf("/tasks/%s/%s", name, task.Name)
	taskJson, err := json.Marshal(&task)
	if err != nil {
		return err
	}
	value := string(taskJson)
	_, err = Client().Set(key, value, expire)
	return err
}
