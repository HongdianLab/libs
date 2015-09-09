package filter

import (
	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"testing"
)

func TestUniq(t *testing.T) {
	Convey("uniq stream", t, func() {
		in := make(chan *Packet)
		out, err := NewUniq(in)
		So(out, ShouldNotBeNil)
		So(err, ShouldBeNil)

		var n int = 5
		go func() {
			for i := 0; i < n; i++ {
				p := &Packet{
					streamId: "111",
					seq:      i,
					data:     nil,
				}
				fmt.Printf("input: %v\n", p)
				in <- p
			}
		}()

		var m int = 0
		for j := 0; j < n; j++ {
			select {
			case p := <-out:
				fmt.Printf("output: %v\n", p)
				m++
			}
		}

		So(n, ShouldEqual, m)
	})
}

func TestDuplicate(t *testing.T) {
	Convey("uniq stream", t, func() {
		in := make(chan *Packet)
		out, err := NewUniq(in)
		So(out, ShouldNotBeNil)
		So(err, ShouldBeNil)

		var n int = 5
		go func() {
			for i := 0; i < n; i++ {
				p := &Packet{
					streamId: "111",
					seq:      i,
					data:     nil,
				}
				fmt.Printf("input: %v\n", p)
				in <- p
				fmt.Printf("input: %v\n", p)
				in <- p
			}
		}()

		var m int = 0
		for j := 0; j < n; j++ {
			select {
			case p := <-out:
				fmt.Printf("output: %v\n", p)
				m++
			}
		}

		So(n, ShouldEqual, m)
	})
}
