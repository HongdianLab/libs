package task

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaskUrl(t *testing.T) {
	Convey("Given a specific task", t, func() {
		task := NewTask("task", nil)
		So(task, ShouldNotBeNil)
	})
}
