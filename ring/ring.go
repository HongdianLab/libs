package ring

import (
	"git.hongdian.fm/media/etcd-discovery/service"
	"github.com/HongdianLab/beego/cache"
	"github.com/HongdianLab/hashring"
	"github.com/coreos/go-etcd/etcd"

	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	logger *log.Logger
)

type Ring struct {
	sync.RWMutex
	selfServerId string
	servicename  string
	hashring     *hashring.HashRing
	stop         chan struct{}
	cache        cache.Cache
}

func New(servicename, serviceport string) (*Ring, error) {
	// XXX 在Mac上调试时，Mac没有eth0接口，会导致ring.New()失败。
	// 建议增加一个参数，允许用户传入指定的Interface（因此用户程序可以将该配置写入配置文件）。
	// 目前考虑到API兼容性，在上述方法获得一致同意前，不改动API，而是从环境变量读取。
	iface, ok := os.LookupEnv("RING_SERVER_INTERFACE")
	if !ok {
		iface = "eth0"
	}
	localIp, err := externalIP(iface)
	if err != nil {
		return nil, err
	}
	mc, err := cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		return nil, err
	}

	serverId := localIp + ":" + serviceport + ":" + strconv.Itoa(os.Getpid()) + ":" + strconv.FormatInt(startTime, 10)

	r := Ring{
		selfServerId: serverId,
		servicename:  servicename,
		hashring:     hashring.New([]string{serverId}),
		stop:         make(chan struct{}),
		cache:        mc,
	}
	r.register()
	r.subscribe()
	return &r, nil
}

func externalIP(name string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if strings.Compare(iface.Name, name) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

var startTime int64

func init() {
	startTime = time.Now().Unix()
	logger = log.New(os.Stderr, "[HongdianLab-ring]", log.LstdFlags)
}

func (this *Ring) GetNode(key string) (string, bool) {
	this.RLock()
	defer this.RUnlock()

	var server string
	ok := true
	value := this.cache.Get(key)
	if value == nil {
		server, ok = this.hashring.GetNode(key)
		this.cache.Put(key, server, 10)
	} else {
		server = value.(string)
	}
	return server, ok
}

func (this *Ring) GetSelf() string {
	return this.selfServerId
}

func (this *Ring) register() {
	host := &service.Host{
		Name:     this.selfServerId,
		Ports:    nil,
		User:     "user",
		Password: "secret",
	}
	service.Register(this.servicename, host, this.stop)
	this.refreshHashring()
}

func (this *Ring) subscribe() {
	newhosts, _ := service.SubscribeNew(this.servicename)
	go this.handleNew(newhosts)
	downhosts, _ := service.SubscribeDown(this.servicename)
	go this.handleDown(downhosts)
}

func (this *Ring) refreshHashring() {
	var serverIds []string
	serverIds = append(serverIds, this.selfServerId)
	hs, err := service.Get(this.servicename)

	for err != nil {
		errEtcd := err.(*etcd.EtcdError)
		logger.Println("Lost etcd registrationg for", this.servicename, ":", errEtcd.ErrorCode)
		time.Sleep(1 * time.Second)
		hs, err = service.Get(this.servicename)
		if err == nil {
			logger.Println("Recover etcd connection for", this.servicename)
		}
	}

	for idx, h := range hs {
		logger.Printf("No.%v host %v\n", idx, h.Name)
		serverIds = append(serverIds, h.Name)
	}

	this.Lock()
	defer this.Unlock()
	this.hashring = hashring.New(serverIds)
	this.cache.ClearAll()
}

func (this *Ring) addNode(serverId string) {
	this.refreshHashring()
}
func (this *Ring) removeNode(serverId string) {
	this.refreshHashring()
}

func (this *Ring) handleNew(hosts <-chan *service.Host) {
	for host := range hosts {
		logger.Printf("newserver join: %v\n", host.Name)
		this.addNode(host.Name)
	}
}

func (this *Ring) handleDown(hosts <-chan string) {
	for server := range hosts {
		logger.Printf("server down: %v\n", server)
		paths := strings.Split(server, "/")
		this.removeNode(paths[len(paths)-1])
	}
}
