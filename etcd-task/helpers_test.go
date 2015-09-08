package task

func genTask(name string) *Task {
	// Empty if no arg, custom name otherways
	task := NewTask(name, nil)
	return task
}
