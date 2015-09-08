package task

import (
	"encoding/json"

	"github.com/coreos/go-etcd/etcd"
)

func buildTasksFromNodes(nodes etcd.Nodes) []*Task {
	tasks := make([]*Task, len(nodes))
	for i, node := range nodes {
		tasks[i] = buildTaskFromNode(node)
	}
	return tasks
}

func buildTaskFromNode(node *etcd.Node) *Task {
	task := &Task{}
	err := json.Unmarshal([]byte(node.Value), &task)
	if err != nil {
		panic(err)
	}
	return task
}
