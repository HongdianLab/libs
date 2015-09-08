package task

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"path"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	Convey("After registering task test", t, func() {
		task := genTask("test-register")
		Convey("It should be available with etcd", func() {
			err := Register("test_register", task)
			So(err, ShouldBeNil)

			res, err := Client().Get("/tasks/test_register/"+task.Name, false, false)
			So(err, ShouldBeNil)

			h := &Task{}
			json.Unmarshal([]byte(res.Node.Value), &h)

			So(path.Base(res.Node.Key), ShouldEqual, task.Name)
			So(h, ShouldResemble, task)
		})

		Convey(fmt.Sprintf("And the ttl must be < %d", HEARTBEAT_DURATION), func() {
			Register("test2_register", task)
			res, err := Client().Get("/tasks/test2_register/"+task.Name, false, false)
			So(err, ShouldBeNil)
			now := time.Now()
			duration := res.Node.Expiration.Sub(now)
			So(duration, ShouldBeLessThanOrEqualTo, 3*HEARTBEAT_DURATION*time.Second)
		})
	})
}
