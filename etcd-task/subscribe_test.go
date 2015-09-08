package task

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// Tests
func TestSubscribe(t *testing.T) {
	Convey("When we subscribe a task, we get all the notifications from it", t, func() {
		responses, _ := Subscribe("test_subs")
		time.Sleep(200 * time.Millisecond)
		Convey("When something happens about this task, the responses must be gathered in the channel", func() {
			_, err := Client().Create("/tasks/test_subs/key", "test", 0)
			So(err, ShouldBeNil)
			r := <-responses
			So(r, ShouldNotBeNil)
			So(r.Node.Key, ShouldEqual, "/tasks/test_subs/key")
			So(r.Action, ShouldEqual, "create")
			_, err = Client().Delete("/tasks/test_subs/key", false)
			So(err, ShouldBeNil)
			r = <-responses
			So(r, ShouldNotBeNil)
			So(r.Node.Key, ShouldEqual, "/tasks/test_subs/key")
			So(r.Action, ShouldEqual, "delete")
		})
	})
}

func TestSubscribeDown(t *testing.T) {
	Convey("When the task 'test' is watched and a task expired", t, func() {
		Register("test_expiration", genTask("test-expiration"))
		tasks, errs := SubscribeDown("test_expiration")
		Convey("The name of the disappeared task should be returned", func() {
			select {
			case err := <-errs:
				panic(err)
			case task, ok := <-tasks:
				So(task, ShouldEqual, "test-expiration")
				So(ok, ShouldBeTrue)
			}
		})
	})
}

func TestSubscribeNew(t *testing.T) {
	Convey("When the task 'test' is watched and a task registered", t, func() {
		tasks, _ := SubscribeNew("test_new")
		time.Sleep(200 * time.Millisecond)
		newTask := genTask("test-new")
		Register("test_new", newTask)
		Convey("A task should be available in the channel", func() {
			task, ok := <-tasks
			So(task, ShouldResemble, newTask)
			So(ok, ShouldBeTrue)
		})
	})
}

func TestSubscribeUpdate(t *testing.T) {
	Convey("When the task 'test' is watched and a task updates its data", t, func() {
		tasks, _ := SubscribeUpdate("test_upd")
		time.Sleep(200 * time.Millisecond)
		newTask := genTask("test-update")
		Register("test_upd", newTask)
		newTask.Params = map[string]interface{}{"foo": "bar"}
		Register("test_upd", newTask)

		Convey("A task should be available in the channel", func() {
			task, ok := <-tasks
			So(task, ShouldResemble, newTask)
			So(ok, ShouldBeTrue)
		})
	})
}
