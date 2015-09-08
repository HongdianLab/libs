package task

import (
	"path"

	"github.com/coreos/go-etcd/etcd"
)

func Subscribe(name string) (<-chan *etcd.Response, <-chan *etcd.EtcdError) {
	stop := make(chan bool)
	responses := make(chan *etcd.Response)
	errors := make(chan *etcd.EtcdError)
	go func() {
		_, err := Client().Watch("/tasks/"+name, 0, true, responses, stop)
		if err != nil {
			errors <- err.(*etcd.EtcdError)
			close(errors)
			close(stop)
			return
		}
	}()
	return responses, errors
}

func SubscribeDown(name string) (<-chan string, <-chan *etcd.EtcdError) {
	expirations := make(chan string)
	responses, errors := Subscribe(name)
	go func() {
		for response := range responses {
			if response.Action == "expire" || response.Action == "delete" {
				expirations <- path.Base(response.Node.Key)
			}
		}
	}()
	return expirations, errors
}

func SubscribeNew(name string) (<-chan *Task, <-chan *etcd.EtcdError) {
	tasks := make(chan *Task)
	responses, errors := Subscribe(name)
	go func() {
		for response := range responses {
			if response.Action == "create" || (response.PrevNode == nil && response.Action == "set") {
				tasks <- buildTaskFromNode(response.Node)
			}
		}
	}()
	return tasks, errors
}

func SubscribeUpdate(name string) (<-chan *Task, <-chan *etcd.EtcdError) {
	tasks := make(chan *Task)
	responses, errors := Subscribe(name)
	go func() {
		for response := range responses {
			if response.Action == "update" || (response.PrevNode != nil && response.Action == "set") {
				tasks <- buildTaskFromNode(response.Node)
			}
		}
	}()
	return tasks, errors
}
