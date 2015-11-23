package filter

import (
	"container/list"
	"github.com/astaxie/beego/cache"
	"github.com/satori/go.uuid"
	"strconv"
)

type Packet struct {
	StreamId string
	Seq      uint64
	Index    uint64
	Data     interface{}
}

type waitingBuffer struct {
	nextIndex  uint64
	bufferList *list.List
}

type stream struct {
	id      string
	in      chan *Packet
	out     chan *Packet
	cache   cache.Cache
	lostMap map[string]*waitingBuffer
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
		id:      uuid.NewV4().String(),
		in:      input,
		out:     make(chan *Packet),
		cache:   mc,
		lostMap: make(map[string]*waitingBuffer),
	}
	s.run()
	return s.out, nil
}

func (this *stream) run() {
	go func() {
		for {
			select {
			case packet := <-this.in:
				packetId := this.id + packet.StreamId + strconv.FormatUint(packet.Seq, 10)
				if this.cache.Get(packetId) == nil {
					this.cache.Put(packetId, true, STREAM_TIME)
					if this.lostMap[packet.StreamId] == nil {
						this.lostMap[packet.StreamId] = &waitingBuffer{
							nextIndex:  packet.Index,
							bufferList: list.New(),
						}
					}
					if packet.Index == this.lostMap[packet.StreamId].nextIndex {
						this.out <- packet
						this.lostMap[packet.StreamId].nextIndex++
					} else {
						this.lostMap[packet.StreamId].push(packet)
						popPacket := this.lostMap[packet.StreamId].pop()
						for popPacket != nil {
							this.out <- popPacket
							popPacket = this.lostMap[packet.StreamId].pop()
						}
					}
					this.cache.Put(packetId, true, STREAM_TIME)
				}
			}
		}
	}()
}

func (wb *waitingBuffer) push(pak *Packet) {
	posPacket := wb.bufferList.Front()
	for posPacket != nil {
		if pak.Index < posPacket.Value.(*Packet).Index {
			wb.bufferList.InsertBefore(pak, posPacket)
			return
		}
		posPacket = posPacket.Next()
	}
	wb.bufferList.PushBack(pak)
}

func (wb *waitingBuffer) pop() (pak *Packet) {
	ele := wb.bufferList.Front()
	if ele == nil {
		pak = nil
		return
	}
	if ele.Value.(*Packet).Index == wb.nextIndex {
		pak = ele.Value.(*Packet)
		wb.bufferList.Remove(ele)
		wb.nextIndex++
		return
	}
	if wb.bufferList.Len() > 10 {
		pak = ele.Value.(*Packet)
		wb.bufferList.Remove(ele)
		wb.nextIndex = pak.Index + 1
		return
	}
	pak = nil
	return
}
