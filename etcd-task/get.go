package task

func Get(task string) ([]*Task, error) {
	res, err := Client().Get("/tasks/"+task, false, true)
	if err != nil {
		if IsKeyNotFoundError(err) {
			return []*Task{}, nil
		}
		return nil, err
	}

	return buildTasksFromNodes(res.Node.Nodes), nil
}
