ring
============================

Implements consistent hashing that can be used when
the number of server nodes can increase or decrease (like in memcached).
The hashing ring is built using hashring and etcd.

Using
============================

Importing ::

```go
import "github.com/Hongdianlan/libs/ring"
```

Basic example usage ::

```go
func init() {
    r= ring.New("servername", "serverport")
}

func isMy(topic string) bool {
    node, ok := ts.r.GetNode(topic)
    if !ok {
        return false
    }
    return node == ts.r.GetSelf()
} 
```
