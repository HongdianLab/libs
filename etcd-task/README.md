etcd-task
==============

Implement task manager using [etcd](https://github.com/coreos/etcd)

To install it:

`go get github.com/Hongdianlab/libs/etcd-task`

API
---

### Register a task

```go
func registerTask(taskname, servicename string) {
    t := task.NewTask(taskname, nil)                                                                                
    task.Register("task:"+servicename, t)
}
```

### Subscribe for a new task

```go
func subscribeNew(servicename string) {
    newTasks, err := task.SubscribeNew("task:" + servicename )
    if err==nil {
        for task := range newTasks {
            fmt.Println(task.Name, "has registered")
        }
    }
}
```

### Watch down tasks

```go
func subscribeDown(servicename string) {
    downTasks, err := task.SubscribeDown("task:" + servicename )
    if err==nil {
        for task := range downTasks {
            fmt.Println(task.Name, "has down")
        }
    }
}
```

### Get all tasks
```go
func getAll() {
    allTasks, err := task.Get("task:" + servicename )
    if err==nil {
        for task := range allTasks {
            fmt.Println(task.Name)
        }
    }
}
```
