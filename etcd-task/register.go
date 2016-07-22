package task

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	HEARTBEAT_DURATION = 10
)

var (
	taskMap map[string]*Task
	taskC   chan *Task
)

func init() {
	taskMap = make(map[string]*Task)
	taskC = make(chan *Task)
	go run()
}

//register or update
func Register(servicename string, task *Task, params ...uint64) error {
	task.ServiceName = servicename
	var expire uint64 = 3 * HEARTBEAT_DURATION
	if len(params) > 0 {
		expire = params[0]
	}

	task.Expire = expire
	taskC <- task
	return nil
}

func register(task *Task) error {
	key := fmt.Sprintf("/tasks/%s/%s", task.ServiceName, task.Name)
	taskJson, err := json.Marshal(&task)
	if err != nil {
		return err
	}
	value := string(taskJson)
	_, err = Client().Set(key, value, task.Expire)
	return err
}

func run() {
	ticker := time.NewTicker(HEARTBEAT_DURATION * time.Second)
	for {
		select {
		case <-ticker.C:
			for key, task := range taskMap {
				register(task)
				delete(taskMap, key)
			}
		case task := <-taskC:
			taskMap[task.ServiceName+task.Name] = task
		}
	}
}
