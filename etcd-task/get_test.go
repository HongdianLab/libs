package task

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetNoTask(t *testing.T) {
	Convey("Without any task", t, func() {
		Convey("Get should return an empty slice", func() {
			tasks, err := Get("test_no_task")
			So(len(tasks), ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestGet(t *testing.T) {
	Convey("Given two registered tasks", t, func() {
		task1, task2 := genTask("task1"), genTask("task2")
		Register("test_task", task1)
		Register("test_task", task2)
		Convey("We should have 2 tasks", func() {
			tasks, err := Get("test_task")
			So(len(tasks), ShouldEqual, 2)
			So(tasks[0], ShouldResemble, task1)
			So(tasks[1], ShouldResemble, task2)
			So(err, ShouldBeNil)
		})
	})
}
