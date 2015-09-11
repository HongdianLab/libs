package task

import (
	"github.com/astaxie/beego/cache"

	"encoding/json"
	"fmt"
)

const (
	HEARTBEAT_DURATION = 3
)

var (
	mc cache.Cache
)

func init() {
	mc, _ = cache.NewCache("memory", `{"interval":60}`)
}

//register or update
func Register(name string, task *Task, params ...uint64) error {
	var expire uint64 = 3 * HEARTBEAT_DURATION
	if len(params) > 0 {
		expire = params[0]
	}
	if mc.Get(name+task.Name) == nil {
		err := register(name, task, expire)
		err = mc.Put(name+task.Name, true, int64(expire))
		return err
	}
	return nil
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
