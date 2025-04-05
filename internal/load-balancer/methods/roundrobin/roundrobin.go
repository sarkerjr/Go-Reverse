package roundrobin

import (
	"slices"
	"sync"
)

type RoundRobin struct {
	serverList []string
	index      uint
	health     map[string]bool
	sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		serverList: make([]string, 0),
		index:      0,
		health:     make(map[string]bool),
	}
}

func (r *RoundRobin) GetNextServer() string {
	if len(r.serverList) == 0 {
		return ""
	}
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	server := r.serverList[r.index]
	r.index = (r.index + 1) % uint(len(r.serverList))
	return server
}

func (r *RoundRobin) AddServer(server string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.serverList = append(r.serverList, server)
	r.health[server] = true
}

func (r *RoundRobin) RemoveServer(server string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for i, s := range r.serverList {
		if s == server {
			r.serverList = slices.Delete(r.serverList, i, i+1)
			delete(r.health, server)
			return
		}
	}
}

func (r *RoundRobin) GetServers() []string {
	return r.serverList
}

func (r *RoundRobin) GetServerCount() int {
	return len(r.serverList)
}

func (r *RoundRobin) GetServerHealth() map[string]bool {
	return r.health
}

func (r *RoundRobin) SetServerHealth(server string, health bool) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.health[server] = health
}
