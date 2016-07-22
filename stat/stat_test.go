package stat

import (
	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	Convey("stat data", t, func() {
		s, err := NewStat(5, "http://localhost")
		fmt.Printf("%v, %v\n", s, err)
		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		var n int = 10
		m := make(map[string]string)
		m["room_id"] = "123"
		m["user_id"] = "456"
		for i := 0; i < n; i++ {
			fmt.Printf("input: %v\n", i)
			s.Acc("test", m, 8192)
			time.Sleep(1 * time.Second)
		}
	})
}
