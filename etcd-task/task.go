package task

type ParamType map[string]interface{}

type Task struct {
	Name   string    `json:"name"`
	Params ParamType `json:"params,omitempty"`
}

func NewTask(name string, params ParamType) *Task {
	h := &Task{Name: name, Params: params}
	return h
}
