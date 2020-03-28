package service

import (
	"chatroom/utils"
	"sync"
)

// 群组结构体
type Group struct {
	Name  string             // 名称
	Type  int                // 公共或者私有 0 / 1
	Admin string             // 群主
	List  map[string]*Client // 私有群的用户列表
}

type Groups struct {
	gs map[string]*Group
	sync.RWMutex
}

func (g Groups) GetGroup(name string) (group *Group, ok bool) {
	g.RLock()
	group, ok = g.gs[name]
	g.RUnlock()
	return
}

func (g Groups) SetGroup(name string, group *Group) {
	g.Lock()
	g.gs[name] = group
	g.Unlock()
}

func (g Groups) RangeGroups(f func(key string, group *Group)) {
	g.Lock()
	for key, value := range g.gs {
		f(key, value)
	}
	g.Unlock()
}

// 并发 map
type ConcurrentGroupMap struct {
	shards []*Groups
	len    int
}

func NewConcurrentGroupMap(num int) *ConcurrentGroupMap {
	m := make([]*Groups, num)
	for i := 0; i < num; i++ {
		m[i] = &Groups{
			gs: make(map[string]*Group),
		}
	}
	return &ConcurrentGroupMap{
		shards: m,
		len:    num,
	}
}

func (cgm ConcurrentGroupMap) GetShard(key string) *Groups {
	return cgm.shards[uint(utils.Fnv32(key))%uint(cgm.len)]
}

func (cgm ConcurrentGroupMap) Get(key string) (g *Group, ok bool) {
	shard := cgm.GetShard(key)
	g, ok = shard.GetGroup(key)
	return
}

func (cgm ConcurrentGroupMap) Set(key string, client *Group) {
	shard := cgm.GetShard(key)
	shard.SetGroup(key, client)
	return
}

func (cgm ConcurrentGroupMap) Range(f func(key string, client *Group)) {
	var wg sync.WaitGroup
	wg.Add(cgm.len)
	for _, shard := range cgm.shards {
		go func(s *Groups) {
			defer wg.Done()
			s.RangeGroups(f)
		}(shard)
	}

	wg.Wait()
}
