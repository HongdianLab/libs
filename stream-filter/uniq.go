package filter

import (
	"github.com/astaxie/beego/cache"
	"github.com/satori/go.uuid"

	"strconv"
)

type Packet struct {
	streamId string
	seq      int
	data     interface{}
}

type stream struct {
	id    string
	in    chan *Packet
	out   chan *Packet
	cache cache.Cache
}

var (
	STREAM_TIME int64 = 5
)

func NewUniq(input chan *Packet) (<-chan *Packet, error) {
	mc, err := cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		return nil, err
	}

	s := stream{
		id:    uuid.NewV4().String(),
		in:    input,
		out:   make(chan *Packet),
		cache: mc,
	}
	s.run()
	return s.out, nil
}

func (this *stream) run() {
	go func() {
		for {
			select {
			case packet := <-this.in:
				packetId := this.id + packet.streamId + strconv.Itoa(packet.seq)
				if this.cache.Get(packetId) == nil {
					this.out <- packet
					this.cache.Put(packetId, true, STREAM_TIME)
				}
			}
		}
	}()
}
