package task

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	Convey("After initialization", t, func() {
		Convey("client should be set", func() {
			So(Client(), ShouldNotBeNil)
		})
		Convey("the logger should be correctly parameterised", func() {
			So(logger.Prefix(), ShouldEqual, "[etcd-task]")
			So(logger.Flags(), ShouldEqual, log.LstdFlags)
		})
	})
}
