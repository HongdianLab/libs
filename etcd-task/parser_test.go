package task

import (
	"testing"

	"github.com/coreos/go-etcd/etcd"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	sampleNode = &etcd.Node{
		Key: "/tasks/test/example.org",
		Value: `
		{
			"Name": "example.org",
			"User": "user",
			"Password": "password",
			"Ports": {
				"http": "111"
			}
		}
		`,
	}
	sampleNodes = etcd.Nodes{sampleNode, sampleNode}
)

var (
	sampleResult = NewTask("example.org", nil)
)

func TestBuildTasksFromNodes(t *testing.T) {
	tasks := buildTasksFromNodes(sampleNodes)
	Convey("Given a sample response with 2 nodes, we got 2 tasks", t, func() {
		So(len(tasks), ShouldEqual, 2)
		So(tasks[0], ShouldResemble, sampleResult)
		So(tasks[1], ShouldResemble, sampleResult)
	})
}

func TestBuildTaskFromNode(t *testing.T) {
	task := buildTaskFromNode(sampleNode)
	Convey("Given a sample response, we got a filled Task", t, func() {
		So(task, ShouldResemble, sampleResult)
	})
}
