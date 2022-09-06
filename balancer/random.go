package balancer

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
	//如果传入一个没有实现的函数，编译器不会报错，但是执行的时候会panic
	//var t Factory
	//factories[RandomBalancer] = t
}

type Random struct {
	sync.RWMutex
	hosts []string
	rnd   *rand.Rand
}

func NewRandom(hosts []string) Balancer {
	return &Random{hosts: hosts,
		rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (r *Random) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

func (r *Random) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

func (r *Random) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

func (r *Random) Inc(_ string) {}

func (r *Random) Done(_ string) {}
