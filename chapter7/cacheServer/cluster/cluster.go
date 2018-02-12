package cluster

import (
	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
	"time"
)

type Node interface {
	ShouldProcess(key string) bool
}

type node struct {
	circle *consistent.Consistent
	addr   string
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
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	circle := consistent.New()
	go func() {
		for {
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil
}

func (n *node) ShouldProcess(key string) bool {
	addr, _ := n.circle.Get(key)
	return addr == n.addr
}
