package stat

import (
	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"github.com/satori/go.uuid"

	"log"
	"os"
	"time"
)

var (
	logger *log.Logger
)

type Key struct {
	name string
	tags map[string]string
	t    int64
	json *simplejson.Json
}
type Task struct {
	key   *Key
	value uint32
}

type Result map[string]*Task

type Stat struct {
	url   string
	dict  Result
	cycle int64
	in    chan *Task
	out   chan *simplejson.Json
}

func newKey(name string, tags map[string]string, t int64) *Key {
	json := simplejson.New()
	json.Set("metric", name)
	json.Set("timestamp", t)
	json.Set("tags", tags)
	k := &Key{
		name: name,
		tags: tags,
		t:    t,
		json: json,
	}
	return k
}
func NewStat(cycle int64, url string) (*Stat, error) {
	s := Stat{
		url:   url,
		dict:  make(map[string]*Task),
		cycle: cycle,
		in:    make(chan *Task),
		out:   make(chan *simplejson.Json),
	}
	go s.run()
	return &s, nil
}

func (this *Stat) Acc(name string, tags map[string]string, value uint32) {
	//logger.Printf("Acc: %v,%v,%v\n", name, tags, value)
	now := time.Now().Unix()
	t := (now/this.cycle + 1) * this.cycle

	task := &Task{
		key:   newKey(name, tags, t),
		value: value,
	}
	this.in <- task
}

func init() {
	logger = log.New(os.Stderr, "[HongdianLab-stat]", log.LstdFlags|log.Lshortfile)
}

func (this *Stat) send() {
	for {
		select {
		case js := <-this.out:
			json := simplejson.New()
			now := time.Now().Unix()
			json.Set("timestamp_send", now)
			json.Set("data", []*simplejson.Json{js})
			buf, err := json.Encode()
			if err != nil {
				//TODO error
				logger.Println("encode json error %v,%v\n", buf, err)
				continue
			}
			logger.Println(string(buf))
			request := gorequest.New()
			resp, body, errs := request.Post(this.url).
				Send(string(buf)).
				End()
			if errs != nil || resp.StatusCode != 200 {
				//TODO error
				logger.Printf("%v,%v,%v\n", resp.StatusCode, body, errs)
				continue
			}
		}
	}
}

func (this *Stat) run() {
	go this.send()
	timer := time.NewTicker(time.Duration(this.cycle/2) * time.Second)
	for {
		select {
		case task := <-this.in:
			keybyte, err := task.key.json.Encode()
			if err != nil {
				//TODO error
				logger.Printf("run %v,%v\n", keybyte, err)
				continue
			}
			keystr := string(keybyte)
			if v, ok := this.dict[keystr]; ok {
				v.value += task.value
				this.dict[keystr] = v
			} else {
				this.dict[keystr] = task
			}

		case <-timer.C:
			now := time.Now().Unix()
			base := (now / this.cycle) * this.cycle
			for k, v := range this.dict {
				if v.key.t < base {
					u := uuid.NewV1().String()
					js := v.key.json
					js.Set("uuid", u)
					js.Set("value", v.value)
					this.out <- js
					delete(this.dict, k)
				}
			}
		}
	}
}
