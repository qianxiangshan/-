package nothing

import (
	//	"container/heap"
	"container/list"
	"sync"
	//	"time"
)

//timeout handler
type TimeoutHandler func(data interface{})

type HeapNode struct {
	//current time
	Time       int64
	ListHeader list.List
	//using to protect ListHeader
	Locker sync.RWMutex
}

type TimerManager []*HeapNode

func (h TimerManager) Len() int           { return len(h) }
func (h TimerManager) Less(i, j int) bool { return h[i].Time < h[j].Time }
func (h TimerManager) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TimerManager) Push(x interface{}) {
	*h = append(*h, x.(*HeapNode))
}

func (h *TimerManager) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
