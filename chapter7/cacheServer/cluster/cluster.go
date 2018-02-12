package cluster

import (
    "stathat.com/c/consistent"
    "time"
    "github.com/hashicorp/memberlist"
)

type Node interface {
    ShouldProcess(key string) bool
}

type node struct {
    circle *consistent.Consistent
    addr string
}

func New(addr, cluster string) (Node, error) {
    conf := memberlist.DefaultLocalConfig()
    conf.Name = addr
    conf.BindAddr = addr
    l, e := memberlist.Create(conf)
    if e != nil {
        return nil, e
    }
    clu := []string{cluster}
    _, e := l.Join(clu)
    if e != nil {
        return nil, e
    }
    circle = consistent.New()
    circle.Set(l.Members())
    go func() {
        for {
            time.Sleep(time.Second)
            circle.Set(l.Members())
        }
    }()
    return &node{circle, addr}
}

func (n *node) ShouldProcess(key string) bool {
    addr, _ := n.circle.Get(key)
    return addr == n.addr
}
